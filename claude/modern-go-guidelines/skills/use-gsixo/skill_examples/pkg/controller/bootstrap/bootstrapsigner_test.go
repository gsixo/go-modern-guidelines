// Example: pkg/controller/bootstrap/bootstrapsigner_test.go
// Patterns: factory function, fake.NewSimpleClientset, GetIndexer().Add(),
//           verifyActions helper, action-based test verification
package bootstrap

import (
	"context"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	core "k8s.io/client-go/testing"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"k8s.io/apiserver/pkg/storage/names"
)

// ── Factory function ──────────────────────────────────────────────────────────
// Одна функция создаёт всё нужное для теста.
// Возвращает все зависимости явно — тест контролирует их полностью.
// Не используй глобальный setup/teardown.

func newSigner(t *testing.T) (*Signer, *fake.Clientset) {
	t.Helper()
	cl := fake.NewSimpleClientset()
	inf := informers.NewSharedInformerFactory(cl, 0)

	signer, err := NewSigner(
		cl,
		inf.Core().V1().Secrets(),
		inf.Core().V1().ConfigMaps(),
		DefaultSignerOptions(),
	)
	if err != nil {
		t.Fatalf("error creating Signer: %v", err)
	}
	return signer, cl
}

// ── Test fixtures — модификаторы объектов ────────────────────────────────────
// Конструктор создаёт базовый валидный объект.
// Отдельные функции добавляют опциональные поля — не перегружают конструктор.

func newTokenSecret(tokenID, tokenSecret string) *v1.Secret {
	return &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:       "kube-system",
			Name:            "bootstrap-token-" + tokenID,
			ResourceVersion: "1",
		},
		Type: "bootstrap.kubernetes.io/token",
		Data: map[string][]byte{
			"token-id":     []byte(tokenID),
			"token-secret": []byte(tokenSecret),
		},
	}
}

func addSecretSigningUsage(s *v1.Secret, value string) {
	s.Data["usage-bootstrap-signing"] = []byte(value)
}

func addSecretExpiration(s *v1.Secret, expiration string) {
	s.Data["expiration"] = []byte(expiration)
}

func newConfigMap(tokenID, signature string) *v1.ConfigMap {
	cm := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "kube-public",
			Name:      "cluster-info",
		},
		Data: map[string]string{},
	}
	if tokenID != "" {
		cm.Data["jws-kubeconfig-"+tokenID] = signature
	}
	return cm
}

// ── TestSimpleSign ────────────────────────────────────────────────────────────
// Setup → Act → Assert — три чёткие фазы.
// Объекты добавляются в indexer напрямую, не через fake API —
// это быстрее и не создаёт лишних actions в cl.Actions().

func TestSimpleSign(t *testing.T) {
	signer, cl := newSigner(t)
	inf := informers.NewSharedInformerFactory(cl, 0)

	// SETUP: добавляем объекты в кэш информера, минуя API-сервер
	cm := newConfigMap("", "")
	inf.Core().V1().ConfigMaps().Informer().GetIndexer().Add(cm)

	secret := newTokenSecret("abc123", "secret456")
	addSecretSigningUsage(secret, "true")
	inf.Core().V1().Secrets().Informer().GetIndexer().Add(secret)

	// ACT
	_ = signer.signConfigMap(context.TODO())

	// ASSERT: проверяем что именно был сделан Update, не что-то другое
	expected := []core.Action{
		core.NewUpdateAction(
			schema.GroupVersionResource{Version: "v1", Resource: "configmaps"},
			"kube-public",
			newConfigMap("abc123", "expected-signature"),
		),
	}
	verifyActions(t, expected, cl.Actions())
}

// ── verifyActions ─────────────────────────────────────────────────────────────
// Вспомогательная функция для сравнения списков actions.
// t.Helper() — ошибка показывается в вызывающем тесте, не здесь.
// Проверяет как лишние, так и отсутствующие actions.

func verifyActions(t *testing.T, expected, actual []core.Action) {
	t.Helper()

	for i, a := range actual {
		if len(expected) < i+1 {
			t.Errorf("%d unexpected actions:\n%v", len(actual)-len(expected), actual[i:])
			break
		}
		// Используем Semantic.DeepEqual а не reflect.DeepEqual —
		// он понимает семантические эквиваленты k8s-объектов
		if e := expected[i]; e.GetVerb() != a.GetVerb() ||
			e.GetResource() != a.GetResource() ||
			e.GetNamespace() != a.GetNamespace() {
			t.Errorf("action[%d]:\n  expected: %v\n  got:      %v", i, e, a)
		}
	}

	if len(expected) > len(actual) {
		t.Errorf("%d expected actions not executed:", len(expected)-len(actual))
		for _, a := range expected[len(actual):] {
			t.Logf("  missing: %v", a)
		}
	}
}

// ── PrependReactor — перехват side-effects ───────────────────────────────────
// Используй когда fake-клиент должен выполнять side-effect при операции
// (например, обновлять indexer при create, возвращать ошибку при update).

func TestWithReactor(t *testing.T) {
	cl := fake.NewSimpleClientset()
	inf := informers.NewSharedInformerFactory(cl, 0)
	secretIndexer := inf.Core().V1().Secrets().Informer().GetIndexer()

	// При создании Secret — автоматически добавляем его в indexer
	cl.PrependReactor("create", "secrets", func(action core.Action) (bool, interface{}, error) {
		obj := action.(core.CreateAction).GetObject().(*v1.Secret)
		_ = secretIndexer.Add(obj)
		return false, obj, nil // false = продолжить стандартную обработку
	})

	// При попытке update — возвращаем ошибку (immutable ресурс)
	cl.PrependReactor("update", "secrets", func(action core.Action) (bool, interface{}, error) {
		return true, nil, names.ErrPrefixMissing // true = остановить обработку
	})

	_, _ = cl, secretIndexer
}

// ── Проверка что informer синхронизирован перед тестом ───────────────────────

func waitForCacheSync(t *testing.T, synced ...cache.InformerSynced) {
	t.Helper()
	stopCh := make(chan struct{})
	defer close(stopCh)
	if !cache.WaitForCacheSync(stopCh, synced...) {
		t.Fatal("timed out waiting for caches to sync")
	}
}
