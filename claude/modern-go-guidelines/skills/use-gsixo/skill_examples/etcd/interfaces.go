// Example: etcd interface design patterns
// Patterns: minimal KV interface, Client embedding, context enrichment, OpOption variadic
package etcd

import (
	"context"

	"google.golang.org/grpc/metadata"
)

// ── Minimal, focused interface ─────────────────────────────────────────────────
// Каждый интерфейс описывает одну ответственность.
// Не добавляй методы "на вырост" — интерфейс легко расширить, трудно сузить.
// Именование: существительное (KV, Cluster, Lease), не глагол.

// OpOption — variadic functional option для запросов
type OpOption func(*Op)
type Op struct {
	key    []byte
	endKey []byte
	limit  int64
}

type PutResponse struct{}
type GetResponse struct{}
type DeleteResponse struct{}
type Txn interface{ Commit() error }

type KV interface {
	Put(ctx context.Context, key, val string, opts ...OpOption) (*PutResponse, error)
	Get(ctx context.Context, key string, opts ...OpOption) (*GetResponse, error)
	Delete(ctx context.Context, key string, opts ...OpOption) (*DeleteResponse, error)
	Txn(ctx context.Context) Txn
}

type Cluster interface {
	MemberList(ctx context.Context) (*MemberListResponse, error)
	MemberAdd(ctx context.Context, peerAddrs []string) (*MemberAddResponse, error)
}

type MemberListResponse struct{}
type MemberAddResponse struct{}

// ── Client через embedding ────────────────────────────────────────────────────
// Embedding интерфейсов: клиент реализует все методы всех вложенных интерфейсов.
// Преимущество: можно передавать Client туда где ожидается только KV.

type Lease interface {
	Grant(ctx context.Context, ttl int64) (*LeaseGrantResponse, error)
}
type LeaseGrantResponse struct{}

type Watcher interface {
	Watch(ctx context.Context, key string, opts ...OpOption) WatchChan
}
type WatchChan <-chan WatchResponse
type WatchResponse struct{ Err error }

type Auth interface {
	AuthEnable(ctx context.Context) error
}

// Client агрегирует все sub-интерфейсы через embedding.
// Вызывающий код может принять *Client или любой вложенный интерфейс.
type Client struct {
	Cluster
	KV
	Lease
	Watcher
	Auth
}

// Compile-time check: убеждаемся что concrete type реализует интерфейс
var _ KV = (*kvImpl)(nil)

type kvImpl struct {
	// реальная реализация
}

func (k *kvImpl) Put(_ context.Context, _, _ string, _ ...OpOption) (*PutResponse, error) {
	return nil, nil
}
func (k *kvImpl) Get(_ context.Context, _ string, _ ...OpOption) (*GetResponse, error) {
	return nil, nil
}
func (k *kvImpl) Delete(_ context.Context, _ string, _ ...OpOption) (*DeleteResponse, error) {
	return nil, nil
}
func (k *kvImpl) Txn(_ context.Context) Txn { return nil }

// ── Context enrichment ────────────────────────────────────────────────────────
// Вместо добавления параметра в каждый метод — обогащай контекст.
// Позволяет прокидывать поведение через любое количество слоёв без изменения сигнатур.
//
// gRPC metadata в исходящем контексте — стандарт etcd client.

const (
	MetadataRequireLeaderKey = "hasleader"
)

// WithRequireLeader оборачивает контекст тегом "только к лидеру".
// Используется во всей цепочке: клиент → транспорт → балансировщик.
func WithRequireLeader(ctx context.Context) context.Context {
	md := metadata.Pairs(MetadataRequireLeaderKey, "true")
	return metadata.NewOutgoingContext(ctx, md)
}
