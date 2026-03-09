---
name: use-gsixo
description: Apply modern Go syntax guidelines based on project's Go version. Use when user ask for modern Go code guidelines. Learned on kubernetess and etcd public repos.
---

# Go Development Guide

## Detected Go Version

!`grep -rh "^go " --include="go.mod" . 2>/dev/null | cut -d' ' -f2 | sort | uniq -c | sort -nr | head -1 | xargs | cut -d' ' -f2 | grep . || echo unknown`

## How to Use This Skill

DO NOT search for go.mod files or try to detect the version yourself. Use ONLY the version shown above.

**If version detected (not "unknown"):**
- Say: "This project is using Go X.XX, so I'll stick to modern Go best practices and freely use language features up to and including this version. If you'd prefer a different target version, just let me know."
- Do NOT list features, do NOT ask for confirmation

**If version is "unknown":**
- Say: "Could not detect Go version in this repository"
- Use AskUserQuestion: "Which Go version should I target?" → [1.23] / [1.24] / [1.25] / [1.26]

**When writing Go code**, use ALL features from this document up to the target version:
- Prefer modern built-ins and packages (`slices`, `maps`, `cmp`) over legacy patterns
- Never use features from newer Go versions than the target
- Never use outdated patterns when a modern alternative is available
- All refered files in MD are stored in skill_examples. Say if you are seeing them

---

## Features by Go Version

### Go 1.0+

- `time.Since`: `time.Since(start)` instead of `time.Now().Sub(start)`

### Go 1.8+

- `time.Until`: `time.Until(deadline)` instead of `deadline.Sub(time.Now())`

### Go 1.13+

- `errors.Is`: `errors.Is(err, target)` instead of `err == target` (works with wrapped errors)

### Go 1.18+

- `any`: Use `any` instead of `interface{}`
- `bytes.Cut`: `before, after, found := bytes.Cut(b, sep)` instead of Index+slice
- `strings.Cut`: `before, after, found := strings.Cut(s, sep)`

### Go 1.19+

- `fmt.Appendf`: `buf = fmt.Appendf(buf, "x=%d", x)` instead of `[]byte(fmt.Sprintf(...))`
- `atomic.Bool`/`atomic.Int64`/`atomic.Pointer[T]`: Type-safe atomics instead of `atomic.StoreInt32`

```go
var flag atomic.Bool
flag.Store(true)
if flag.Load() { ... }

var ptr atomic.Pointer[Config]
ptr.Store(cfg)
```

### Go 1.20+

- `strings.Clone`: `strings.Clone(s)` to copy string without sharing memory
- `bytes.Clone`: `bytes.Clone(b)` to copy byte slice
- `strings.CutPrefix/CutSuffix`: `if rest, ok := strings.CutPrefix(s, "pre:"); ok { ... }`
- `errors.Join`: `errors.Join(err1, err2)` to combine multiple errors
- `context.WithCancelCause`: `ctx, cancel := context.WithCancelCause(parent)` then `cancel(err)`
- `context.Cause`: `context.Cause(ctx)` to get the error that caused cancellation

### Go 1.21+

**Built-ins:**
- `min`/`max`: `max(a, b)` instead of if/else comparisons
- `clear`: `clear(m)` to delete all map entries, `clear(s)` to zero slice elements

**slices package:**
- `slices.Contains`: `slices.Contains(items, x)` instead of manual loops
- `slices.Index`: `slices.Index(items, x)` returns index (-1 if not found)
- `slices.IndexFunc`: `slices.IndexFunc(items, func(item T) bool { return item.ID == id })`
- `slices.SortFunc`: `slices.SortFunc(items, func(a, b T) int { return cmp.Compare(a.X, b.X) })`
- `slices.Sort`: `slices.Sort(items)` for ordered types
- `slices.Max`/`slices.Min`: `slices.Max(items)` instead of manual loop
- `slices.Reverse`: `slices.Reverse(items)` instead of manual swap loop
- `slices.Compact`: `slices.Compact(items)` removes consecutive duplicates in-place
- `slices.Clip`: `slices.Clip(s)` removes unused capacity
- `slices.Clone`: `slices.Clone(s)` creates a copy

**maps package:**
- `maps.Clone`: `maps.Clone(m)` instead of manual map iteration
- `maps.Copy`: `maps.Copy(dst, src)` copies entries from src to dst
- `maps.DeleteFunc`: `maps.DeleteFunc(m, func(k K, v V) bool { return condition })`

**sync package:**
- `sync.OnceFunc`: `f := sync.OnceFunc(func() { ... })` instead of `sync.Once` + wrapper
- `sync.OnceValue`: `getter := sync.OnceValue(func() T { return computeValue() })`

**context package:**
- `context.AfterFunc`: `stop := context.AfterFunc(ctx, cleanup)` runs cleanup on cancellation
- `context.WithTimeoutCause`: `ctx, cancel := context.WithTimeoutCause(parent, d, err)`
- `context.WithDeadlineCause`: Similar with deadline instead of duration

### Go 1.22+

**Loops:**
- `for i := range n`: `for i := range len(items)` instead of `for i := 0; i < len(items); i++`
- Loop variables are now safe to capture in goroutines (each iteration has its own copy)

**cmp package:**
- `cmp.Or`: `cmp.Or(flag, env, config, "default")` returns first non-zero value

```go
// Instead of:
name := os.Getenv("NAME")
if name == "" {
    name = "default"
}
// Use:
name := cmp.Or(os.Getenv("NAME"), "default")
```

**reflect package:**
- `reflect.TypeFor`: `reflect.TypeFor[T]()` instead of `reflect.TypeOf((*T)(nil)).Elem()`

**net/http:**
- Enhanced `http.ServeMux` patterns: `mux.HandleFunc("GET /api/{id}", handler)` with method and path params
- `r.PathValue("id")` to get path parameters

### Go 1.23+

- `maps.Keys(m)` / `maps.Values(m)` return iterators
- `slices.Collect(iter)` not manual loop to build slice from iterator
- `slices.Sorted(iter)` to collect and sort in one step

```go
keys := slices.Collect(maps.Keys(m))       // not: for k := range m { keys = append(keys, k) }
sortedKeys := slices.Sorted(maps.Keys(m))  // collect + sort
for k := range maps.Keys(m) { process(k) } // iterate directly
```

