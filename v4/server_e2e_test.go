//go:build v4
package web

import "testing"

// 这里放着端到端测试的代码

func TestServer(t *testing.T) {
	s := NewHTTPServer()
	
	// 添加路由
	s.GET("/", func(ctx *Context) {
		ctx.Resp.Write([]byte("Hello World"))
	})
	s.POST("/user", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello user"))
	})
	
	err := s.Start(":8081")
	if err != nil {
		panic(err)
	}
}