package test

import (
	"fmt"
	"github.com/Ai-feier/geek-web"
	"testing"
)

// 这里放着端到端测试的代码

func TestServer(t *testing.T) {
	s := web.NewHTTPServer()
	s.GET("/", func(ctx *web.Context) {
		ctx.Resp.Write([]byte("hello, world"))
	})
	s.GET("/user", func(ctx *web.Context) {
		ctx.Resp.Write([]byte("hello, user"))
	})

	s.POST("/form", func(ctx *web.Context) {
		err := ctx.Req.ParseForm()
		if err != nil {
			fmt.Println(err)
		}
	})

	s.Start(":8081")
}