**time package**

- `time.Tick`: Use `time.Tick` freely — as of Go 1.23, the garbage collector can recover unreferenced tickers, even if they haven't been stopped.

### Go 1.24+

- `t.Context()` not `context.WithCancel(context.Background())` in tests.

```go
// Before:
func TestFoo(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    result := doSomething(ctx)
}
// After:
func TestFoo(t *testing.T) {
    ctx := t.Context()
    result := doSomething(ctx)
}
```

- `omitzero` not `omitempty` for time.Duration, time.Time, structs, slices, maps in JSON tags.

```go
type Config struct {
    Timeout time.Duration `json:"timeout,omitzero"` // not omitempty — doesn't work for Duration
}
```

- `b.Loop()` not `for i := 0; i < b.N; i++` in benchmarks.

```go
func BenchmarkFoo(b *testing.B) {
    for b.Loop() {
        doWork()
    }
}
```

- `strings.SplitSeq` not `strings.Split` when iterating.

```go
for part := range strings.SplitSeq(s, ",") {
    process(part)
}
// Also: strings.FieldsSeq, bytes.SplitSeq, bytes.FieldsSeq
```

### Go 1.25+

- `wg.Go(fn)` not `wg.Add(1)` + `go func() { defer wg.Done(); ... }()`.

```go
var wg sync.WaitGroup
for _, item := range items {
    wg.Go(func() { process(item) })
}
wg.Wait()
```

### Go 1.26+

- `new(val)` not `x := val; &x` — returns pointer to any value.

```go
cfg := Config{
    Timeout: new(30),   // *int
    Debug:   new(true), // *bool
}
```

- `errors.AsType[T](err)` not `errors.As(err, &target)`.

```go
if pathErr, ok := errors.AsType[*os.PathError](err); ok {
    handle(pathErr)
}
```

---

## Best Practices from Open Source

### Kubernetes

Паттерны из `kubernetes/kubernetes`.

#### Error Handling

```go
// Sentinel errors — pkg/kubelet/container/sync_result.go
var (
    ErrCrashLoopBackOff  = errors.New("CrashLoopBackOff")
    ErrContainerNotFound = errors.New("no matching container")
    ErrRunContainer      = errors.New("RunContainerError")
)

// Typed error with context
type BackoffError struct {
    error
    backoffTime time.Time
}
func NewBackoffError(err error, t time.Time) *BackoffError {
    return &BackoffError{error: err, backoffTime: t}
}

// Error aggregation
return utilerrors.NewAggregate([]error{err1, err2})

// API errors — специфичное перед общим
switch {
case apierrors.IsNotFound(err):
    return nil, nil
case apierrors.IsConflict(err):
    return retryOnConflict(...)
case err != nil:
    return nil, err
}
```

#### Controller Pattern (Informer → Workqueue → Worker → Sync)

```go
// pkg/controller/cronjob/cronjob_controllerv2.go
type ControllerV2 struct {
    queue          workqueue.TypedRateLimitingInterface[string]
    kubeClient     clientset.Interface
    recorder       record.EventRecorder
    jobControl     jobControlInterface   // инъекция для тестов
    cronJobControl cjControlInterface
    jobLister      batchv1listers.JobLister
    jobListerSynced cache.InformerSynced
    now            func() time.Time     // инъекция времени для тестов
}

// Run: defer shutdown → sync кэша → воркеры → ctx.Done()
func (e *Signer) Run(ctx context.Context) {
    defer utilruntime.HandleCrash()
    var wg sync.WaitGroup
    defer func() { e.syncQueue.ShutDown(); wg.Wait() }()
    if !cache.WaitForNamedCacheSyncWithContext(ctx, e.configMapSynced, e.secretSynced) {
        return
    }
    wg.Go(func() { wait.UntilWithContext(ctx, e.worker, 0) })
    <-ctx.Done()
}

// processNextWorkItem
func (jm *ControllerV2) processNextWorkItem(ctx context.Context) bool {
    key, quit := jm.queue.Get()
    if quit { return false }
    defer jm.queue.Done(key)
    requeueAfter, err := jm.sync(ctx, key)
    switch {
    case err != nil:
        utilruntime.HandleError(fmt.Errorf("error syncing %v: %w", key, err))
        jm.queue.AddRateLimited(key)
    case requeueAfter != nil:
        jm.queue.Forget(key)
        jm.queue.AddAfter(key, *requeueAfter)
    default:
        jm.queue.Forget(key)
    }
    return true
}

// Sync: split key → get from lister → IsNotFound → DeepCopy → reconcile
func (jm *ControllerV2) sync(ctx context.Context, key string) (*time.Duration, error) {
    ns, name, _ := cache.SplitMetaNamespaceKey(key)
    obj, err := jm.cronJobLister.CronJobs(ns).Get(name)
    switch {
    case apierrors.IsNotFound(err):
        return nil, nil
    case err != nil:
        return nil, err
    }
    copy := obj.DeepCopy() // всегда работай с копией из кэша
    return jm.syncCronJob(ctx, copy, ...)
}
```

#### API Type Definitions

```go
// staging/src/k8s.io/api/apps/v1/types.go
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:prerelease-lifecycle-gen:introduced=1.9
type StatefulSet struct {
    metav1.TypeMeta   `json:",inline"`
    metav1.ObjectMeta `json:"metadata,omitempty"`
    Spec   StatefulSetSpec   `json:"spec,omitempty"`
    Status StatefulSetStatus `json:"status,omitempty"`
}

// +enum
type PodManagementPolicyType string
const (
    OrderedReadyPodManagement PodManagementPolicyType = "OrderedReady"
    ParallelPodManagement     PodManagementPolicyType = "Parallel"
)

// +featureGate=MaxUnavailableStatefulSet
// +optional
MaxUnavailable *intstr.IntOrString `json:"maxUnavailable,omitempty"`
```

#### Feature Gates

```go
// pkg/features/kube_features.go
const (
    // owner: @aojea
    // Allow kubelet to request a certificate without any Node IP available.
    AllowDNSOnlyNodeCSR featuregate.Feature = "AllowDNSOnlyNodeCSR"
)

if utilfeature.DefaultFeatureGate.Enabled(features.AllowDNSOnlyNodeCSR) { ... }
```

#### Interface Design

