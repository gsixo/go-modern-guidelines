// Example: pkg/controller/cronjob/cronjob_controllerv2.go
// Patterns: controller struct layout, Run(), processNextWorkItem(), sync/reconcile,
//           klog.FromContext, DeepCopy before mutation, switch{} for errors
package cronjob

import (
	"context"
	"fmt"
	"sync"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/wait"
	batchv1listers "k8s.io/client-go/listers/batch/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/record"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/klog/v2"
)

// ── Struct layout ─────────────────────────────────────────────────────────────
// Порядок полей: queue → clients → control interfaces → listers → synced → helpers.
// Интерфейсы вместо конкретных типов для тестируемости.
// now — инъекция времени, подменяется в тестах.

type ControllerV2 struct {
	queue workqueue.TypedRateLimitingInterface[string] // 1. очередь

	kubeClient  kubernetes.Interface  // 2. clients
	recorder    record.EventRecorder
	broadcaster record.EventBroadcaster

	jobControl     jobControlInterface  // 3. control interfaces (инъекция для тестов)
	cronJobControl cjControlInterface

	jobLister     batchv1listers.JobLister      // 4. listers
	cronJobLister batchv1listers.CronJobLister

	jobListerSynced     cache.InformerSynced // 5. synced checks
	cronJobListerSynced cache.InformerSynced

	now func() time.Time // 6. тест-хелпер; в продакшне: time.Now
}

// ── Run ───────────────────────────────────────────────────────────────────────
// Шаблон: defer HandleCrash → defer ShutDown+Wait → cache sync → воркеры → ctx.Done()
// Этот порядок гарантирует корректное завершение при отмене контекста.

func (jm *ControllerV2) Run(ctx context.Context, workers int) {
	defer utilruntime.HandleCrash() // панику превращаем в лог, не падаем

	logger := klog.FromContext(ctx)
	logger.Info("Starting CronJob controller")

	var wg sync.WaitGroup
	defer func() {
		logger.Info("Shutting down CronJob controller")
		jm.queue.ShutDown()  // сначала останавливаем очередь...
		wg.Wait()            // ...потом ждём завершения воркеров
	}()

	// Не стартуем до синхронизации кэша информеров
	if !cache.WaitForNamedCacheSync("cronjob", ctx.Done(),
		jm.jobListerSynced, jm.cronJobListerSynced) {
		return
	}

	for i := 0; i < workers; i++ {
		wg.Go(func() {
			wait.UntilWithContext(ctx, jm.worker, time.Second)
		})
	}

	<-ctx.Done() // блокируемся до отмены контекста
}

// ── Worker loop ───────────────────────────────────────────────────────────────

func (jm *ControllerV2) worker(ctx context.Context) {
	for jm.processNextWorkItem(ctx) {
	}
}

// ── processNextWorkItem ───────────────────────────────────────────────────────
// Стандартный паттерн workqueue:
//   Get → defer Done → sync → Forget/AddRateLimited/AddAfter
// switch{} вместо if/else if — специфичные случаи явно разделены.

func (jm *ControllerV2) processNextWorkItem(ctx context.Context) bool {
	key, quit := jm.queue.Get()
	if quit {
		return false
	}
	defer jm.queue.Done(key)

	requeueAfter, err := jm.sync(ctx, key)
	switch {
	case err != nil:
		// Логируем через HandleError (не fmt.Println!) и ставим в rate-limited очередь
		utilruntime.HandleError(fmt.Errorf("error syncing CronJob %v, requeuing: %w", key, err))
		jm.queue.AddRateLimited(key) // экспоненциальный backoff встроен в queue

	case requeueAfter != nil:
		jm.queue.Forget(key)                    // сбрасываем счётчик rate-limit
		jm.queue.AddAfter(key, *requeueAfter)   // запланированный requeue

	default:
		jm.queue.Forget(key)
	}
	return true
}

// ── Sync / Reconcile ──────────────────────────────────────────────────────────
// Стандартная структура: split key → get from lister → IsNotFound → DeepCopy → reconcile.
// IsNotFound — не ошибка, объект был удалён.
// DeepCopy() обязателен: lister отдаёт указатель в shared кэш, мутировать его нельзя.

func (jm *ControllerV2) sync(ctx context.Context, cronJobKey string) (*time.Duration, error) {
	ns, name, err := cache.SplitMetaNamespaceKey(cronJobKey)
	if err != nil {
		return nil, err
	}

	logger := klog.FromContext(ctx)

	cronJob, err := jm.cronJobLister.CronJobs(ns).Get(name)
	switch {
	case apierrors.IsNotFound(err):
		// Объект удалён — это нормально, ничего делать не надо
		logger.V(4).Info("CronJob not found, may have been deleted", "key", cronJobKey)
		return nil, nil
	case err != nil:
		return nil, err
	}

	// Всегда работаем с копией из кэша!
	// cronJobCopy объединяет все обновления для одного API-вызова в конце
	cronJobCopy := cronJob.DeepCopy()

	jobs, err := jm.getJobsToBeReconciled(cronJob)
	if err != nil {
		return nil, err
	}

	requeueAfter, updateStatus, syncErr := jm.syncCronJob(ctx, cronJobCopy, jobs)
	_ = updateStatus // используется для определения нужно ли вызывать UpdateStatus
	return requeueAfter, syncErr
}

// заглушки для компиляции примера
func (jm *ControllerV2) getJobsToBeReconciled(_ *batchv1.CronJob) ([]*batchv1.Job, error) {
	return nil, nil
}
func (jm *ControllerV2) syncCronJob(_ context.Context, _ *batchv1.CronJob, _ []*batchv1.Job) (*time.Duration, bool, error) {
	return nil, false, nil
}
