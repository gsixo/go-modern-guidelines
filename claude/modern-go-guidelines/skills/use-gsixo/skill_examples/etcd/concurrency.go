// Example: etcd concurrency patterns
// Patterns: typed atomics, fine-grained RWMutex, Notifier broadcast, ID-based Wait
package etcd

import (
	"sync"
	"sync/atomic"
)

// ── Typed atomics вместо sync/atomic функций ──────────────────────────────────
// atomic.Uint64, atomic.Bool, atomic.Pointer[T] — Go 1.19+
// Преимущества: type-safe, нет случайной передачи неверного значения,
// легче читается, не требует unsafe.Pointer.

type server struct {
	// ── Atomics — в НАЧАЛЕ структуры ──────────────────────────────────────────
	// Обеспечивает правильное 64-битное выравнивание на 32-битных платформах.
	// Комментируй что именно хранится и кем читается/пишется.
	appliedIndex   atomic.Uint64  // индекс последнего применённого лог-entry
	committedIndex atomic.Uint64  // индекс последнего коммит-entry
	isLeader       atomic.Bool    // true если этот узел является лидером

	// ── Fine-grained RWMutex ──────────────────────────────────────────────────
	// Один RWMutex на одно поле/группу данных.
	// НЕ один большой mutex на весь объект — это bottleneck.
	// Комментируй ТОЧНО что защищает каждый mutex.
	lgMu   sync.RWMutex // protects lg (logger replacement)
	readMu sync.RWMutex // protects readwaitc and readNotifier
	bemu   sync.RWMutex // protects backend
}

// ── Notifier pattern ───────────────────────────────────────────────────────────
// Thread-safe broadcast уведомления без закрытия оригинального канала.
// Принцип: при Notify() создаём НОВЫЙ канал, старый закрываем (все читатели проснутся).
// Читатели хранят ссылку на канал — при следующем Notify() получат новый канал.

type Notifier struct {
	mu      sync.RWMutex
	channel chan struct{}
}

func NewNotifier() *Notifier {
	return &Notifier{channel: make(chan struct{})}
}

// Receive возвращает текущий канал для ожидания.
// Вызывающий должен ждать закрытия канала, потом вызвать Receive снова.
func (n *Notifier) Receive() <-chan struct{} {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.channel
}

// Notify будит всех текущих ожидающих и создаёт новый канал для следующего раунда.
func (n *Notifier) Notify() {
	newCh := make(chan struct{})
	n.mu.Lock()
	ch := n.channel
	n.channel = newCh
	n.mu.Unlock()
	close(ch) // будит всех кто ждёт на старом канале
}

// Пример использования Notifier:
//
//   for {
//       ch := notifier.Receive()
//       select {
//       case <-ch:
//           // получили уведомление — обработать
//       case <-ctx.Done():
//           return ctx.Err()
//       }
//   }

// ── ID-based Wait ─────────────────────────────────────────────────────────────
// Асинхронная координация по числовому ID.
// Используется в etcd для ожидания результата конкретного запроса
// в raft-цепочке: запрос → ID → ждём → получаем ответ.

// Wait позволяет горутинам ждать результата по ID и триггерить его.
type Wait interface {
	// Register регистрирует ожидание ID; возвращает канал для получения результата.
	Register(id uint64) <-chan any
	// Trigger отправляет результат x всем ждущим с данным ID.
	Trigger(id uint64, x any)
	// IsRegistered возвращает true если кто-то ждёт на id.
	IsRegistered(id uint64) bool
}

type waitList struct {
	mu sync.Mutex
	l  map[uint64]chan any
}

func NewWait() Wait {
	return &waitList{l: make(map[uint64]chan any)}
}

func (w *waitList) Register(id uint64) <-chan any {
	w.mu.Lock()
	defer w.mu.Unlock()
	ch := make(chan any, 1)
	if _, ok := w.l[id]; !ok {
		w.l[id] = ch
	}
	return w.l[id]
}

func (w *waitList) Trigger(id uint64, x any) {
	w.mu.Lock()
	ch := w.l[id]
	delete(w.l, id)
	w.mu.Unlock()
	if ch != nil {
		ch <- x
	}
}

func (w *waitList) IsRegistered(id uint64) bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	_, ok := w.l[id]
	return ok
}

// ── Lock-free ID generation ───────────────────────────────────────────────────
// Используется в etcd для генерации уникальных ID без мьютекса.
// prefix уникален для каждого экземпляра, suffix монотонно растёт.

const suffixMask = ^uint64(0) >> 20 // нижние 44 бита

type Generator struct {
	prefix uint64 // уникален для узла/процесса
	suffix atomic.Uint64
}

func (g *Generator) Next() uint64 {
	return g.prefix | (g.suffix.Add(1) & suffixMask)
}
