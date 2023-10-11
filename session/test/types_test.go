package test

import (
	web "github.com/Ai-feier/geek-web"
	"github.com/Ai-feier/geek-web/session"
	"github.com/Ai-feier/geek-web/session/cookie"
	"github.com/Ai-feier/geek-web/session/memory"
	"net/http"
	"testing"
	"time"
)

func TestSession(t *testing.T) {
	// 进行简单的登录校验
	var m *session.Manager = &session.Manager{
		Propagator: cookie.NewPropagator(),
		Store: memory.NewStore(30*time.Minute),
		CtxSessKey: "sessId",
	}
	// 为 http 服务注册登录校验 middleware
	server := web.NewHTTPServer(web.ServerWithMiddleware(func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			// 进行登录校验
			if ctx.Req.URL.Path == "/login" {
				// 登录不进行校验
				next(ctx)
				return
			}

			// 从请求中提取 sessionId  -- Propagator 模块
			// 根据 sessionId 从 Store 模块查询对应的 Session -- Store 模块
			// -- 整合到 Manage 进行统一管理
			_, err := m.GetSession(ctx)
			if err != nil {
				ctx.RespStatusCode = http.StatusUnauthorized
				ctx.RespData = []byte("请重新登录")
				return
			}
			
			// 刷新 session
			_ = m.RefreshSession(ctx)
			
			next(ctx)
		}
	}))

	// 登录
	server.POST("/login", func(ctx *web.Context) {
		// 进行用户名和密码的校验
		
		sess, err := m.InitSession(ctx)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("登录失败")
			return
		}
		// 将常用数据存入 session
		err = sess.Set(ctx.Req.Context(), "nickname", "mkt")
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("登录失败")
			return
		}

		ctx.RespStatusCode = http.StatusOK
		ctx.RespData = []byte("登录成功")
		return
	})
	
	// 退出登录
	server.POST("/logout", func(ctx *web.Context) {
		// 清楚数据
		err := m.RemoveSession(ctx)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("退出失败")
			return
		}
		ctx.RespStatusCode = http.StatusOK
		ctx.RespData = []byte("退出成功")
	})
	
	server.GET("/user", func(ctx *web.Context) {
		sess, _ := m.GetSession(ctx)
		// 从 session 中取值
		val, _ := sess.Get(ctx.Req.Context(), "nickname")
		ctx.RespData = []byte(val.(string))
	})
	server.Start(":8081")
}
