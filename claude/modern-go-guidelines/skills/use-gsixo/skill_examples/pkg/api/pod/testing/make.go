// Example: pkg/api/pod/testing/make.go
// Patterns: Tweak functional option pattern for test fixture construction.
// Позволяет собирать тестовые объекты декларативно без перегруженного конструктора.
package testing

import (
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

// ── Tweak — функциональная опция ─────────────────────────────────────────────
// type Tweak func(*T) — стандартный паттерн для построения тестовых объектов.
// Преимущества перед конструктором с параметрами:
//   - добавление новой опции не ломает существующие вызовы
//   - тест читается как список намерений, не как список аргументов
//   - легко комбинируется: MakePod("name", SetNS("ns"), SetNode("n"), SetPhase(Running))

type Tweak func(*v1.Pod)
type TweakContainer func(*v1.Container)

// ── Базовый конструктор ───────────────────────────────────────────────────────
// Создаёт минимально валидный объект, затем применяет все tweaks.

func MakePod(name string, tweaks ...Tweak) *v1.Pod {
	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: metav1.NamespaceDefault,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				MakeContainer("ctr"),
			},
			DNSPolicy:                     v1.DNSClusterFirst,
			RestartPolicy:                 v1.RestartPolicyAlways,
			TerminationGracePeriodSeconds: ptr.To[int64](30),
		},
	}
	for _, tweak := range tweaks {
		tweak(pod)
	}
	return pod
}

func MakeContainer(name string, tweaks ...TweakContainer) v1.Container {
	c := v1.Container{
		Name:  name,
		Image: "image:latest",
	}
	for _, t := range tweaks {
		t(&c)
	}
	return c
}

// ── Tweak-функции ─────────────────────────────────────────────────────────────
// Каждая функция отвечает за одно поле/аспект.
// Именование: Set* для простых полей, With* для сложных добавлений.

func SetNamespace(ns string) Tweak {
	return func(pod *v1.Pod) { pod.Namespace = ns }
}

func SetNodeName(name string) Tweak {
	return func(pod *v1.Pod) { pod.Spec.NodeName = name }
}

func SetPhase(phase v1.PodPhase) Tweak {
	return func(pod *v1.Pod) { pod.Status.Phase = phase }
}

func SetLabels(labels map[string]string) Tweak {
	return func(pod *v1.Pod) { pod.Labels = labels }
}

func SetAnnotations(ann map[string]string) Tweak {
	return func(pod *v1.Pod) { pod.Annotations = ann }
}

func SetDeletionTimestamp(t time.Time) Tweak {
	return func(pod *v1.Pod) {
		mt := metav1.NewTime(t)
		pod.DeletionTimestamp = &mt
	}
}

func SetRestartPolicy(policy v1.RestartPolicy) Tweak {
	return func(pod *v1.Pod) { pod.Spec.RestartPolicy = policy }
}

func WithContainer(tweak TweakContainer) Tweak {
	return func(pod *v1.Pod) {
		if len(pod.Spec.Containers) == 0 {
			pod.Spec.Containers = append(pod.Spec.Containers, MakeContainer("ctr"))
		}
		tweak(&pod.Spec.Containers[0])
	}
}

func SetContainerImage(image string) TweakContainer {
	return func(c *v1.Container) { c.Image = image }
}

func SetContainerResources(req, lim v1.ResourceList) TweakContainer {
	return func(c *v1.Container) {
		c.Resources = v1.ResourceRequirements{
			Requests: req,
			Limits:   lim,
		}
	}
}

// ── Использование ────────────────────────────────────────────────────────────
// Пример из теста:
//
//   pod := MakePod("test-pod",
//       SetNamespace("kube-system"),
//       SetNodeName("node-1"),
//       SetPhase(v1.PodRunning),
//       WithContainer(SetContainerImage("nginx:1.21")),
//   )
//
// Намного читаемее чем:
//
//   pod := newPod("test-pod", "kube-system", "node-1", v1.PodRunning, "nginx:1.21", nil, nil)

// ── Вспомогательные функции для ресурсов ─────────────────────────────────────

func ResourceList(cpu, memory string) v1.ResourceList {
	return v1.ResourceList{
		v1.ResourceCPU:    resource.MustParse(cpu),
		v1.ResourceMemory: resource.MustParse(memory),
	}
}
