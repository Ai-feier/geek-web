//go:build v5
package web

import (
	"net/http"
	"strings"
)

type Server interface {
	http.Handler
	// Start 启动服务器
	// addr 是监听地址。如果只指定端口，可以使用 ":8081"
	// 或者 "localhost:8082"
	Start(addr string) error

	// addRoute 注册一个路由
	// method 是 HTTP 方法
	addRoute(method string, path string, handler HandleFunc)
	// 我们并不采取这种设计方案
	// addRoute(method string, path string, handlers... HandleFunc)
}

var _ Server = &HTTPServer{}

type HTTPServer struct {
	// 组合 router
	router
}

func NewHTTPServer() *HTTPServer {
	return &HTTPServer{
		router: newRouter(),
	}
}

// Start 启动服务器
func (h *HTTPServer) Start(addr string) error {
	err := http.ListenAndServe(addr, h)
	return err
}

// ServeHTTP HTTPServer 处理请求的入口
func (h *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// 封装请求与响应
	ctx := &Context{
		Req: request,
		Resp: writer,
	}
	h.server(ctx)
}

func (h *HTTPServer) server(ctx *Context) {
	// 查找路由
	n, ok := h.findRoute(ctx.Req.Method, ctx.Req.URL.Path)
	if !ok || n.n.handler == nil {
		ctx.Resp.WriteHeader(http.StatusNotFound)
		_, _ = ctx.Resp.Write([]byte("NOT FOUND"))
		return
	}
	n.n.handler(ctx)
}

func (h *HTTPServer) GET(path string, handler HandleFunc) {
	h.addRoute(http.MethodGet, path, handler)
}

func (h *HTTPServer) POST(path string, handler HandleFunc) {
	h.addRoute(http.MethodPost, path, handler)
}

func minOperations(s1 string, s2 string, x int) int {
	if strings.Count(s1, "1") != strings.Count(s2, "1") {
		return -1
	}
	return 0
}

func countOnes(s string) int {
	count := strings.Count(s, "1")
	return count
}








