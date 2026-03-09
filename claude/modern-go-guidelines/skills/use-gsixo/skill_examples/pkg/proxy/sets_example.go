// Example: pkg/proxy/endpointschangetracker_test.go + nftables/proxier.go
// Patterns: sets.New[string](), typed set operations in production and test code
package proxy

import (
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
)

// ── sets.Set[T] в production-коде ────────────────────────────────────────────
// pkg/proxy/nftables/proxier.go использует sets для отслеживания локальных IP.
// Преимущество перед map[string]struct{}: читаемые операции (Union, Difference, Has).

type endpointTracker struct {
	// map: NamespacedName → множество IP-адресов этого эндпоинта
	localIPs map[types.NamespacedName]sets.Set[string]
}

func newEndpointTracker() *endpointTracker {
	return &endpointTracker{
		localIPs: make(map[types.NamespacedName]sets.Set[string]),
	}
}

func (t *endpointTracker) update(key types.NamespacedName, ips []string) {
	// Создаём Set из слайса — O(n), без дублей
	t.localIPs[key] = sets.New[string](ips...)
}

func (t *endpointTracker) isLocal(key types.NamespacedName, ip string) bool {
	s, ok := t.localIPs[key]
	if !ok {
		return false
	}
	return s.Has(ip) // O(1) поиск
}

func (t *endpointTracker) allLocalIPs() sets.Set[string] {
	result := sets.New[string]()
	for _, ips := range t.localIPs {
		result = result.Union(ips) // объединяем все множества
	}
	return result
}

// ── sets.Set[T] в тестах ──────────────────────────────────────────────────────
// pkg/proxy/endpointschangetracker_test.go: expected-значения как sets.

func exampleTestSetUsage() {
	// Было: ручное сравнение map[string]struct{}
	// expected := map[string]struct{}{"1.1.1.1": {}, "2.2.2.2": {}}
	// for k := range got { if _, ok := expected[k]; !ok { /* error */ } }

	// Стало: читаемое декларативное expected
	expectedLocalIPs := map[types.NamespacedName]sets.Set[string]{
		{Namespace: "ns1", Name: "ep1"}: sets.New[string]("1.1.1.1"),
		{Namespace: "ns2", Name: "ep2"}: sets.New[string]("2.2.2.2", "2.2.2.3"),
		{Namespace: "ns3", Name: "ep3"}: sets.New[string]("3.3.3.3", "3.3.3.30", "3.3.3.31"),
	}

	_ = expectedLocalIPs
}

// ── Полный список операций sets.Set[T] ───────────────────────────────────────

func setOperationsCheatsheet() {
	// Создание
	s1 := sets.New[string]("a", "b", "c")
	s2 := sets.New[string]("b", "c", "d")

	_ = s1.Has("a")               // true
	_ = s1.HasAll("a", "b")       // true
	_ = s1.HasAny("x", "a")       // true

	s1.Insert("d", "e")           // добавить элементы
	s1.Delete("a")                // удалить элемент

	_ = s1.Union(s2)              // {a,b,c,d,e} — все из обоих
	_ = s1.Intersection(s2)       // {b,c,d}     — только общие
	_ = s1.Difference(s2)         // {a,e}       — только в s1
	_ = s1.IsSuperset(s2)         // все элементы s2 есть в s1?
	_ = s1.Equal(s2)              // одинаковые множества?
	_ = s1.Len()                  // количество элементов

	list := s1.UnsortedList()     // []string (порядок не определён)
	sorted := sets.List(s1)       // []string (отсортированный)

	_ = list
	_ = sorted
}
