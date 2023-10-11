//go:build v7
package errhdl

import web "github.com/Ai-feier/geek-web/v6"

type MiddlewareBuilder struct {
	resp map[int][]byte
}

func NewMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{
		// 这里可以非常大方，因为在预计中用户会关心的错误码不可能超过 64
		resp: make(map[int][]byte, 64),
	}
}

func (m *MiddlewareBuilder) RegisterError(code int, resp []byte) *MiddlewareBuilder {
	m.resp[code] = resp
	return m 
}

func (m MiddlewareBuilder) Build() web.Middleware {
	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			next(ctx)
			data, ok := m.resp[ctx.RespStatusCode]
			if ok {
				ctx.RespData = data
			}
		}
	}
}
