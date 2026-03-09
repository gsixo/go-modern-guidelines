// Example: pkg/kubelet/container/sync_result.go
// Patterns: sentinel errors, typed error with context, error aggregation, %w wrapping
package container

import (
	"errors"
	"fmt"
	"time"

	utilerrors "k8s.io/apimachinery/pkg/util/errors"
)

// ── Sentinel errors ──────────────────────────────────────────────────────────
// Package-level vars с Err-префиксом.
// Сравниваются через errors.Is(), не через ==.

var (
	ErrCrashLoopBackOff   = errors.New("CrashLoopBackOff")
	ErrContainerNotFound  = errors.New("no matching container")
	ErrRunContainer       = errors.New("RunContainerError")
	ErrKillContainer      = errors.New("KillContainerError")
	ErrVerifyNonRoot      = errors.New("VerifyNonRootError")
	ErrPreStartHook       = errors.New("PreStartHookError")
	ErrPostStartHook      = errors.New("PostStartHookError")
)

// ── Typed error with metadata ─────────────────────────────────────────────────
// Структура-ошибка несёт дополнительный контекст (время backoff),
// нужный вызывающему коду — без необходимости парсить строку.

type BackoffError struct {
	error
	backoffTime time.Time
}

func NewBackoffError(err error, backoffTime time.Time) *BackoffError {
	return &BackoffError{
		error:       err,
		backoffTime: backoffTime,
	}
}

func (e *BackoffError) BackoffTime() time.Time { return e.backoffTime }

// ── Error aggregation ─────────────────────────────────────────────────────────
// utilerrors.NewAggregate объединяет несколько ошибок в одну.
// Возвращает nil если список пустой — удобно в конце цикла сборки ошибок.

type PodSyncResult struct {
	SyncError   error
	SyncResults []*SyncResult
}

type SyncResult struct {
	Action  string
	Target  string
	Error   error
	Message string
}

func (p *PodSyncResult) Error() error {
	var errlist []error
	if p.SyncError != nil {
		// %w сохраняет цепочку для errors.Is / errors.As
		errlist = append(errlist, fmt.Errorf("failed to SyncPod: %w", p.SyncError))
	}
	for _, result := range p.SyncResults {
		if result.Error != nil {
			errlist = append(errlist, fmt.Errorf(
				"failed to %q for %q with %w: %q",
				result.Action, result.Target, result.Error, result.Message,
			))
		}
	}
	return utilerrors.NewAggregate(errlist)
}

// ── Recursive unwrap для поиска BackoffError ─────────────────────────────────
// errors.As обходит цепочку через Unwrap; для Aggregate — рекурсия вручную.

func MinBackoffExpiration(err error) (time.Time, bool) {
	var be *BackoffError
	var ae utilerrors.Aggregate
	switch {
	case errors.As(err, &be):
		return be.BackoffTime(), true
	case errors.As(err, &ae):
		var minTime time.Time
		found := false
		for _, e := range ae.Errors() {
			if t, ok := MinBackoffExpiration(e); ok {
				if !found || t.Before(minTime) {
					minTime = t
					found = true
				}
			}
		}
		return minTime, found
	default:
		if e := errors.Unwrap(err); e != nil {
			return MinBackoffExpiration(e)
		}
		return time.Time{}, false
	}
}
