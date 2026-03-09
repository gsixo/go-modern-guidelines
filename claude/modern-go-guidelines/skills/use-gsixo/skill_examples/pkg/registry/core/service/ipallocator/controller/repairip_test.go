// Example: pkg/registry/core/service/ipallocator/controller/repairip_test.go
// Patterns: testclock.NewFakeClock, PrependReactor для side-effects,
//           GetIndexer().Add(), Action verification
package controller

import (
	"context"
	"testing"
	"time"

	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
	testclock "k8s.io/utils/clock/testing"
)

// ── FakeClock — детерминированное время в тестах ──────────────────────────────
// Проблема: тесты с time.Now() нестабильны и медленны.
// Решение: инъекция clock через поле или параметр конструктора.
//
// testclock.NewFakeClock(t time.Time) — возвращает часы застывшие в t.
// fakeClock.Step(d) — переводит время вперёд.
// fakeClock.SetTime(t) — устанавливает конкретное время.

func TestRepairWithFakeClock(t *testing.T) {
	// Фиксированное время — тест детерминирован
	fixedTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	fakeClock := testclock.NewFakeClock(fixedTime)

	// Инъекция clock в контроллер
	ctrl := &repairController{
		clock: fakeClock,
	}

	// Проверяем поведение ДО истечения deadline
	_ = ctrl.isExpired(fixedTime.Add(-time.Hour))

	// Переводим время вперёд — симулируем истечение TTL
	fakeClock.Step(2 * time.Hour)

	// Теперь проверяем поведение ПОСЛЕ истечения
	_ = ctrl.isExpired(fixedTime.Add(-time.Hour))
}

// ── PrependReactor — side-effects при API-вызовах ─────────────────────────────
// Проблема: fake client не обновляет indexer при Create/Update.
// Решение: PrependReactor перехватывает вызов, добавляет в indexer вручную,
//          возвращает false — продолжить обработку оригинальным fake handler.
//
// PrependReactor(verb, resource, fn) — добавляет reactor ПЕРЕД дефолтными.
// fn возвращает (handled bool, obj runtime.Object, err error).
// handled=false → передать следующему reactor.
// handled=true → остановиться, вернуть obj, err.

func TestRepairWithReactor(t *testing.T) {
	fakeClient := fake.NewSimpleClientset()
	factory := informers.NewSharedInformerFactory(fakeClient, 0)
	ipInformer := factory.Networking().V1().IPAddresses()
	ipIndexer := ipInformer.Informer().GetIndexer()

	// PrependReactor: при создании IPAddress — добавить в indexer
	fakeClient.PrependReactor("create", "ipaddresses",
		func(action k8stesting.Action) (bool, runtime.Object, error) {
			ip := action.(k8stesting.CreateAction).GetObject().(*networkingv1.IPAddress)
			// Добавляем в indexer синхронно — доступно сразу
			if err := ipIndexer.Add(ip); err != nil {
				t.Errorf("failed to add to indexer: %v", err)
			}
			return false, ip, nil // false = продолжить обработку fake client
		},
	)

	ctx := context.Background()

	// Создаём через API
	ip := &networkingv1.IPAddress{
		ObjectMeta: metav1.ObjectMeta{Name: "10.0.0.1"},
	}
	_, err := fakeClient.NetworkingV1().IPAddresses().Create(ctx, ip, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// Теперь объект доступен через indexer (без ожидания синхронизации informer)
	item, exists, err := ipIndexer.GetByKey("10.0.0.1")
	if err != nil || !exists {
		t.Errorf("expected item in indexer: exists=%v, err=%v", exists, err)
	}
	_ = item
}

// ── GetIndexer().Add() для setup без API-вызовов ─────────────────────────────
// Добавляй объекты напрямую в indexer — быстрее и проще чем через fake client.
// Используй когда нужно предзаполнить кэш без записи в "API".

func TestWithPrefilledIndexer(t *testing.T) {
	fakeClient := fake.NewSimpleClientset()
	factory := informers.NewSharedInformerFactory(fakeClient, 0)
	ipInformer := factory.Networking().V1().IPAddresses()

	// Предзаполняем indexer напрямую — не через API
	existingIP := &networkingv1.IPAddress{
		ObjectMeta: metav1.ObjectMeta{Name: "192.168.1.1"},
	}
	_ = ipInformer.Informer().GetIndexer().Add(existingIP)

	// Теперь lister видит объект без запуска informer
	result, err := ipInformer.Lister().Get("192.168.1.1")
	if err != nil {
		t.Fatalf("expected to find IP: %v", err)
	}
	if result.Name != "192.168.1.1" {
		t.Errorf("wrong name: %s", result.Name)
	}
}

// ── stubs ─────────────────────────────────────────────────────────────────────

type repairController struct {
	clock interface {
		Now() time.Time
	}
}

func (c *repairController) isExpired(creationTime time.Time) bool {
	return c.clock.Now().Sub(creationTime) > time.Hour
}
