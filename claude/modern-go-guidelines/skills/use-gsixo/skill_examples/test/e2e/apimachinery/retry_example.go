// Example: test/e2e/apimachinery/validatingadmissionpolicy.go
// Patterns: retry.RetryOnConflict, DeepCopy before update, ResourceVersion handling
package apimachinery

import (
	"context"
	"fmt"

	admissionv1 "k8s.io/api/admissionregistration/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
)

// ── retry.RetryOnConflict ─────────────────────────────────────────────────────
// Стандартный паттерн для Update-операций в Kubernetes.
// Проблема: между Get и Update другой контроллер может изменить объект.
// Решение: при 409 Conflict — перечитать объект и повторить.
//
// retry.DefaultRetry = 5 попыток с экспоненциальным backoff (100ms base).
// retry.DefaultBackoff — то же самое но без jitter.

func updateVAPWithRetry(
	ctx context.Context,
	client kubernetes.Interface,
	name string,
	mutateFn func(*admissionv1.ValidatingAdmissionPolicy),
) (*admissionv1.ValidatingAdmissionPolicy, error) {
	var updated *admissionv1.ValidatingAdmissionPolicy

	err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		// 1. Читаем свежую версию (важно: ResourceVersion актуален)
		current, err := client.AdmissionregistrationV1().
			ValidatingAdmissionPolicies().
			Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return err // не-conflict ошибка — не ретраим
		}

		// 2. Мутируем КОПИЮ — никогда не мутируем объект из кэша
		toUpdate := current.DeepCopy()
		mutateFn(toUpdate)

		// 3. Update — может вернуть 409 Conflict
		updated, err = client.AdmissionregistrationV1().
			ValidatingAdmissionPolicies().
			Update(ctx, toUpdate, metav1.UpdateOptions{})
		return err // при 409 RetryOnConflict повторит цикл
	})

	return updated, err
}

// ── Пример использования ─────────────────────────────────────────────────────

func exampleVAPUpdate(ctx context.Context, client kubernetes.Interface) error {
	_, err := updateVAPWithRetry(ctx, client, "my-policy", func(vap *admissionv1.ValidatingAdmissionPolicy) {
		if vap.Annotations == nil {
			vap.Annotations = map[string]string{}
		}
		vap.Annotations["updated-by"] = "controller"

		fail := admissionv1.Fail
		vap.Spec.FailurePolicy = &fail
	})
	return err
}

// ── Custom backoff ────────────────────────────────────────────────────────────
// Для конкретных случаев можно задать свой backoff.
// Пример из pkg/controller/controller_utils.go:

// var UpdateTaintBackoff = wait.Backoff{
//     Steps:    5,
//     Duration: 100 * time.Millisecond,
//     Jitter:   1.0,
// }

// ── RetryOnConflict с первым чтением из кэша ─────────────────────────────────
// Оптимизация: первый Get читает из кэша (ResourceVersion="0"),
// последующие — с актуальным ResourceVersion напрямую с API-сервера.

func updateNodeWithCacheOptimization(
	ctx context.Context,
	client kubernetes.Interface,
	nodeName string,
	mutateFn func(*metav1.ObjectMeta),
) error {
	firstTry := true
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		opts := metav1.GetOptions{}
		if firstTry {
			opts.ResourceVersion = "0" // читаем из kube-apiserver кэша
			firstTry = false
		}
		// при следующих попытках opts.ResourceVersion пустой — идём напрямую

		node, err := client.CoreV1().Nodes().Get(ctx, nodeName, opts)
		if err != nil {
			return fmt.Errorf("get node %s: %w", nodeName, err)
		}

		nodeCopy := node.DeepCopy()
		mutateFn(&nodeCopy.ObjectMeta)

		_, err = client.CoreV1().Nodes().Update(ctx, nodeCopy, metav1.UpdateOptions{})
		return err
	})
}
