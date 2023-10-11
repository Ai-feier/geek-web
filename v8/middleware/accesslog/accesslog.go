//go:build v8
package accesslog

import (
	"encoding/json"
	"fmt"
	web "github.com/Ai-feier/geek-web/v6"
)

type MiddlewareBuilder struct {
	logFunc func(accessLog string)
}

func (b *MiddlewareBuilder) LogFunc(f func(accessLog string)) *MiddlewareBuilder {
	b.logFunc = f 
	return b 
}

func NewBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{
		logFunc: func(accessLog string) {
			fmt.Println(accessLog)
		},
	}
}

func (b *MiddlewareBuilder) Build() web.Middleware {
	return func(next web.HandleFunc) web.HandleFunc {  // 返回 middleware
		return func(ctx *web.Context) {  // 返回 handlefunc
			defer func() {
				l := accessLog{
					Host:       ctx.Req.Host,
					Route:      ctx.MatchedRoute,
					HTTPMethod: ctx.Req.Method,
					Path:       ctx.Req.URL.Path,
				}
				data, _ := json.Marshal(l)
				b.logFunc(string(data))
			}()
			// 将context 向下传
			next(ctx)
		}
	}
}

type accessLog struct {
	Host       string
	Route      string
	HTTPMethod string `json:"http_method"`
	Path       string
}

