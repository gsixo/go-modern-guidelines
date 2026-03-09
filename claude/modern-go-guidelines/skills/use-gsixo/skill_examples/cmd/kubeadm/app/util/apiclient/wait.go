// Example: cmd/kubeadm/app/util/apiclient/wait.go
// Patterns: wait.PollUntilContextTimeout, wait.UntilWithContext,
//           ktesting.NewTestContext, structured wait patterns
package apiclient

import (
	"context"
	"fmt"
	"time"

	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
)

// ── wait.PollUntilContextTimeout ─────────────────────────────────────────────
// Было: ручной ticker + select + time.After + break-логика
// Стало: декларативный polling с автоматической отменой.
//
// Сигнатура: PollUntilContextTimeout(ctx, interval, timeout, immediate, condition)
//   - immediate=true: condition вызывается сразу, потом через interval
//   - immediate=false: первый вызов через interval
//
// condition возвращает (done bool, err error):
//   - (true, nil)  → успех, выходим
//   - (false, nil) → не готово, ждём следующего интервала
//   - (_, err)     → немедленный выход с ошибкой

func WaitForAPIServer(ctx context.Context, client kubernetes.Interface, timeout time.Duration) error {
	err := wait.PollUntilContextTimeout(ctx, time.Second, timeout, true,
		func(ctx context.Context) (bool, error) {
			// Лёгкий запрос для проверки живости API сервера
			_, err := client.Discovery().ServerVersion()
			if err != nil {
				return false, nil // не ошибка — просто ещё не готов
			}
			return true, nil
		},
	)
	if err != nil {
		return fmt.Errorf("timeout waiting for API server: %w", err)
	}
	return nil
}

// WaitForPodsWithLabel ждёт пока все поды с меткой не станут Ready.
func WaitForPodsWithLabel(ctx context.Context, client kubernetes.Interface, ns, label string) error {
	return wait.PollUntilContextTimeout(ctx, 2*time.Second, 5*time.Minute, false,
		func(ctx context.Context) (bool, error) {
			pods, err := client.CoreV1().Pods(ns).List(ctx, metav1.ListOptions{
				LabelSelector: label,
			})
			if err != nil {
				if apierrors.IsNotFound(err) {
					return false, nil
				}
				return false, err // реальная ошибка — прерываем polling
			}
			if len(pods.Items) == 0 {
				return false, nil
			}
			for _, pod := range pods.Items {
				if !isPodReady(&pod) {
					return false, nil // хотя бы один не готов — ждём
				}
			}
			return true, nil
		},
	)
}

func isPodReady(pod *v1.Pod) bool {
	for _, cond := range pod.Status.Conditions {
		if cond.Type == v1.PodReady && cond.Status == v1.ConditionTrue {
			return true
		}
	}
	return false
}

// ── wait.UntilWithContext ─────────────────────────────────────────────────────
// Используется в controller.Run() для запуска worker-горутин.
// period=0: функция вызывается немедленно после возврата (tight loop).
// period>0: ждём period между вызовами.
//
// ВАЖНО: fn должна проверять ctx.Done() внутри если долго работает.

func startWorker(ctx context.Context, workerFn func(ctx context.Context)) {
	// period=0 → запускается снова сразу после возврата
	wait.UntilWithContext(ctx, workerFn, 0)
	// Когда ctx отменяется — UntilWithContext возвращается
}

// ── Backoff с ExponentialWithJitter ──────────────────────────────────────────

var defaultWaitBackoff = wait.Backoff{
	Duration: 500 * time.Millisecond,
	Factor:   1.5,
	Jitter:   0.5,
	Steps:    10,
	Cap:      30 * time.Second,
}

func retryWithBackoff(fn func() error) error {
	return wait.ExponentialBackoffWithContext(context.Background(), defaultWaitBackoff,
		func(ctx context.Context) (bool, error) {
			if err := fn(); err != nil {
				return false, nil // retry
			}
			return true, nil
		},
	)
}
