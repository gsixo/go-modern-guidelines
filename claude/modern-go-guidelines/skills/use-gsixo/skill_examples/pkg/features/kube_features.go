// Example: pkg/features/kube_features.go
// Patterns: feature gate constants, owner/KEP comments, alphabetical order,
//           usage at call sites with utilfeature.DefaultFeatureGate.Enabled()
package features

import "k8s.io/component-base/featuregate"

// ── Feature gate константы ────────────────────────────────────────────────────
// Правила оформления (из реального kube_features.go):
//   1. Тип: featuregate.Feature (type alias для string — type-safety)
//   2. Комментарий: // owner: @github-username
//   3. Ссылка на KEP: // kep: https://kep.k8s.io/NNN (опционально)
//   4. Описание: что делает фича
//   5. Алфавитный порядок (с учётом регистра: заглавные перед строчными)
//
// Регистрация дефолтов — в отдельном вызове AddFeatureGates() при инициализации.

const (
	// owner: @aojea
	//
	// Allow kubelet to request a certificate without any Node IP available, only
	// with DNS names.
	AllowDNSOnlyNodeCSR featuregate.Feature = "AllowDNSOnlyNodeCSR"

	// owner: @micahhausler
	//
	// Setting AllowInsecureKubeletCertificateSigningRequests to true disables node
	// admission validation of CSRs for kubelet signers.
	AllowInsecureKubeletCertificateSigningRequests featuregate.Feature = "AllowInsecureKubeletCertificateSigningRequests"

	// owner: @deads2k
	// kep: https://kep.k8s.io/4460
	//
	// APIServerTracing enables tracing calls to the API server.
	APIServerTracing featuregate.Feature = "APIServerTracing"

	// owner: @jpbetz
	// kep: https://kep.k8s.io/4355
	//
	// MaxUnavailableStatefulSet enables the maxUnavailable field for
	// RollingUpdateStatefulSetStrategy.
	MaxUnavailableStatefulSet featuregate.Feature = "MaxUnavailableStatefulSet"
)

// ── Использование feature gate в коде ────────────────────────────────────────
// Импортируй утилиту и проверяй через Enabled().
// Blank import для регистрации фич из других пакетов.

// В реальном коде:
//
//   import (
//       utilfeature "k8s.io/apiserver/pkg/util/feature"
//       _ "k8s.io/kubernetes/pkg/features" // регистрирует фичи через init()
//   )
//
//   func handleRequest(ctx context.Context, req *Request) {
//       if utilfeature.DefaultFeatureGate.Enabled(features.APIServerTracing) {
//           span := startTrace(ctx, "handleRequest")
//           defer span.End()
//       }
//       // ...
//   }
//
//   func createStatefulSet(spec *StatefulSetSpec) {
//       if utilfeature.DefaultFeatureGate.Enabled(features.MaxUnavailableStatefulSet) {
//           // применяем maxUnavailable из spec
//       }
//   }
