// Example: staging/src/k8s.io/api/apps/v1/types.go
// Patterns: TypeMeta→ObjectMeta→Spec→Status, typed enums (+enum),
//           +optional / +featureGate markers, *T for optional fields,
//           codegen markers (+genclient, +k8s:deepcopy-gen)
package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// ── Codegen markers ───────────────────────────────────────────────────────────
// Эти комментарии читаются генераторами, не компилятором.
// +genclient                         — сгенерировать typed client
// +k8s:deepcopy-gen:interfaces=...   — сгенерировать DeepCopyObject()
// +k8s:prerelease-lifecycle-gen      — сгенерировать lifecycle методы

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:prerelease-lifecycle-gen:introduced=1.9

// StatefulSet represents a set of pods with consistent identities.
// Структура API-объекта: TypeMeta → ObjectMeta → Spec → Status.
// Это единственный допустимый порядок полей для ресурсов Kubernetes.
type StatefulSet struct {
	metav1.TypeMeta `json:",inline"` // GroupVersionKind — встраивается

	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the desired identities of pods in this set.
	// +optional
	Spec StatefulSetSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`

	// Status is the current status of Pods in this StatefulSet.
	// This data may be out of date by some window of time.
	// +optional
	Status StatefulSetStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// ── Typed enum ────────────────────────────────────────────────────────────────
// type T string + const блок — стандарт для перечислений в Kubernetes API.
// +enum — маркер для генерации validation schema.
// Имя константы = имя типа + значение (PascalCase).

// PodManagementPolicyType defines the policy for creating pods under a stateful set.
// +enum
type PodManagementPolicyType string

const (
	// OrderedReadyPodManagement will create pods in strictly increasing order on
	// scale up and strictly decreasing order on scale down, progressing only when
	// the previous pod is ready or terminated.
	OrderedReadyPodManagement PodManagementPolicyType = "OrderedReady"

	// ParallelPodManagement will create and delete pods as soon as the stateful set
	// replica count is changed, and will not wait for pods to be ready or complete
	// termination.
	ParallelPodManagement PodManagementPolicyType = "Parallel"
)

// StatefulSetUpdateStrategyType is a string enumeration type that enumerates
// all possible update strategies for the StatefulSet controller.
// +enum
type StatefulSetUpdateStrategyType string

const (
	RollingUpdateStatefulSetStrategyType StatefulSetUpdateStrategyType = "RollingUpdate"
	OnDeleteStatefulSetStrategyType      StatefulSetUpdateStrategyType = "OnDelete"
)

// ── Spec с опциональными полями ───────────────────────────────────────────────
// *int32 вместо int32 — позволяет отличить 0 от "не задано".
// +optional + omitempty — согласованно на уровне комментария и JSON-тега.
// +featureGate=Xxx — поле доступно только при включённой фиче.

// StatefulSetSpec is the specification of a StatefulSet.
type StatefulSetSpec struct {
	// replicas is the desired number of replicas of the given Template.
	// If unspecified, defaults to 1.
	// +optional
	Replicas *int32 `json:"replicas,omitempty" protobuf:"varint,1,opt,name=replicas"`

	// updateStrategy indicates the StatefulSetUpdateStrategy that will be
	// employed to update Pods in the StatefulSet when a revision is made to
	// Template.
	UpdateStrategy StatefulSetUpdateStrategy `json:"updateStrategy,omitempty" protobuf:"bytes,7,opt,name=updateStrategy"`

	// revisionHistoryLimit is the maximum number of revisions that will be
	// maintained in the StatefulSet's revision history.
	// If not set, default is 10.
	// +optional
	RevisionHistoryLimit *int32 `json:"revisionHistoryLimit,omitempty" protobuf:"varint,8,opt,name=revisionHistoryLimit"`

	// ordinals controls the numbering of replica indices in a StatefulSet.
	// +featureGate=StatefulSetAutoDeletePVC
	// +optional
	Ordinals *StatefulSetOrdinals `json:"ordinals,omitempty" protobuf:"bytes,11,opt,name=ordinals"`
}

// StatefulSetUpdateStrategy indicates the strategy that the StatefulSet
// controller will use to perform updates.
type StatefulSetUpdateStrategy struct {
	// Type indicates the type of the StatefulSetUpdateStrategy.
	// +optional
	Type StatefulSetUpdateStrategyType `json:"type,omitempty" protobuf:"bytes,1,opt,name=type"`

	// RollingUpdate is used to communicate parameters when Type is RollingUpdateStatefulSetStrategyType.
	// +optional
	RollingUpdate *RollingUpdateStatefulSetStrategy `json:"rollingUpdate,omitempty" protobuf:"bytes,2,opt,name=rollingUpdate"`
}

type RollingUpdateStatefulSetStrategy struct {
	// Partition indicates the ordinal at which the StatefulSet should be partitioned
	// for updates. Default is 0.
	// +optional
	Partition *int32 `json:"partition,omitempty" protobuf:"varint,1,opt,name=partition"`

	// maxUnavailable is the maximum number of pods that can be unavailable during the update.
	// +featureGate=MaxUnavailableStatefulSet
	// +optional
	MaxUnavailable *intstr.IntOrString `json:"maxUnavailable,omitempty" protobuf:"bytes,2,opt,name=maxUnavailable"`
}

type StatefulSetOrdinals struct {
	// start is the number representing the first replica's index.
	// +optional
	Start int32 `json:"start" protobuf:"varint,1,opt,name=start"`
}

type StatefulSetStatus struct {
	// observedGeneration is the most recent generation observed.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	Replicas           int32 `json:"replicas"`
	ReadyReplicas      int32 `json:"readyReplicas,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// StatefulSetList is a collection of StatefulSets.
type StatefulSetList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []StatefulSet `json:"items"`
}