```go
// pkg/controller/cronjob/injection.go
type jobControlInterface interface {
    GetJob(namespace, name string) (*batchv1.Job, error)
    CreateJob(namespace string, job *batchv1.Job) (*batchv1.Job, error)
    DeleteJob(namespace string, name string) error
}

// Compile-time check
var _ jobControlInterface = &realJobControl{}
var _ jobControlInterface = &fakeJobControl{}
```

#### Retry & Backoff

```go
// pkg/controller/controller_utils.go
var UpdateTaintBackoff = wait.Backoff{Steps: 5, Duration: 100 * time.Millisecond, Jitter: 1.0}

err = clientretry.RetryOnConflict(UpdateTaintBackoff, func() error {
    node, err := client.CoreV1().Nodes().Get(ctx, name, metav1.GetOptions{})
    if err != nil { return err }
    // ... изменить node ...
    _, err = client.CoreV1().Nodes().Update(ctx, node, metav1.UpdateOptions{})
    return err
})
```

---

### etcd

Паттерны из `etcd-io/etcd`.

#### Error Handling

```go
// Sentinel errors
var (
    ErrStopped  = errors.New("etcdserver: server stopped")
    ErrNoLeader = errors.New("etcdserver: no leader")
    ErrTimeout  = errors.New("etcdserver: request timed out")
)

// Custom error type
type DiscoveryError struct{ Op string; Err error }
func (e DiscoveryError) Error() string {
    return fmt.Sprintf("discovery %s: %v", e.Op, e.Err)
}

// Must — для операций которые не должны падать (marshal известных типов)
func MustMarshal(m Marshaler) []byte {
    d, err := m.Marshal()
    if err != nil { panic(fmt.Sprintf("marshal should never fail (%v)", err)) }
    return d
}
```

#### Interface Design

```go
// Minimal, focused
type KV interface {
    Put(ctx context.Context, key, val string, opts ...OpOption) (*PutResponse, error)
    Get(ctx context.Context, key string, opts ...OpOption) (*GetResponse, error)
    Delete(ctx context.Context, key string, opts ...OpOption) (*DeleteResponse, error)
    Txn(ctx context.Context) Txn
}

// Composition через embedding
type Client struct { Cluster; KV; Lease; Watcher; Auth }

// Context enrichment вместо доп. параметров
func WithRequireLeader(ctx context.Context) context.Context {
    md := metadata.Pairs(rpctypes.MetadataRequireLeaderKey, "true")
    return metadata.NewOutgoingContext(ctx, md)
}
```

#### Concurrency

```go
// Typed atomics вместо sync/atomic функций
type server struct {
    appliedIndex   atomic.Uint64
    committedIndex atomic.Uint64
    isLeader       atomic.Bool
}

// Fine-grained RWMutex — отдельный на каждое поле
type server struct {
    lgMu       sync.RWMutex // protects logger
    readMu     sync.RWMutex // protects read index
    bemu       sync.RWMutex // protects backend
}

// Notifier pattern — thread-safe broadcast
type Notifier struct {
    mu      sync.RWMutex
    channel chan struct{}
}
func (n *Notifier) Notify() {
    newCh := make(chan struct{})
    n.mu.Lock()
    ch := n.channel
    n.channel = newCh
    n.mu.Unlock()
    close(ch)
}

// ID-based wait для async coordination
type Wait interface {
    Register(id uint64) <-chan any
    Trigger(id uint64, x any)
    IsRegistered(id uint64) bool
}
```

#### Configuration

```go
// Struct-based config
type Config struct {
    Endpoints        []string
    DialTimeout      time.Duration
    AutoSyncInterval time.Duration // 0 disables
    MaxCallSendMsgSize int
}

// Named constants for defaults
const (
    DefaultMaxTxnOps            = uint(128)
    DefaultMaxRequestBytes      = 1.5 * 1024 * 1024
    DefaultWarningApplyDuration = 100 * time.Millisecond
)

// Functional options для расширяемости
type BackendConfigOption func(*BackendConfig)
func WithMmapSize(size uint64) BackendConfigOption {
    return func(bcfg *BackendConfig) { bcfg.MmapSize = size }
}
func NewBackend(cfg BackendConfig, opts ...BackendConfigOption) Backend {
    for _, opt := range opts { opt(&cfg) }
    // ...
}
```

#### Metrics & Observability

```go
// Prometheus metrics — register in init()
var commitSec = prometheus.NewHistogram(prometheus.HistogramOpts{
    Namespace: "etcd", Subsystem: "disk",
    Name:    "backend_commit_duration_seconds",
    Buckets: prometheus.ExponentialBuckets(0.001, 2, 14),
})
func init() { prometheus.MustRegister(commitSec) }

// Multi-dimensional labels
var promoteFailed = prometheus.NewCounterVec(prometheus.CounterOpts{
    Name: "learner_promote_failed_total",
}, []string{"reason"})
promoteFailed.WithLabelValues("not_synced").Inc()

// expvar для runtime state
expvar.Publish("raft.status", expvar.Func(func() any {
    raftStatusMu.Lock()
    defer raftStatusMu.Unlock()
    return raftStatus()
}))
```

#### Package Organization

```
pkg/               — generic utilities (wait, notify, idutil)
client/v3/         — public client library
client/v3/internal/ — hidden implementation details
server/etcdserver/ — server core
server/storage/    — storage layer (mvcc, backend, wal)
```

- `internal/` — скрывает детали реализации
- `pkg/` — только переносимые утилиты без etcd-зависимостей
- `doc.go` в каждом пакете

#### Performance

```go
// Lock-free ID generation
type Generator struct{ prefix, suffix uint64 }
func (g *Generator) Next() uint64 {
    return g.prefix | (atomic.AddUint64(&g.suffix, 1) & suffixMask)
}

// Batch writes
type BatchConfig struct {
    batchInterval time.Duration
    batchLimit    int
}

// Contention detection
type TimeoutDetector struct{ maxDuration time.Duration; last time.Time }
func (td *TimeoutDetector) Check(now time.Time) error {
    if now.Sub(td.last) > td.maxDuration {
        return fmt.Errorf("waited %v, exceeds %v", now.Sub(td.last), td.maxDuration)
    }
    td.last = now
    return nil
}
```

---

## Useful Libraries (было → стало)

