// Example: pkg/apis/batch/validation/validation_test.go
// Patterns: map[string]struct{} table-driven tests, update func in test case,
//           cmp.Diff + cmpopts.IgnoreFields, ptr.To, t.Parallel inside t.Run
package validation_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	batch "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/ptr"
)

// ── Игнорируемые поля при сравнении ошибок ────────────────────────────────────
// BadValue и Detail часто нестабильны — игнорируем их в assertions.
// Origin зависит от версии API — тоже игнорируем.

var ignoreErrValueDetail = cmpopts.IgnoreFields(field.Error{}, "BadValue", "Detail", "Origin")

// ── getValid* — базовые конструкторы ─────────────────────────────────────────
// Возвращают минимально валидные объекты.
// Называем getValid* чтобы было ясно: это не любой объект, а валидный.

func getValidManualSelector() *metav1.LabelSelector {
	return &metav1.LabelSelector{
		MatchLabels: map[string]string{"a": "b"},
	}
}

func getValidJob() batch.Job {
	selector := getValidManualSelector()
	return batch.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-job",
			Namespace: metav1.NamespaceDefault,
		},
		Spec: batch.JobSpec{
			Selector:    selector,
			Parallelism: ptr.To[int32](1),
			Completions: ptr.To[int32](1),
			Template: batch.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: selector.MatchLabels},
			},
		},
	}
}

// ── Table-driven test: map[string]struct{} ────────────────────────────────────
// Ключ = название кейса (используется как t.Run имя).
// Порядок не важен — каждый кейс независим.
// Предпочтительно когда кейсов много и порядок исполнения не важен.

func TestValidateJob(t *testing.T) {
	successCases := map[string]struct {
		job batch.Job
	}{
		"valid basic job": {
			job: getValidJob(),
		},
		"valid job with parallelism 0": {
			job: func() batch.Job {
				j := getValidJob()
				j.Spec.Parallelism = ptr.To[int32](0)
				return j
			}(),
		},
		"nil parallelism defaults to 1": {
			job: func() batch.Job {
				j := getValidJob()
				j.Spec.Parallelism = nil
				return j
			}(),
		},
	}

	for name, tc := range successCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel() // независимые кейсы — запускаем параллельно
			// errs := ValidateJob(&tc.job, JobValidationOptions{})
			// if len(errs) != 0 {
			//     t.Errorf("unexpected errors: %v", errs)
			// }
			_ = tc
		})
	}

	errorCases := map[string]struct {
		job     batch.Job
		wantErr field.Error
	}{
		"missing selector": {
			job: func() batch.Job {
				j := getValidJob()
				j.Spec.Selector = nil
				return j
			}(),
			wantErr: field.Error{
				Type:  field.ErrorTypeRequired,
				Field: "spec.selector",
			},
		},
		"negative parallelism": {
			job: func() batch.Job {
				j := getValidJob()
				j.Spec.Parallelism = ptr.To[int32](-1)
				return j
			}(),
			wantErr: field.Error{
				Type:  field.ErrorTypeInvalid,
				Field: "spec.parallelism",
			},
		},
	}

	for name, tc := range errorCases {
		t.Run(name, func(t *testing.T) {
			// errs := ValidateJob(&tc.job, JobValidationOptions{})
			// wantErrs := field.ErrorList{&tc.wantErr}
			// if diff := cmp.Diff(wantErrs, errs, ignoreErrValueDetail); diff != "" {
			//     t.Errorf("unexpected errors (-want,+got):\n%s", diff)
			// }
			_ = tc
		})
	}
}

// ── Table-driven test: update func в кейсе ────────────────────────────────────
// Используется для тестирования ValidateXxxUpdate.
// old — базовый объект до обновления.
// update — функция-мутатор: применяет только изменения конкретного кейса.
// Избегаем дублирования: не копируем весь объект для каждого кейса.

func TestValidateJobUpdate(t *testing.T) {
	cases := map[string]struct {
		old    batch.Job
		update func(*batch.Job)
		err    *field.Error // nil = ожидаем успех
	}{
		"mutable: parallelism": {
			old: getValidJob(),
			update: func(job *batch.Job) {
				job.Spec.Parallelism = ptr.To[int32](5) // разрешено менять
			},
		},
		"immutable: selector": {
			old: getValidJob(),
			update: func(job *batch.Job) {
				job.Spec.Selector = &metav1.LabelSelector{
					MatchLabels: map[string]string{"new": "selector"},
				}
			},
			err: &field.Error{
				Type:  field.ErrorTypeInvalid,
				Field: "spec.selector",
			},
		},
		"immutable: completions when indexed": {
			old: func() batch.Job {
				j := getValidJob()
				mode := batch.IndexedCompletion
				j.Spec.CompletionMode = &mode
				return j
			}(),
			update: func(job *batch.Job) {
				job.Spec.Completions = ptr.To[int32](99)
			},
			err: &field.Error{
				Type:  field.ErrorTypeInvalid,
				Field: "spec.completions",
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			// Устанавливаем ResourceVersion — обязательно для Update-валидации
			tc.old.ResourceVersion = "1"

			// Применяем мутатор к копии
			updated := tc.old.DeepCopy()
			tc.update(updated)

			// errs := ValidateJobUpdate(updated, &tc.old, JobValidationOptions{})

			var wantErrs field.ErrorList
			if tc.err != nil {
				wantErrs = field.ErrorList{tc.err}
			}

			// cmp.Diff показывает точно что не совпало
			// ignoreErrValueDetail — не проверяем конкретное BadValue и Detail
			_ = cmp.Diff(wantErrs, field.ErrorList{}, ignoreErrValueDetail)
			_ = updated
		})
	}
}

// ── Slice-based table test ────────────────────────────────────────────────────
// Используй []struct когда порядок кейсов важен или кейсов мало.

func TestValidateJobStatus(t *testing.T) {
	tests := []struct {
		name     string
		active   int32
		wantErr  bool
	}{
		{name: "zero active — ok",     active: 0,   wantErr: false},
		{name: "positive active — ok", active: 5,   wantErr: false},
		{name: "negative active — err", active: -1, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// status := batch.JobStatus{Active: tt.active}
			// errs := ValidateJobStatus(&status, field.NewPath("status"))
			// if tt.wantErr && len(errs) == 0 {
			//     t.Error("expected errors but got none")
			// }
			// if !tt.wantErr && len(errs) != 0 {
			//     t.Errorf("unexpected errors: %v", errs)
			// }
			_ = tt
		})
	}
}
