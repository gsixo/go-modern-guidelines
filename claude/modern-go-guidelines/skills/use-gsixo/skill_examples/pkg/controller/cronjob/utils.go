// Example: pkg/controller/cronjob/utils.go
// Patterns: iota enum with String(), sort.Interface, unexported helper types,
//           package-level vars for named constants
package cronjob

import (
	"fmt"
	"sort"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ── Package-level vars для именованных констант ───────────────────────────────
// Не магические числа — именованные переменные с комментарием.
// var, а не const, потому что time.Duration не является константой в Go.

var (
	// controllerKind используется при формировании OwnerReference
	controllerKind = batchv1.SchemeGroupVersion.WithKind("CronJob")

	// nextScheduleDelta — небольшой сдвиг при планировании следующего запуска,
	// чтобы избежать гонки с часами
	nextScheduleDelta = 100 * time.Millisecond
)

// ── iota-enum с String() ─────────────────────────────────────────────────────
// Тип-обёртка вместо bare int:
//   - компилятор проверяет правильность передачи значений
//   - String() автоматически используется в %v и klog-полях
//   - default в switch защищает от добавления новых значений без обработки

type missedSchedulesType int

const (
	noneMissed missedSchedulesType = iota
	fewMissed
	manyMissed
)

func (e missedSchedulesType) String() string {
	switch e {
	case noneMissed:
		return "none"
	case fewMissed:
		return "few"
	case manyMissed:
		return "many"
	default:
		// Всегда добавляй default — защита при расширении enum
		return fmt.Sprintf("unknown(%d)", int(e))
	}
}

// ── sort.Interface через type alias ──────────────────────────────────────────
// byJobStartTime реализует sort.Interface для сортировки джобов по времени старта.
// Тай-брейкер по имени гарантирует детерминированный порядок.
// Используется с sort.Sort(byJobStartTime(jobs)).

type byJobStartTime []*batchv1.Job

func (o byJobStartTime) Len() int      { return len(o) }
func (o byJobStartTime) Swap(i, j int) { o[i], o[j] = o[j], o[i] }
func (o byJobStartTime) Less(i, j int) bool {
	if o[i].Status.StartTime == nil && o[j].Status.StartTime == nil {
		return o[i].Name < o[j].Name
	}
	if o[i].Status.StartTime == nil {
		return false
	}
	if o[j].Status.StartTime == nil {
		return true
	}
	if o[i].Status.StartTime.Equal(o[j].Status.StartTime) {
		return o[i].Name < o[j].Name
	}
	return o[i].Status.StartTime.Before(o[j].Status.StartTime)
}

// ── Unexported helper functions ───────────────────────────────────────────────
// Функции не экспортируются — они детали реализации пакета.
// Именование: глагол + существительное, без префикса пакета.

func inActiveList(cj *batchv1.CronJob, uid interface{ String() string }) bool {
	for _, j := range cj.Status.Active {
		if j.UID == uid.String() {
			return true
		}
	}
	return false
}

func inActiveListByName(cj *batchv1.CronJob, job *batchv1.Job) bool {
	for _, j := range cj.Status.Active {
		if j.Name == job.Name && j.Namespace == job.Namespace {
			return true
		}
	}
	return false
}

// sortJobsByStartTime — удобная обёртка над sort.Sort
func sortJobsByStartTime(jobs []*batchv1.Job) {
	sort.Sort(byJobStartTime(jobs))
}

// getJobStartTime возвращает время старта джоба или zero time
func getJobStartTime(job *batchv1.Job) time.Time {
	if job.Status.StartTime != nil {
		return job.Status.StartTime.Time
	}
	return time.Time{}
}

// isJobFinished проверяет завершён ли джоб (успешно или нет)
func isJobFinished(j *batchv1.Job) (bool, batchv1.JobConditionType) {
	for _, c := range j.Status.Conditions {
		if (c.Type == batchv1.JobComplete || c.Type == batchv1.JobFailed) &&
			c.Status == "True" {
			return true, c.Type
		}
	}
	return false, ""
}

// deleteFromActiveList удаляет джоб из Active-списка CronJob
func deleteFromActiveList(cj *batchv1.CronJob, uid string) {
	if cj == nil {
		return
	}
	newActive := cj.Status.Active[:0]
	for _, j := range cj.Status.Active {
		if string(j.UID) != uid {
			newActive = append(newActive, j)
		}
	}
	cj.Status.Active = newActive
	cj.Status.LastSuccessfulTime = &metav1.Time{Time: time.Now()}
}