### k8s.io/utils/ptr

```go
// Было: boolPtr, int32Ptr вспомогательные функции
// Стало:
enabled = ptr.To(false)
defaultBit := ptr.Deref(bit, int32(14))
```

### k8s.io/apimachinery/pkg/util/sets

```go
// Было: map[string]struct{} + ручные методы
// Стало:
seen := sets.New[string]("a", "b")
seen.Insert("c")
seen.Has("a")
diff  := seen.Difference(other)
union := seen.Union(other)
```

### k8s.io/apimachinery/pkg/util/errors

```go
// Было: ручная конкатенация строк ошибок
// Стало:
var errs []error
for _, c := range components {
    if err := check(c); err != nil { errs = append(errs, err) }
}
return utilerrors.NewAggregate(errs) // nil если errs пуст
```

### k8s.io/apimachinery/pkg/util/wait

```go
// Было: ручной ticker + таймаут + select
// Стало:
err := wait.PollUntilContextTimeout(ctx, time.Second, timeout, true,
    func(ctx context.Context) (bool, error) {
        resp, err := client.Get(url)
        if err != nil { return false, nil }
        return resp.StatusCode == 200, nil
    },
)
```

### k8s.io/client-go/util/retry

```go
// Было: ручной цикл с IsConflict
// Стало:
err = retry.RetryOnConflict(retry.DefaultRetry, func() error {
    obj, err := client.Get(ctx, name, metav1.GetOptions{})
    if err != nil { return err }
    _, err = client.Update(ctx, mutate(obj.DeepCopy()), metav1.UpdateOptions{})
    return err
})
```

### github.com/google/go-cmp/cmp

```go
// Было: reflect.DeepEqual — непонятно что не совпало
// Стало:
var ignoreErrDetail = cmpopts.IgnoreFields(field.Error{}, "BadValue", "Detail")

if diff := cmp.Diff(want, got, ignoreErrDetail); diff != "" {
    t.Errorf("(-want,+got):\n%s", diff)
}
// Игнорировать временны́е метки:
cmp.Diff(want, got, cmpopts.IgnoreFields(v1.PodCondition{}, "LastTransitionTime"))
```

### k8s.io/apimachinery/pkg/util/runtime

```go
// Было: ручной recover в каждой горутине
// Стало:
go func() {
    defer utilruntime.HandleCrash()
    doWork()
}()
// В обработчиках событий:
utilruntime.HandleError(fmt.Errorf("unexpected nil object"))
```

### k8s.io/apimachinery/pkg/api/resource

```go
// Было: ручной парсинг единиц (Mi, Gi...)
// Стало:
mem  := resource.MustParse("256Mi")
swap := resource.MustParse("512Mi")
mem.Sub(resource.MustParse("50Mi"))
mem.Cmp(swap) // -1, 0, 1
```

### k8s.io/client-go/tools/cache (key utils)

```go
// Было: ручная конкатенация "ns/name"
// Стало:
key, _          := cache.MetaNamespaceKeyFunc(pod)        // "default/my-pod"
ns, name, _     := cache.SplitMetaNamespaceKey(key)
```

### k8s.io/klog/v2

```go
// Было: log.Printf("error for pod %s/%s: %v", ns, name, err)
// Стало:
logger := klog.FromContext(ctx)
logger.Error(err, "failed to process pod", "pod", klog.KObj(pod))
logger.V(4).Info("syncing", "key", key)
```

### github.com/stretchr/testify

```go
// Было: t.Fatal + ручные if-проверки
// Стало:
require.NoError(t, err)               // останавливает тест
assert.Equal(t, want, got)            // продолжает тест
assert.Len(t, items, 3)
require.ErrorIs(t, err, context.DeadlineExceeded)
// require — для setup; assert — для результатов
```

### go.uber.org/zap

```go
// Было: log.Printf("retrying to %s, err: %v", target, err)
// Стало:
lg.Debug("retrying request",
    zap.String("target", target),
    zap.Uint("attempt", attempt),
    zap.Error(err),
)
// В тестах:
lg := zaptest.NewLogger(t) // привязан к t.Log
```

### golang.org/x/sync/errgroup

```go
// Было: WaitGroup + канал ошибок + ручная сборка
// Стало:
g, ctx := errgroup.WithContext(parentCtx)
for _, m := range members {
    g.Go(func() error { return m.StartWithContext(ctx) })
}
return g.Wait() // ждёт всех, при первой ошибке отменяет ctx
```

### go.uber.org/multierr

```go
// Было: теряли все ошибки кроме первой
// Стало:
var combined error
for _, task := range tasks {
    combined = multierr.Append(combined, task.Close())
}
return combined // nil если все успешно, иначе все ошибки
```

### github.com/spf13/cobra

```go
// Было: ручной switch по os.Args
// Стало:
func NewMemberCommand() *cobra.Command {
    mc := &cobra.Command{Use: "member <subcommand>", Short: "Membership related commands"}
    mc.AddCommand(NewMemberAddCommand())
    return mc
}
func NewMemberAddCommand() *cobra.Command {
    cc := &cobra.Command{Use: "add <name>", Run: memberAddCommandFunc}
    cc.Flags().StringVar(&peerURLs, "peer-urls", "", "peer URLs for the new member")
    return cc
}
```

### github.com/prometheus/client_golang

```go
// Было: самодельные счётчики + ручной /metrics
// Стало:
var commitSec = prometheus.NewHistogram(prometheus.HistogramOpts{
    Namespace: "etcd", Subsystem: "disk",
    Name:    "backend_commit_duration_seconds",
    Buckets: prometheus.ExponentialBuckets(0.001, 2, 14),
})
func init() { prometheus.MustRegister(commitSec) }
func (b *backend) commit() {
    start := time.Now()
    commitSec.Observe(time.Since(start).Seconds())
}
```

---

## Unit Testing Recommendations

### Table-driven Tests

