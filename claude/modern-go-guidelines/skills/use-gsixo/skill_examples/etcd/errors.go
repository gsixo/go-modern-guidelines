// Example: etcd sentinel errors + custom error types + Must pattern
// Patterns: sentinel vars с "etcdserver:" prefix, DiscoveryError, MustMarshal
package etcd

import (
	"errors"
	"fmt"
)

// ── Sentinel errors ───────────────────────────────────────────────────────────
// Стандарт etcd: префикс "etcdserver: " в строке ошибки.
// Naming: ErrXxx, package-level var.
// Используются через errors.Is() — поддерживают wrapping.

var (
	ErrStopped  = errors.New("etcdserver: server stopped")
	ErrNoLeader = errors.New("etcdserver: no leader")
	ErrTimeout  = errors.New("etcdserver: request timed out")
)

// Ошибки операций с членством
var (
	ErrIDRemoved  = errors.New("membership: ID removed")
	ErrIDExists   = errors.New("membership: ID exists")
	ErrIDNotFound = errors.New("membership: ID not found")
)

// ── Custom error type ─────────────────────────────────────────────────────────
// Используется когда нужен дополнительный структурированный контекст.
// Реализует error интерфейс через Error() string.
// Именование: XxxError struct (не ErrXxx — это для переменных).

type DiscoveryError struct {
	Op  string
	Err error
}

func (e DiscoveryError) Error() string {
	return fmt.Sprintf("discovery %s: %v", e.Op, e.Err)
}

// Unwrap позволяет errors.Is/As заглядывать внутрь
func (e DiscoveryError) Unwrap() error { return e.Err }

func newDiscoveryError(op string, err error) DiscoveryError {
	return DiscoveryError{Op: op, Err: err}
}

// ── Must — для операций которые не должны падать ──────────────────────────────
// Паттерн: panic при ошибке которая "никогда не должна случиться".
// Типичное применение: marshal известных типов, регистрация метрик.
// НЕ использовать для операций с пользовательским вводом или I/O.

type Marshaler interface {
	Marshal() ([]byte, error)
}

func MustMarshal(m Marshaler) []byte {
	d, err := m.Marshal()
	if err != nil {
		// panic с явным сообщением — легко найти по трассировке стека
		panic(fmt.Sprintf("marshal should never fail (%v)", err))
	}
	return d
}

// ── Проверки через errors.Is ──────────────────────────────────────────────────

func exampleErrorChecks(err error) {
	switch {
	case errors.Is(err, ErrNoLeader):
		// Специфичное — перед общим
	case errors.Is(err, ErrStopped):
		// ...
	case err != nil:
		// Общий случай — последним
	}
}
