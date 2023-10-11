//go:build v1
package v1

import "net/http"

type HandleFunc func(ctx *Context)

type Server interface {
	http.Handler
	// Start 启动服务器
	// addr 是监听地址。如果只指定端口，可以使用 ":8081"
	// 或者 "localhost:8082"
	Start(addr string) error
	// AddRoute 注册一个路由
	// method 是 HTTP 方法
	// path 是路径，必须以 / 为开头
	AddRoute(method string, path string, handler HandleFunc)
	// 我们并不采取这种设计方案
	// addRoute(method string, path string, handlers... HandleFunc)
}

var _ Server = &HttpServer{}

type HttpServer struct {
	
}
// ServeHTTP HTTPServer 处理请求的入口
func (s *HttpServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := &Context{
		Req:  request,
		Resp: writer,
	}
	
	// 可加入 web 服务启动前的 hook 方法
	// - 可用于服务注册
	
	s.serve(ctx)
}

func (s *HttpServer) Start(addr string) error {
	err := http.ListenAndServe(":8080", s)
	return err
}

func (s *HttpServer) AddRoute(method string, path string, handler HandleFunc) {
	//TODO implement me
	panic("implement me")
}

func (s *HttpServer) Get(path string, handler HandleFunc) {
	s.AddRoute(http.MethodGet, path, handler)
}

func (s *HttpServer) Post(path string, handler HandleFunc) {
	s.AddRoute(http.MethodPost, path, handler)
}

func (s *HttpServer) serve(ctx *Context) {
	
}





