```go
// map[string]struct{} для независимых кейсов — ключ = имя теста
cases := map[string]struct {
    opts JobValidationOptions
    job  batch.Job
}{
    "valid success policy": { opts: ..., job: ... },
    "invalid completion mode": { ... },
}
for name, tc := range cases {
    t.Run(name, func(t *testing.T) { ... })
}

// []struct{name string, ...} когда важен порядок
tests := []struct{ name string; input, want string }{{
    name: "basic", input: "foo", want: "FOO",
}}
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) { ... })
}

// update-функция в кейсе для мутационных тестов
cases := map[string]struct {
    old    batch.Job
    update func(*batch.Job)
    err    *field.Error
}{
    "immutable selector": {
        old:    validJob(),
        update: func(job *batch.Job) { job.Spec.Selector = newSelector },
        err:    &field.Error{Type: field.ErrorTypeInvalid, Field: "spec.selector"},
    },
}
for name, tc := range cases {
    t.Run(name, func(t *testing.T) {
        updated := tc.old.DeepCopy()
        tc.update(updated)
        errs := ValidateJobUpdate(updated, &tc.old, opts)
        if diff := cmp.Diff(tc.err, errs[0], ignoreErrDetail); diff != "" {
            t.Errorf("(-want,+got):\n%s", diff)
        }
    })
}
```

### Fixtures & Builders

```go
// Tweak-паттерн (k8s)
type Tweak func(*api.Pod)
func MakePod(name string, tweaks ...Tweak) *api.Pod {
    pod := &api.Pod{ ObjectMeta: metav1.ObjectMeta{Name: name} }
    for _, t := range tweaks { t(pod) }
    return pod
}
func SetNamespace(ns string) Tweak { return func(p *api.Pod) { p.Namespace = ns } }

// getValid* конструкторы
func getValidManualSelector() *metav1.LabelSelector {
    return &metav1.LabelSelector{MatchLabels: map[string]string{"a": "b"}}
}

// Модификаторы для добавления опций
func addSecretExpiration(s *v1.Secret, exp string) {
    s.Data[bootstrapapi.BootstrapTokenExpirationKey] = []byte(exp)
}

// ptr.To для pointer-полей
job.Spec.Parallelism = ptr.To[int32](5)
```

### Fake Clients & Action Verification

```go
// Setup
func newSigner() (*Signer, *fake.Clientset, ...) {
    cl := fake.NewSimpleClientset()
    inf := informers.NewSharedInformerFactory(cl, 0)
    bsc, _ := NewSigner(cl, inf.Core().V1().Secrets(), inf.Core().V1().ConfigMaps(), DefaultSignerOptions())
    return bsc, cl, ...
}

// Добавлять объекты через indexer, не через API
inf.GetIndexer().Add(obj)

// Верифицировать API-вызовы
func verifyActions(t *testing.T, expected, actual []core.Action) {
    t.Helper()
    for i, a := range actual {
        if len(expected) < i+1 {
            t.Errorf("%d unexpected actions", len(actual)-len(expected)); break
        }
        if !helper.Semantic.DeepEqual(expected[i], a) {
            t.Errorf("Expected\n\t%s\ngot\n\t%s", dump.Pretty(expected[i]), dump.Pretty(a))
        }
    }
}

// PrependReactor для side-effects
fakeClient.PrependReactor("create", "ipaddresses", func(action k8stesting.Action) (bool, runtime.Object, error) {
    ip := action.(k8stesting.CreateAction).GetObject().(*networkingv1.IPAddress)
    _ = ipIndexer.Add(ip)
    return false, ip, nil
})
```

### Error Assertion

```go
// field.ErrorList + cmp.Diff + cmpopts
var ignoreErrDetail = cmpopts.IgnoreFields(field.Error{}, "BadValue", "Detail", "Origin")

wantErrs := field.ErrorList{
    {Type: field.ErrorTypeInvalid, Field: "spec.selector"},
    {Type: field.ErrorTypeRequired, Field: "spec.template"},
}
if diff := cmp.Diff(wantErrs, gotErrs, ignoreErrDetail); diff != "" {
    t.Errorf("(-want,+got):\n%s", diff)
}

// API errors
if !apierrors.IsNotFound(err) {
    t.Errorf("expected NotFound, got: %v", err)
}
```

### Assertion Style

```go
// Предпочтительно — stdlib
if got != want { t.Errorf("got %v, want %v", got, want) }

// Допустимо — testify (меньше бойлерплейта)
require.NoError(t, err)    // останавливает (для setup)
assert.Equal(t, want, got) // продолжает (для результатов)

// t.Fatal — когда дальше нет смысла; t.Error — когда можно продолжить
obj, err := getObject()
if err != nil { t.Fatalf("setup failed: %v", err) }
if obj.Name != want { t.Errorf("wrong name: %s", obj.Name) }
```

### Time, Clocks, Concurrency in Tests

```go
// Инъекция времени через func() time.Time
jm.now = func() time.Time { return tc.now }

// Fake clock
import testclock "k8s.io/utils/clock/testing"
fakeClock := testclock.NewFakeClock(time.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC))
fakeClock.Step(5 * time.Minute)

// ktesting.NewTestContext вместо context.Background()
_, ctx := ktesting.NewTestContext(t)
ctx, cancel := context.WithCancel(ctx)
defer cancel()
```

### Hand-written Fakes vs Mocks

```go
// Предпочтительно — hand-written fake (k8s подход)
type fakeJobControl struct {
    sync.Mutex
    Jobs      []batchv1.Job
    CreateErr error
}
func (f *fakeJobControl) CreateJob(_ context.Context, ns string, job *batchv1.Job) (*batchv1.Job, error) {
    f.Lock(); defer f.Unlock()
    if f.CreateErr != nil { return nil, f.CreateErr }
    f.Jobs = append(f.Jobs, *job)
    return job, nil
}
var _ jobControlInterface = &fakeJobControl{}

// Recorder pattern для проверки вызовов (etcd подход)
type WaitRecorder struct {
    wait.Wait
    Recorder
}
func (w *WaitRecorder) Register(id uint64) <-chan any {
    w.Record(Action{Name: "Register", Params: []any{id}})
    return w.Wait.Register(id)
}

// mockery (vektra) — только для сложных интерфейсов
func NewMockDevicesProvider(t mock.TestingT) *MockDevicesProvider {
    m := &MockDevicesProvider{}
    m.Mock.Test(t)
    t.Cleanup(func() { m.AssertExpectations(t) })
    return m
}
```

### Parallel & Benchmarks

