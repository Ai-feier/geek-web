package session

import (
	web "github.com/Ai-feier/geek-web"
	"github.com/google/uuid"
)

type Manager struct {
	Store
	Propagator
	// sessionId
	CtxSessKey string
}

func (m *Manager) GetSession(ctx *web.Context) (Session, error) {
	if ctx.UserValues == nil {
		ctx.UserValues = make(map[string]any, 1)
	}
	val, ok := ctx.UserValues[m.CtxSessKey]  // 查缓存
	if ok {
		return val.(Session), nil
	}
	// 缓存中不存在 Session, 从 c.Req 中提取
	sessionId, err := m.Extract(ctx.Req)
	if err != nil {
		return nil, err
	}
	// 从 Store 模块中获取 sessionId
	sess, err := m.Get(ctx.Req.Context(), sessionId)
	if err != nil {
		return nil, err
	}
	ctx.UserValues[sessionId] = sessionId
	return sess, nil
}

func (m *Manager) InitSession(ctx *web.Context) (Session, error) {
	id := uuid.New().String()
	sess, err := m.Generate(ctx.Req.Context(), id)
	if err != nil {
		return nil, err
	}
	// 把 sess_id 注入到响应中
	err = m.Inject(id, ctx.Resp)
	return sess, err
}

func (m *Manager) RefreshSession(ctx *web.Context) error {
	session, err := m.GetSession(ctx)
	if err != nil {
		return err
	}
	// 将 Store 模块的 session 的过期时间进行重新设置
	return m.Refresh(ctx.Req.Context(), session.ID())
}

func (m *Manager) RemoveSession(ctx *web.Context) error {
	// 获取 session
	session, err := m.GetSession(ctx)
	if err != nil {
		return err
	}
	// 从 Store 模块删除
	err = m.Store.Remove(ctx.Req.Context(), session.ID())
	if err != nil {
		return err
	}
	// 从响应模块删除
	return m.Propagator.Remove(ctx.Resp)
}














