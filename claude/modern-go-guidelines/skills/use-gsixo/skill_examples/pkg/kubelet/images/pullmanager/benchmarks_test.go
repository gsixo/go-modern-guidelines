// Example: pkg/kubelet/images/image_manager_test.go (benchmarks)
// Patterns: b.Loop() Go 1.24+, BenchmarkXxx naming, b.ReportAllocs()
package pullmanager

import (
	"testing"
)

// ── Benchmark pattern — Go 1.24+ ──────────────────────────────────────────────
// b.Loop() — новый способ (Go 1.24+); заменяет for i := 0; i < b.N; i++
// Преимущества b.Loop():
//   - корректно учитывает время setup/teardown за пределами цикла
//   - автоматически управляет b.N итерациями

func BenchmarkImageManagerAdd(b *testing.B) {
	// Setup — НЕ входит в замер
	manager := newTestImageManager()

	b.ReportAllocs() // показывает allocations/op и B/op в отчёте
	b.ResetTimer()   // сбрасываем таймер после setup

	for b.Loop() { // Go 1.24+
		manager.AddImage("sha256:abc123", "nginx:latest")
	}
}

// ── Старый стиль — Go < 1.24 ─────────────────────────────────────────────────

func BenchmarkImageManagerAddLegacy(b *testing.B) {
	manager := newTestImageManager()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ { // до Go 1.24
		manager.AddImage("sha256:abc123", "nginx:latest")
	}
}

// ── Параллельный бенчмарк ─────────────────────────────────────────────────────
// RunParallel используется для измерения contention при конкурентном доступе.
// -benchtime=Xs или -benchtime=NxB задают продолжительность/итерации.

func BenchmarkImageManagerConcurrent(b *testing.B) {
	manager := newTestImageManager()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			manager.GetImage("sha256:abc123")
		}
	})
}

// ── Sub-benchmarks ────────────────────────────────────────────────────────────
// Позволяют сравнить несколько реализаций или вариантов данных в одном прогоне.
// Запуск конкретного: go test -bench=BenchmarkImageOps/large

func BenchmarkImageOps(b *testing.B) {
	sizes := []struct {
		name  string
		count int
	}{
		{"small", 10},
		{"medium", 100},
		{"large", 1000},
	}

	for _, s := range sizes {
		b.Run(s.name, func(b *testing.B) {
			manager := newTestImageManagerWithN(s.count)
			b.ResetTimer()
			for b.Loop() {
				manager.ListImages()
			}
		})
	}
}

// ── Тест на panic ─────────────────────────────────────────────────────────────
// Паттерн из etcd: проверяем что функция паникует при неверных аргументах.
// defer+recover должен быть ПЕРВЫМ в функции теста.

func TestMustMarshalPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic but did not panic")
		}
	}()

	// Вызываем что должно паниковать
	mustNotFail(nil)
}

func mustNotFail(v any) {
	if v == nil {
		panic("must not be nil")
	}
}

// ── helpers ───────────────────────────────────────────────────────────────────

type testImageManager struct {
	images map[string]string
}

func newTestImageManager() *testImageManager {
	return &testImageManager{images: make(map[string]string)}
}

func newTestImageManagerWithN(n int) *testImageManager {
	m := newTestImageManager()
	for i := 0; i < n; i++ {
		m.images[string(rune('a'+i))] = "tag"
	}
	return m
}

func (m *testImageManager) AddImage(id, tag string)    { m.images[id] = tag }
func (m *testImageManager) GetImage(id string) string  { return m.images[id] }
func (m *testImageManager) ListImages() []string {
	out := make([]string, 0, len(m.images))
	for k := range m.images {
		out = append(out, k)
	}
	return out
}