```go
// t.Parallel() внутри t.Run()
for _, tc := range testCases {
    tc := tc // захват переменной (Go < 1.22)
    t.Run(tc.name, func(t *testing.T) {
        t.Parallel()
        // ...
    })
}

// Тест на panic (etcd подход)
defer func() {
    if r := recover(); r == nil { t.Error("expected panic") }
}()

// Benchmark
func BenchmarkFoo(b *testing.B) {
    for b.Loop() { doWork() } // Go 1.24+
}
```

---

## Code Conventions (CC)

> Unified rules from `kubernetes/kubernetes` and `etcd-io/etcd`.

---

### CC-NAMING

- **Константы** — `PascalCase`, не `UPPER_SNAKE_CASE`: `CronJobScheduledTimestampAnnotation`, `DefaultMaxTxnOps`
- **Enum-константы** — имя типа как префикс значения: `PodFailurePolicyActionFailJob`, `OrderedReadyPodManagement`
- **Интерфейсы** — семантическое существительное: `MetricsClient`, `KV`, `Backend`, `Wait`; реализация — строчная: `type kv struct`
- **Sentinel errors** — `Err`-префикс, package-level var: `ErrNoLeader`, `ErrCrashLoopBackOff`
- **Булевы** — `is*`, `has*` или глагол: `isFinished`, `updateStatus`, `found`
- **Ресиверы** — 1-2 символа, аббревиатура типа: `jm` (job manager), `s` (server), `c` (client)
- **Аббревиатуры** — PascalCase в составных именах: `ControllerUid`, `HTTPHandler`
- **Функциональные опции** — тип `XxxOption`, конструктор `WithXxx`: `BackendConfigOption`, `WithMmapSize`
- **Пакет** — одно слово, строчное, без подчёркиваний: `wait`, `notify`, `idutil`
- **Channel-переменные** — суффикс `c`/`C`: `applyc`, `msgSnapC`, `readStateC`

---

### CC-FUNC

- **Сигнатура**: `ctx context.Context` первым, `error` последним
- **Конструктор**: `New*`, возвращает `(T, error)` или интерфейс, не конкретный тип
- **Длинные параметры** — каждый на своей строке
- **Named return values** — не использовать
- **Variadic options** для необязательной конфигурации: `opts ...OpOption`
- **Options-struct** когда параметров ≥ 4: `SignerOptions`, `CIDRAllocatorParams`
- **`DefaultXxxOptions()`** для разумных дефолтов рядом с options-struct
- **Несколько конструкторов**: `New` — основной, `NewXxx` — специализированный
- **Тест-хелперы**: `t.Helper()` первой строкой, `*testing.T` первым аргументом

---

### CC-COMMENT

- **Godoc** начинается с имени символа: `// NewControllerV2 creates...`, `// Job represents...`
- **Поля struct** — краткое описание + `+optional`, ссылки на документацию
- **Mutex** — комментарий рядом объясняет что защищает: `// readMu protects concurrent access to readwaitc`
- **Маркеры**: `+genclient`, `+k8s:deepcopy-gen`, `+optional`, `+featureGate=Xxx`, `+enum`
- **Feature gate** — `// owner: @username` и `// kep: https://kep.k8s.io/NNN`
- **TODO** — простой формат с объяснением: `// TODO: group and encapsulate the reads in a single tx`
- **Inline** — только там, где логика неочевидна; объясняй ПОЧЕМУ, не ЧТО
- **`doc.go`** в каждом пакете: `// Package wait provides utility functions for polling and notification.`

---

### CC-FORMAT

- **Импорты** — три группы, разделены пустой строкой:
  1. stdlib (`context`, `fmt`, `sync`, ...)
  2. внешние (`k8s.io/`, `github.com/`, `go.uber.org/`, ...)
  3. внутренние (`k8s.io/kubernetes/...`, `go.etcd.io/etcd/...`)
- **Длина строки** — ~100 символов, прагматично; сигнатуры и logging-вызовы могут превышать
- **`gofmt`** обязателен; `goimports` для управления импортами
- **Поля структуры** — логические группы разделены пустой строкой; атомики и мьютексы в начале

---

### CC-OBJECTS

- **Контроллер**: queue → clients → control interfaces → listers → synced checks → test helpers
- **API-тип**: `TypeMeta` → `ObjectMeta` → `Spec` → `Status`
- **Опциональные поля** — указатели `*T` + `// +optional` + `omitempty`
- **Mutex** — перед полем которое защищает; embedded `sync.Mutex` первым в тестовых структурах
- **Атомики** — в начале структуры; атомики типа `atomic.Int64`, не `sync/atomic` функции
- **`RWMutex`** по умолчанию для read-heavy данных
- **Каналы** — инициализируются в конструкторе с явным размером буфера и комментарием
- **Embedding** — только с объяснительным комментарием; для публичного API через интерфейсы

---

### CC-ERROR

- **Строки ошибок** — lowercase, без точки: `"no matching container"`, `"etcdserver: no leader"`
- **Оборачивай с `%w`** когда добавляешь контекст; возвращай as-is если контекст уже есть
- **`switch {}`** для разных типов ошибок — специфичное перед общим (`IsNotFound` → `IsConflict` → `err != nil`)
- **Не логируй и не возвращай** одновременно — выбери одно; `utilruntime.HandleError` + requeue или `return err`
- **`panic`** — только для нарушения внутренних инвариантов; для дисковых неисправимых ошибок — `lg.Fatal()`
- **Не проглатывай** — каждый `err` либо возвращается, либо логируется, либо `panic`
- **Кастомный тип** для структурированного контекста: `type DiscoveryError struct{ Op string; Err error }`

---

### CC-BOUNDARY

- **Экспортируй минимум**: только `New*`, `Run()` и то, что нужно снаружи
- **Интерфейс для каждой внешней зависимости** — для тестируемости
- **Compile-time check**: `var _ Interface = &Implementation{}`
- **Конструктор возвращает интерфейс**, не конкретный тип: `func New(...) Backend`
- **`internal/`** для скрытия деталей реализации
- **`pkg/`** — только переносимые утилиты без domain-специфичных зависимостей
- **Context как контракт**: отсутствие `ctx` в публичном методе — нарушение границы
- **Options-struct** для группировки параметров конфигурации

---

### CC-TEST

