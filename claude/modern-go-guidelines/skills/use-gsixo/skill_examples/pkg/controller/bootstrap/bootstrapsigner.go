// Example: pkg/controller/bootstrap/bootstrapsigner.go
// Patterns: SignerOptions + DefaultSignerOptions(), FilteringResourceEventHandler,
//           Run() lifecycle pattern, cache sync, graceful shutdown
package bootstrap

import (
	"context"
	"sync"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	coreinformers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	corelisters "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
)

// ── Options struct + DefaultOptions() ────────────────────────────────────────
// Группируем конфигурацию в struct — не перечисляем параметры в New().
// DefaultSignerOptions() даёт разумные дефолты; вызывающий меняет только нужное.

type SignerOptions struct {
	ConfigMapNamespace   string
	ConfigMapName        string
	TokenSecretNamespace string
	ConfigMapResync      time.Duration
	SecretResync         time.Duration
}

func DefaultSignerOptions() SignerOptions {
	return SignerOptions{
		ConfigMapNamespace:   "kube-public",
		ConfigMapName:        "cluster-info",
		TokenSecretNamespace: "kube-system",
	}
}

// ── Struct layout ─────────────────────────────────────────────────────────────
// clients → config → queue → listers → synced checks

type Signer struct {
	client             kubernetes.Interface // 1. client
	configMapNamespace string               // 2. config
	configMapName      string
	secretNamespace    string

	syncQueue workqueue.TypedRateLimitingInterface[string] // 3. queue

	secretLister    corelisters.SecretLister      // 4. listers
	secretSynced    cache.InformerSynced
	configMapLister corelisters.ConfigMapLister
	configMapSynced cache.InformerSynced
}

func NewSigner(
	cl kubernetes.Interface,
	secrets coreinformers.SecretInformer,
	configMaps coreinformers.ConfigMapInformer,
	options SignerOptions,
) (*Signer, error) {
	e := &Signer{
		client:             cl,
		configMapNamespace: options.ConfigMapNamespace,
		configMapName:      options.ConfigMapName,
		secretNamespace:    options.TokenSecretNamespace,
		syncQueue:          workqueue.NewTypedRateLimitingQueue[string](workqueue.DefaultTypedControllerRateLimiter[string]()),
		secretLister:       secrets.Lister(),
		secretSynced:       secrets.Informer().HasSynced,
		configMapLister:    configMaps.Lister(),
		configMapSynced:    configMaps.Informer().HasSynced,
	}

	// ── FilteringResourceEventHandler ────────────────────────────────────────
	// Регистрируем только события для конкретного ConfigMap.
	// FilterFunc отсекает шум — воркер не просыпается по чужим событиям.
	configMaps.Informer().AddEventHandlerWithResyncPeriod(
		cache.FilteringResourceEventHandler{
			FilterFunc: func(obj interface{}) bool {
				switch t := obj.(type) {
				case *v1.ConfigMap:
					return t.Name == options.ConfigMapName &&
						t.Namespace == options.ConfigMapNamespace
				default:
					utilruntime.HandleError(nil)
					return false
				}
			},
			Handler: cache.ResourceEventHandlerFuncs{
				AddFunc:    func(_ interface{}) { e.pokeConfigMapSync() },
				UpdateFunc: func(_, _ interface{}) { e.pokeConfigMapSync() },
			},
		},
		options.ConfigMapResync,
	)

	return e, nil
}

func (e *Signer) pokeConfigMapSync() {
	e.syncQueue.Add(e.configMapNamespace + "/" + e.configMapName)
}

// ── Run() — стандартный lifecycle ─────────────────────────────────────────────
// 1. defer HandleCrash    — паника не роняет процесс
// 2. defer ShutDown+Wait  — корректное завершение воркеров
// 3. WaitForCacheSync     — не стартуем до синхронизации кэша
// 4. wg.Go(worker)        — запускаем воркеры
// 5. <-ctx.Done()         — блокируемся до отмены контекста

func (e *Signer) Run(ctx context.Context) {
	defer utilruntime.HandleCrash()

	logger := klog.FromContext(ctx)
	logger.V(5).Info("Starting bootstrap signer")

	var wg sync.WaitGroup
	defer func() {
		logger.V(1).Info("Shutting down bootstrap signer")
		e.syncQueue.ShutDown() // сначала очередь
		wg.Wait()              // потом ждём воркеров
	}()

	if !cache.WaitForNamedCacheSync("bootstrap-signer", ctx.Done(),
		e.configMapSynced, e.secretSynced) {
		return // контекст отменён во время ожидания
	}

	wg.Go(func() {
		wait.UntilWithContext(ctx, e.serviceConfigMapQueue, 0)
	})

	<-ctx.Done()
}

func (e *Signer) serviceConfigMapQueue(ctx context.Context) {
	key, quit := e.syncQueue.Get()
	if quit {
		return
	}
	defer e.syncQueue.Done(key)

	if err := e.signConfigMap(ctx); err != nil {
		utilruntime.HandleError(err)
		e.syncQueue.AddRateLimited(key)
		return
	}
	e.syncQueue.Forget(key)
}

func (e *Signer) signConfigMap(_ context.Context) error {
	// реализация подписи ConfigMap
	return nil
}
