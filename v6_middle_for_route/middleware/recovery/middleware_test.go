//go:build middle_for_route
package recovery

import (
	web "github.com/Ai-feier/geek-web/v6"
	"log"
	"testing"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	s := web.NewHTTPServer()
	s.GET("/", func(ctx *web.Context) {
		ctx.Resp.Write([]byte("hello world"))
	})
	s.GET("/panic", func(ctx *web.Context) {
		panic("测试 panic 恢复的 middleware")
	})
	builder := &MiddlewareBuilder{
		StatusCode: 500,
		ErrMsg:     "成功从 panic 中恢复",
		LogFunc: func(ctx *web.Context) {
			log.Println(ctx.Req.URL.Path)
		},
	}
	s.Use(builder.Build())
	s.Start(":8081")
}
