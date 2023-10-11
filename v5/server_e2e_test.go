//go:build v5
package web

import (
	"fmt"
	"testing"
)

// 这里放着端到端测试的代码

func TestServer(t *testing.T) {
	s := NewHTTPServer()
	s.GET("/", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello, world"))
	})
	s.GET("/user", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello, user"))
	})

	s.POST("/form", func(ctx *Context) {
		err := ctx.Req.ParseForm()
		if err != nil {
			fmt.Println(err)
		}
	})

	s.Start(":8081")
}