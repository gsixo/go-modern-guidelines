// Example: staging/src/k8s.io/apimachinery/pkg/util/errors/errors.go
// Patterns: Aggregate interface, NewAggregate, AggregateGoroutines,
//           usage patterns for collecting multiple errors
package errors

import (
	"errors"
	"fmt"
	"sync"
)

// ── Aggregate interface ───────────────────────────────────────────────────────
// Aggregate объединяет несколько ошибок в одну.
// Реализует error — можно передавать туда где ожидается error.
// Is() позволяет errors.Is() работать с содержимым агрегата.

type Aggregate interface {
	error
	Errors() []error
	Is(error) bool
}

// ── NewAggregate ──────────────────────────────────────────────────────────────
// Ключевые свойства:
//   - возвращает nil если список пустой или все элементы nil
//   - фильтрует nil-ошибки из входного списка
//   - удобен в паттерне: собрать ошибки → return NewAggregate(errs)

func NewAggregate(errList []error) Aggregate {
	var errs []error
	for _, e := range errList {
		if e != nil {
			errs = append(errs, e)
		}
	}
	if len(errs) == 0 {
		return nil // nil Aggregate — это nil error
	}
	return aggregate(errs)
}

type aggregate []error

func (agg aggregate) Error() string {
	if len(agg) == 1 {
		return agg[0].Error()
	}
	msg := fmt.Sprintf("[%s", agg[0].Error())
	for _, e := range agg[1:] {
		msg += fmt.Sprintf(", %s", e.Error())
	}
	return msg + "]"
}

func (agg aggregate) Errors() []error { return []error(agg) }

func (agg aggregate) Is(target error) bool {
	for _, e := range agg {
		if errors.Is(e, target) {
			return true
		}
	}
	return false
}

// ── AggregateGoroutines ───────────────────────────────────────────────────────
// Запускает funcs в параллельных горутинах, собирает все ошибки.
// Аналог errgroup, но возвращает Aggregate с доступом ко всем ошибкам.

func AggregateGoroutines(funcs ...func() error) Aggregate {
	var wg sync.WaitGroup
	errCh := make(chan error, len(funcs))

	for _, f := range funcs {
		wg.Add(1)
		go func(fn func() error) {
			defer wg.Done()
			if err := fn(); err != nil {
				errCh <- err
			}
		}(f)
	}

	wg.Wait()
	close(errCh)

	var errs []error
	for e := range errCh {
		errs = append(errs, e)
	}
	return NewAggregate(errs)
}

// ── Типичные паттерны использования ──────────────────────────────────────────

// Паттерн 1: цикл с накоплением ошибок
func validateAll(items []string, validate func(string) error) error {
	var errs []error
	for _, item := range items {
		if err := validate(item); err != nil {
			errs = append(errs, fmt.Errorf("item %q: %w", item, err))
		}
	}
	return NewAggregate(errs) // nil если errs пустой
}

// Паттерн 2: несколько шагов, все должны выполниться
func closeAll(closers []func() error) error {
	var errs []error
	for i, close := range closers {
		if err := close(); err != nil {
			errs = append(errs, fmt.Errorf("closer[%d]: %w", i, err))
			// НЕ break — продолжаем закрывать остальные
		}
	}
	return NewAggregate(errs)
}

// Паттерн 3: параллельная проверка
func validateParallel(validators []func() error) error {
	return AggregateGoroutines(validators...)
}

// Паттерн 4: проверка через errors.Is с агрегатом
func containsNotFound(err error) bool {
	return errors.Is(err, errNotFound) // работает рекурсивно через Is()
}

var errNotFound = errors.New("not found")
