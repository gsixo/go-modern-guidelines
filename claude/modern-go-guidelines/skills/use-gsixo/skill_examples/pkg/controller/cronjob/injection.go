// Example: pkg/controller/cronjob/injection.go
// Patterns: small focused interfaces, real vs fake implementations,
//           compile-time interface check, context.TODO() anti-pattern
package cronjob

import (
	"context"
	"sync"

	batchv1 "k8s.io/api/batch/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ── Малый интерфейс — одна ответственность ───────────────────────────────────
// Интерфейс определяется у потребителя (контроллера), не у производителя.
// Каждый интерфейс закрывает ровно одну задачу.

type cjControlInterface interface {
	UpdateStatus(ctx context.Context, cj *batchv1.CronJob) (*batchv1.CronJob, error)
	GetCronJob(ctx context.Context, namespace, name string) (*batchv1.CronJob, error)
}

type jobControlInterface interface {
	GetJob(namespace, name string) (*batchv1.Job, error)
	CreateJob(namespace string, job *batchv1.Job) (*batchv1.Job, error)
	DeleteJob(namespace string, name string) error
}

// ── Real implementation ───────────────────────────────────────────────────────
// Оборачивает реальный kubeClient. Поля — только необходимые зависимости.

type realCJControl struct {
	KubeClient kubernetes.Interface
}

// Compile-time check: убеждаемся что realCJControl удовлетворяет интерфейсу.
// Если нет — ошибка компиляции, не runtime-паника.
var _ cjControlInterface = &realCJControl{}

func (c *realCJControl) GetCronJob(ctx context.Context, namespace, name string) (*batchv1.CronJob, error) {
	return c.KubeClient.BatchV1().CronJobs(namespace).Get(ctx, name, metav1.GetOptions{})
}

func (c *realCJControl) UpdateStatus(ctx context.Context, cj *batchv1.CronJob) (*batchv1.CronJob, error) {
	return c.KubeClient.BatchV1().CronJobs(cj.Namespace).UpdateStatus(ctx, cj, metav1.UpdateOptions{})
}

// ── Fake implementation (для тестов) ─────────────────────────────────────────
// Хранит тестовые данные. Exported поля — чтобы тест мог их читать/проверять.

type fakeCJControl struct {
	CronJob *batchv1.CronJob
	Updates []batchv1.CronJob
}

var _ cjControlInterface = &fakeCJControl{}

func (c *fakeCJControl) GetCronJob(_ context.Context, namespace, name string) (*batchv1.CronJob, error) {
	if c.CronJob != nil && name == c.CronJob.Name && namespace == c.CronJob.Namespace {
		return c.CronJob, nil
	}
	return nil, apierrors.NewNotFound(
		schema.GroupResource{Group: "batch", Resource: "cronjobs"}, name,
	)
}

func (c *fakeCJControl) UpdateStatus(_ context.Context, cj *batchv1.CronJob) (*batchv1.CronJob, error) {
	c.Updates = append(c.Updates, *cj) // записываем для последующей проверки в тесте
	return cj, nil
}

// ── Fake job control с mutex ──────────────────────────────────────────────────
// sync.Mutex встроен первым полем — стандарт для thread-safe test doubles.
// Exported поля позволяют тесту проверять что именно было создано/удалено.

type fakeJobControl struct {
	sync.Mutex                   // первое поле — сразу виден как thread-safe
	Job           *batchv1.Job
	Jobs          []batchv1.Job
	DeleteJobName []string
	CreateErr     error
	Err           error
}

var _ jobControlInterface = &fakeJobControl{}

func (f *fakeJobControl) GetJob(namespace, name string) (*batchv1.Job, error) {
	f.Lock()
	defer f.Unlock()
	if f.Err != nil {
		return nil, f.Err
	}
	return f.Job, nil
}

func (f *fakeJobControl) CreateJob(namespace string, job *batchv1.Job) (*batchv1.Job, error) {
	f.Lock()
	defer f.Unlock()
	if f.CreateErr != nil {
		return nil, f.CreateErr
	}
	f.Jobs = append(f.Jobs, *job)
	job.UID = "test-uid"
	return job, nil
}

func (f *fakeJobControl) DeleteJob(namespace, name string) error {
	f.Lock()
	defer f.Unlock()
	if f.Err != nil {
		return f.Err
	}
	f.DeleteJobName = append(f.DeleteJobName, name)
	return nil
}

// ── Real job control — ANTI-PATTERN: context.TODO() ──────────────────────────
// В реальном коде kubernetes/pkg/controller/cronjob/injection.go
// методы jobControlInterface не принимают context — это legacy.
// Правильно: добавить ctx context.Context первым параметром во все методы.

type realJobControl struct {
	KubeClient kubernetes.Interface
}

var _ jobControlInterface = &realJobControl{}

func (r *realJobControl) GetJob(namespace, name string) (*batchv1.Job, error) {
	// ❌ context.TODO() — теряем timeout и cancellation
	// ✓ должно быть: func GetJob(ctx context.Context, ...) с переданным ctx
	return r.KubeClient.BatchV1().Jobs(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func (r *realJobControl) CreateJob(namespace string, job *batchv1.Job) (*batchv1.Job, error) {
	return r.KubeClient.BatchV1().Jobs(namespace).Create(context.TODO(), job, metav1.CreateOptions{})
}

func (r *realJobControl) DeleteJob(namespace, name string) error {
	return r.KubeClient.BatchV1().Jobs(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{
		PropagationPolicy: func() *metav1.DeletionPropagation {
			p := metav1.DeletePropagationBackground
			return &p
		}(),
	})
}