- **Table-driven** с `t.Run(name, ...)` для любых параметризованных случаев
- **`t.Helper()`** — первая строка любой вспомогательной функции теста
- **`require`** для setup и предусловий, **`assert`** для результатов
- **`t.Context()`** (Go 1.24+) вместо `context.WithCancel(context.Background())`
- **`zaptest.NewLogger(t)`** — единственный способ создать логгер в тестах
- **`cmp.Diff`** вместо `reflect.DeepEqual`; `cmpopts.IgnoreFields` для нестабильных полей
- **Добавлять объекты в informer через `GetIndexer().Add()`**, не через fake API
- **Инъекция времени** через `now func() time.Time` / `testclock.NewFakeClock`
- **`testutil.RegisterLeakDetection(t)`** для обнаружения утечек горутин (etcd)

---

### CC-CLASS

- **Typed enum**: `type T string` + `const (...)` + `// +enum`
- **iota-enum** с `String()` методом для внутренних типов — удобен в логах
- **Compile-time check**: `var _ KV = (*kv)(nil)`
- **Новый тип** вместо алиаса когда нужна type-safety: `type LeaseID int64`
- **Псевдоним** только для protobuf abbreviation: `type PutResponse = pb.PutResponse`
- **Функциональные опции** — предпочтительнее config-struct для небольшого числа опций
- **Fluent interface** для builder-паттернов; защищать от двойного вызова через panic

---

### CC-SYSTEM

- **Пакет = одна ответственность** — отдельно: логика, интерфейсы/fakes, утилиты, метрики, конфиг
- **Слои зависимостей строго вниз**, без циклов
- **Feature gates** — константы типа `featuregate.Feature`; регистрация через blank import
- **Инициализация** — конструктор создаёт, `Start()`/`Run()` запускает горутины; `Run()` не блокирует
- **`Stop()` + `done` channel** — стандартный паттерн завершения: `Stop()` сигналит, `<-s.done` ждёт
- **Зависимости через конструктор**, не через глобальные переменные

---

### CC-REFINE

- **`DeepCopy()`** перед любой мутацией объекта из кэша информера
- **`switch {}`** вместо `if/else if` для цепочек error-проверок
- **Магические числа** → именованные package-level переменные с комментарием
- **Тип-обёртка + `String()`** вместо bare-int для читаемости в логах
- **Не возвращай конкретный тип из конструктора** — измени на интерфейс
- **Один `sync.Mutex` на одно логическое поле** — не для несвязанных полей
- **Группируй связанные константы** в `const (...)` блоки с комментарием
- **`_ = err`** с комментарием если игнорирование намеренно; иначе обрабатывай

---

### CC-CONCURRENCY

- **Run()-паттерн**: defer shutdown → sync кэша → запуск воркеров → `<-ctx.Done()`
- **`defer Unlock()`** сразу после `Lock()` — никогда не полагайся на ручной Unlock
- **`wg.Go()`** (Go 1.25+) или цикл с `wg.Add(1)` + горутина для воркеров
- **Каждая горутина имеет выход**: `<-ctx.Done()` или `<-stopped` в select
- **Владение каналом**: создающий — закрывает; получатель никогда не закрывает
- **`chan struct{}`** для сигналов; `close(ch)` — broadcast всем читателям
- **Буфер = 1** для каналов-нотификаторов; unbuffered для backpressure
- **`select` с `default`** только для non-blocking проверки
- **`errgroup`** для параллельных задач с ошибками вместо WaitGroup + канал
- **Fine-grained мьютексы**: отдельный RWMutex на каждое логически независимое поле

---

### CC-SMELL

| ❌ Антипаттерн | ✓ Правильно |
|---------------|-------------|
| `context.TODO()` в методах | Принять `ctx context.Context` параметром |
| `ctx` в поле структуры | Передавать в каждый метод |
| `panic()` для обычных ошибок | Возвращать `error`; panic только для инвариантов |
| `init()` в рукописном коде | Явный `New*()` конструктор |
| Глобальное изменяемое состояние | Инкапсулировать в структуру + mutex |
| 4+ параметра одного типа подряд | Options-struct |
| `reflect.DeepEqual` в тестах | `cmp.Diff` с понятным diff |
| Логировать + возвращать ошибку | Выбрать одно |
| `fmt.Println` в production-коде | `klog` / `zap` |
| Один большой mutex на всю структуру | Fine-grained мьютексы |
| Неинициализированный канал | Всегда `make(chan T)` |
| Горутина без exit-условия | `<-ctx.Done()` или `<-stopped` |
| `time.Sleep` в тестах | Channel / condition / `testutil.Poll` |
| Конкретный тип из конструктора | Возвращать интерфейс |
| Глобальный логгер `zap.L()` в библиотеке | Инжектировать через конструктор |

---

## File Index by Directory

Все файлы Kubernetes, упомянутые в этом guide, сгруппированы по директориям.

---

### `cmd/`

#### `cmd/kube-apiserver/app/`
| Файл | Что демонстрирует |
|------|-------------------|
| `options/options.go` | Embedding struct с комментарием; blank import для регистрации feature gates |

#### `cmd/kube-proxy/app/`
| Файл | Что демонстрирует |
|------|-------------------|
| `options.go` | `ptr.To()`, `ptr.Deref()` для pointer-полей конфигурации; structured logging с klog |

#### `cmd/kubeadm/app/`
| Файл | Что демонстрирует |
|------|-------------------|
| `componentconfigs/kubelet.go` | `ptr.To()` для bool/int32 полей kubelet-конфига |
| `phases/controlplane/manifests_test.go` | `intstr.FromString()` / `intstr.FromInt32()` в тестах probe |
| `util/apiclient/wait.go` | `wait.PollUntilContextTimeout`; `utilerrors.NewAggregate` |

---

### `pkg/`

#### `pkg/api/pod/testing/`
| Файл | Что демонстрирует |
|------|-------------------|
| `make.go` | **Tweak-паттерн**: `type Tweak func(*api.Pod)`, `MakePod(name, ...Tweak)`, функции-твики `SetNamespace`, `SetNodeName` |

#### `pkg/apis/batch/validation/`
| Файл | Что демонстрирует |
|------|-------------------|
| `validation_test.go` | Table-driven tests через `map[string]struct{}`; `update func(*batch.Job)` в кейсе; `cmp.Diff` + `cmpopts.IgnoreFields`; `ptr.To[int32]()` |

