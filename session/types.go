package session

import (
	"context"
	"net/http"
)

// Store 管理 Session 本身
type Store interface {
	// Generate
	// session 对应的 ID 谁来指定？
	// 要不要在接口维度上设置超时时间，以及，要不要让 Store 内部去生成ID，都是可以自由决策
	Generate(ctx context.Context, id string) (Session, error)
	// Refresh 刷新 session
	Refresh(ctx context.Context, id string) error
	// Remove 移除 session
	Remove(ctx context.Context, id string) error
	// Get 获取 session
	Get(ctx context.Context, id string) (Session, error)
	//
	// Refresh(ctx context.Context, sess Session) error
}

// Session session 本身
type Session interface {
	Get(ctx context.Context, key string) (any, error)
	Set(ctx context.Context, key string, val any) error
	ID() string
}

// Propagator 直接与 http 框架耦合
type Propagator interface {
	// Inject 将 sessionId 注入响应
	Inject(id string, writer http.ResponseWriter) error
	// Extract 从请求中提取 sessionId
	Extract(req *http.Request) (string, error)
	// Remove 移除 session
	Remove(writer http.ResponseWriter) error
}