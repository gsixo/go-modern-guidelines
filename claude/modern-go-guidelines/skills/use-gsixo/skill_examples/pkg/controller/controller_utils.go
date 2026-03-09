// Example: pkg/controller/controller_utils.go
// Patterns: named Backoff vars, RetryOnConflict with custom backoff
package controller

import (
	"context"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	clientretry "k8s.io/client-go/util/retry"
	"k8s.io/client-go/util/wait"
)

// ── Именованные переменные для Backoff ───────────────────────────────────────
// Вместо использования retry.DefaultRetry, создавай специализированные backoff
// с понятными именами — это самодокументирует намерение.
//
// wait.Backoff.Jitter: каждый шаг умножается на (1 + Jitter*rand), устраняет thundering herd.

var UpdateTaintBackoff = wait.Backoff{
	Steps:    5,
	Duration: 100 * time.Millisecond,
	Jitter:   1.0,
}

var UpdateLabelBackoff = wait.Backoff{
	Steps:    5,
	Duration: 100 * time.Millisecond,
	Jitter:   1.0,
}

var NodeRequestBackoff = wait.Backoff{
	Steps:    10,
	Duration: 50 * time.Millisecond,
	Factor:   1.5,
	Jitter:   1.0,
}

// ── RetryOnConflict с custom backoff ─────────────────────────────────────────
// clientretry.RetryOnConflict повторяет только при 409 Conflict.
// При любой другой ошибке — немедленный возврат.

func AddOrUpdateTaintOnNode(
	ctx context.Context,
	c kubernetes.Interface,
	nodeName string,
	taints ...*v1.Taint,
) error {
	firstTry := true
	return clientretry.RetryOnConflict(UpdateTaintBackoff, func() error {
		opts := metav1.GetOptions{}
		if firstTry {
			opts.ResourceVersion = "0" // читаем из apiserver кэша при первой попытке
			firstTry = false
		}

		node, err := c.CoreV1().Nodes().Get(ctx, nodeName, opts)
		if err != nil {
			return err
		}

		nodeCopy := node.DeepCopy() // никогда не мутируем объект из кэша
		for _, taint := range taints {
			nodeCopy.Spec.Taints = appendOrReplaceTaint(nodeCopy.Spec.Taints, taint)
		}

		_, err = c.CoreV1().Nodes().Update(ctx, nodeCopy, metav1.UpdateOptions{})
		return err // 409 → RetryOnConflict повторит; другие ошибки — немедленный возврат
	})
}

// appendOrReplaceTaint добавляет или заменяет taint в списке.
func appendOrReplaceTaint(taints []v1.Taint, taint *v1.Taint) []v1.Taint {
	for i := range taints {
		if taints[i].Key == taint.Key && taints[i].Effect == taint.Effect {
			taints[i] = *taint
			return taints
		}
	}
	return append(taints, *taint)
}

// ── wait.PollUntilContextTimeout ─────────────────────────────────────────────
// Было: ручной ticker + таймаут + select
// Стало: декларативный polling с гарантированным завершением при ctx.Done()
//
// Параметры: ctx, interval, timeout, immediate bool, condition func
// immediate=true — сначала вызвать condition немедленно, потом по интервалу

func waitForNodeReady(ctx context.Context, c kubernetes.Interface, nodeName string) error {
	return wait.PollUntilContextTimeout(ctx, time.Second, 5*time.Minute, true,
		func(ctx context.Context) (bool, error) {
			node, err := c.CoreV1().Nodes().Get(ctx, nodeName, metav1.GetOptions{})
			if err != nil {
				return false, nil // не ошибка — просто не готово ещё
			}
			for _, cond := range node.Status.Conditions {
				if cond.Type == v1.NodeReady && cond.Status == v1.ConditionTrue {
					return true, nil // готово
				}
			}
			return false, nil
		},
	)
}
