// Example: etcd Prometheus metrics patterns
// Patterns: init() registration, Histogram with ExponentialBuckets,
//           CounterVec with labels, expvar for runtime state
package etcd

import (
	"expvar"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// ── Histogram для задержек ────────────────────────────────────────────────────
// ExponentialBuckets(start, factor, count) — хорошо для задержек (логарифмическая шкала).
// Namespace + Subsystem + Name — формат: "etcd_disk_backend_commit_duration_seconds".
// Регистрация в init() — выполняется при импорте пакета, до main().

var commitSec = prometheus.NewHistogram(prometheus.HistogramOpts{
	Namespace: "etcd",
	Subsystem: "disk",
	Name:      "backend_commit_duration_seconds",
	Help:      "The latency distributions of commit called by backend.",
	Buckets:   prometheus.ExponentialBuckets(0.001, 2, 14), // 1ms → 8.192s
})

// ── Counter для событий ───────────────────────────────────────────────────────
var applyDurations = prometheus.NewHistogram(prometheus.HistogramOpts{
	Namespace: "etcd",
	Subsystem: "server",
	Name:      "apply_duration_seconds",
	Help:      "The latency distributions of v2 apply called by backend.",
	Buckets:   prometheus.ExponentialBuckets(0.0001, 2, 20),
})

// ── CounterVec с label-ами ────────────────────────────────────────────────────
// Labels = дополнительные измерения. Осторожно с высокой кардинальностью.
// Статические label-значения (enum, status-коды) — ok.
// Динамические (user ID, IP) — нет.

var promoteFailed = prometheus.NewCounterVec(prometheus.CounterOpts{
	Namespace: "etcd",
	Subsystem: "server",
	Name:      "learner_promote_failed_total",
	Help:      "The total number of failed learner promotions (likely leader changed) while this member is leader.",
}, []string{"reason"})

// ── Регистрация в init() ──────────────────────────────────────────────────────
// init() вызывается один раз при импорте пакета.
// prometheus.MustRegister паникует при duplicate — защита от случайного дублирования.
// Альтернатива: prometheus.Register + обработка ошибки для graceful degradation.

func init() {
	prometheus.MustRegister(commitSec)
	prometheus.MustRegister(applyDurations)
	prometheus.MustRegister(promoteFailed)
}

// ── Использование метрик ──────────────────────────────────────────────────────

type diskBackend struct{}

func (b *diskBackend) commit() {
	start := time.Now()
	defer func() {
		commitSec.Observe(time.Since(start).Seconds())
	}()
	// реальная работа...
}

func observePromoteFailed(reason string) {
	promoteFailed.WithLabelValues(reason).Inc()
}

// ── expvar для runtime state ──────────────────────────────────────────────────
// expvar доступен через /debug/vars endpoint (HTTP).
// Используется для сложного состояния которое нельзя выразить через Prometheus.
// Func откладывает вычисление до момента запроса — не обновляется постоянно.

var raftStatusMu sync.Mutex

func raftStatus() any {
	return map[string]any{
		"state": "leader",
		"term":  42,
	}
}

func init() {
	expvar.Publish("raft.status", expvar.Func(func() any {
		raftStatusMu.Lock()
		defer raftStatusMu.Unlock()
		return raftStatus()
	}))
}
