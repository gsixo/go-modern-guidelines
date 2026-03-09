// Example: etcd configuration patterns
// Patterns: struct-based config, named constants for defaults, functional options
package etcd

import "time"

// ── Named constants для дефолтов ──────────────────────────────────────────────
// Экспортируемые константы с типами — часть публичного API.
// Пользователи могут ссылаться на них по имени вместо magic numbers.
// Именование: DefaultXxx, PascalCase.

const (
	DefaultMaxTxnOps            = uint(128)
	DefaultMaxRequestBytes      = 1.5 * 1024 * 1024 // 1.5 MiB
	DefaultWarningApplyDuration = 100 * time.Millisecond
	DefaultDialTimeout          = 5 * time.Second
	DefaultAutoSyncInterval     = 0 // 0 disables auto-sync
)

// ── Struct-based Config ────────────────────────────────────────────────────────
// Используй когда конфигурация сложная или передаётся в несколько функций.
// Ноль-значение должно быть разумным ("zero is valid").
// Комментируй "0 disables" / "empty = default" для нулевых значений.

type Config struct {
	Endpoints   []string
	DialTimeout time.Duration // 0 uses DefaultDialTimeout

	// AutoSyncInterval — интервал синхронизации эндпоинтов.
	// 0 отключает автосинхронизацию.
	AutoSyncInterval time.Duration

	// MaxCallSendMsgSize максимальный размер gRPC сообщения.
	// 0 использует 2 MiB по умолчанию.
	MaxCallSendMsgSize int
}

// ── Backend Config с Functional Options ───────────────────────────────────────
// Functional options: тип BackendConfigOption func(*BackendConfig).
// Функция-конструктор принимает ...BackendConfigOption.
// Преимущества: легко добавить новые опции без изменения сигнатуры NewBackend.

type BackendConfig struct {
	// Path до файла с базой данных
	Path string
	// MmapSize — размер mmap в байтах; 0 = auto
	MmapSize uint64
	// BatchInterval — max время ожидания перед flush батча
	BatchInterval time.Duration
	// BatchLimit — max количество операций в батче
	BatchLimit int
}

// BackendConfigOption — тип для функциональных опций
type BackendConfigOption func(*BackendConfig)

func WithMmapSize(size uint64) BackendConfigOption {
	return func(bcfg *BackendConfig) { bcfg.MmapSize = size }
}

func WithBatchInterval(d time.Duration) BackendConfigOption {
	return func(bcfg *BackendConfig) { bcfg.BatchInterval = d }
}

func WithBatchLimit(limit int) BackendConfigOption {
	return func(bcfg *BackendConfig) { bcfg.BatchLimit = limit }
}

// DefaultBackendConfig возвращает разумные дефолты.
// Располагается рядом со struct — легко найти.
func DefaultBackendConfig(path string) BackendConfig {
	return BackendConfig{
		Path:          path,
		BatchInterval: 100 * time.Millisecond,
		BatchLimit:    10000,
	}
}

// Backend — интерфейс бэкенда хранилища
type Backend interface {
	Close() error
	ForceCommit()
}

// NewBackend создаёт Backend с опциональными overrides.
// Паттерн: cfg передаётся по значению (копия), opts мутируют копию.
func NewBackend(cfg BackendConfig, opts ...BackendConfigOption) Backend {
	for _, opt := range opts {
		opt(&cfg) // мутируем копию, не исходный объект
	}
	return &backendImpl{cfg: cfg}
}

type backendImpl struct{ cfg BackendConfig }

func (b *backendImpl) Close() error    { return nil }
func (b *backendImpl) ForceCommit()    {}

// Использование:
//
//   b := NewBackend(DefaultBackendConfig("/var/etcd/data"),
//       WithMmapSize(10 * 1024 * 1024 * 1024),  // 10 GiB
//       WithBatchInterval(200 * time.Millisecond),
//   )