#### `pkg/apis/core/v1/helper/`
| Файл | Что демонстрирует |
|------|-------------------|
| `helpers_test.go` | `t.Parallel()` внутри `t.Run()` для параллельных подтестов |

#### `pkg/controller/bootstrap/`
| Файл | Что демонстрирует |
|------|-------------------|
| `bootstrapsigner.go` | `SignerOptions` + `DefaultSignerOptions()`; `FilteringResourceEventHandler`; Run-паттерн (defer shutdown → cache sync → воркеры → ctx.Done) |
| `bootstrapsigner_test.go` | Factory-функция `newSigner()`; `fake.NewSimpleClientset()`; `GetIndexer().Add()`; `verifyActions()` |
| `common_test.go` | `newTokenSecret()`, `addSecretExpiration()` — модификаторы тестовых объектов; `verifyActions()` helper |

#### `pkg/controller/`
| Файл | Что демонстрирует |
|------|-------------------|
| `controller_utils.go` | `wait.Backoff` конфиги как package-level vars; `clientretry.RetryOnConflict` |

#### `pkg/controller/cronjob/`
| Файл | Что демонстрирует |
|------|-------------------|
| `cronjob_controllerv2.go` | Полный контроллер-паттерн: структура, `Run()`, `processNextWorkItem()`, `sync()`; `klog.FromContext(ctx)`; `DeepCopy()` перед мутацией; `switch {}` для error-handling |
| `cronjob_controllerv2_test.go` | `ktesting.NewTestContext(t)`; `jm.now = func() time.Time { return tt.now }`; `GetIndexer().Add()` |
| `injection.go` | `jobControlInterface` / `cjControlInterface` — малые интерфейсы; `realCJControl` + `fakeCJControl`; `var _ I = &T{}` compile-time check; `context.TODO()` — пример legacy антипаттерна |
| `utils.go` | `missedSchedulesType` — iota-enum с `String()` методом; `byJobStartTime` — sort.Interface реализация; unexported helper functions |

#### `pkg/controller/nodeipam/`
| Файл | Что демонстрирует |
|------|-------------------|
| `node_ipam_controller.go` | `CIDRAllocatorParams` options-struct; минимальный экспорт (только `New*` + `Run`) |
| `ipam/cidr_allocator.go` | `CIDRAllocator` — Strategy-паттерн через интерфейс |

#### `pkg/controller/podautoscaler/metrics/`
| Файл | Что демонстрирует |
|------|-------------------|
| `interfaces.go` | `MetricsClient` — именование интерфейса существительным, не `-er`-суффиксом |

#### `pkg/features/`
| Файл | Что демонстрирует |
|------|-------------------|
| `kube_features.go` | Feature gate константы: `featuregate.Feature`, `// owner: @username`, алфавитный порядок |

#### `pkg/kubelet/apis/podresources/testing/`
| Файл | Что демонстрирует |
|------|-------------------|
| `mocks.go` | Auto-generated mockery mock: `NewMockDevicesProvider(t)`, `EXPECT()` pattern, `mock.AssertExpectations` |

#### `pkg/kubelet/container/`
| Файл | Что демонстрирует |
|------|-------------------|
| `sync_result.go` | Sentinel errors (`ErrCrashLoopBackOff`...); `BackoffError` typed error; `MinBackoffExpiration` — рекурсивный `errors.Unwrap`; `%w` оборачивание; `utilerrors.NewAggregate` |

#### `pkg/kubelet/images/pullmanager/`
| Файл | Что демонстрирует |
|------|-------------------|
| `benchmarks_test.go` | Benchmark-паттерн: `b.N`, benchmark-helper функции |

#### `pkg/proxy/`
| Файл | Что демонстрирует |
|------|-------------------|
| `endpointschangetracker_test.go` | `sets.New[string]()`, `sets.Set[string]` в тестовых expected-значениях |
| `nftables/proxier.go` | `sets.New[string]()` в production-коде |

#### `pkg/registry/core/service/ipallocator/controller/`
| Файл | Что демонстрирует |
|------|-------------------|
| `repairip_test.go` | `testclock.NewFakeClock()`; `PrependReactor` для create/update/delete; несколько informer-типов в одном тесте |

#### `pkg/registry/core/service/storage/`
| Файл | Что демонстрирует |
|------|-------------------|
| `transaction_test.go` | `[]struct{name string, ...}` slice-based table-driven test; `t.Parallel()` на уровне теста |

#### `pkg/scheduler/`
| Файл | Что демонстрирует |
|------|-------------------|
| `scheduler_test.go` | `cache.MetaNamespaceKeyFunc(pod)` для получения ключа |

---

### `staging/`

#### `staging/src/k8s.io/api/apps/v1/`
| Файл | Что демонстрирует |
|------|-------------------|
| `types.go` | `TypeMeta → ObjectMeta → Spec → Status`; typed enums (`PodManagementPolicyType`); `+genclient`, `+k8s:deepcopy-gen`, `+featureGate`, `+optional` маркеры; `*int32` для optional полей |

#### `staging/src/k8s.io/api/batch/v1/`
| Файл | Что демонстрирует |
|------|-------------------|
| `types.go` | `CompletionMode`, `PodFailurePolicyAction` — enum-константы с документацией; `CronJobScheduledTimestampAnnotation` — PascalCase константа |

#### `staging/src/k8s.io/apimachinery/pkg/util/errors/`
| Файл | Что демонстрирует |
|------|-------------------|
| `errors.go` | `Aggregate` interface; `NewAggregate([]error)` |

#### `staging/src/k8s.io/endpointslice/util/`
| Файл | Что демонстрирует |
|------|-------------------|
| `controller_utils.go` | `utilruntime.HandleError()` в event-handler'ах |

---

### `test/`

#### `test/e2e/apimachinery/`
| Файл | Что демонстрирует |
|------|-------------------|
| `validatingadmissionpolicy.go` | `retry.RetryOnConflict(retry.DefaultRetry, ...)` |

#### `test/e2e_node/`
| Файл | Что демонстрирует |
|------|-------------------|
| `swap_test.go` | `resource.MustParse("256Mi")`, `resource.Quantity` операции |

#### `test/integration/podgc/`
| Файл | Что демонстрирует |
|------|-------------------|
| `podgc_test.go` | `cmp.Diff` с `cmpopts.IgnoreFields` для временны́х меток |
