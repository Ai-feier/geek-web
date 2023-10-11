package web

import (
	"log"
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

type HTTPServerOption func(server *HTTPServer)

type HTTPServer struct {
	// 组合 router
	router
	mdls []Middleware
	tplEngine TemplateEngine
}

func NewHTTPServer(opts ...HTTPServerOption) *HTTPServer {
	server := &HTTPServer{
		router: newRouter(),
	}
	for _, opt := range opts {
		opt(server)
	}
	return server
}

func ServerWithTemplateEngine(tplEngine TemplateEngine) HTTPServerOption {
	return func(server *HTTPServer) {
		server.tplEngine = tplEngine
	}
}

func ServerWithMiddleware(mdls ...Middleware) HTTPServerOption {
	return func(server *HTTPServer) {
		server.mdls = mdls
	}
}

func (h *HTTPServer) Use(mdls ...Middleware) {
	if h.mdls == nil {
		h.mdls = mdls
		return
	}
	h.mdls = append(h.mdls, mdls...)
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
		tplEngine: h.tplEngine,
	}
	// 最后一个应该是 HTTPServer 执行路由匹配，执行用户代码
	root := h.server
	// 从后往前组装
	for i:=len(h.mdls)-1;i>=0;i-- {
		root = h.mdls[i](root)
	}
	// 第一个应该是回写响应的
	// 因为它在调用next之后才回写响应，
	// 所以实际上 flashResp 是最后一个步骤
	var m Middleware = func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			next(ctx)
			h.flashResp(ctx)
		}
	}
	root = m(root)
	root(ctx)
}

func (h *HTTPServer) server(ctx *Context) {
	// 查找路由
	n, ok := h.findRoute(ctx.Req.Method, ctx.Req.URL.Path)
	if !ok || n.n.handler == nil {
		ctx.Resp.WriteHeader(http.StatusNotFound)
		_, _ = ctx.Resp.Write([]byte("NOT FOUND"))
		return
	}
	ctx.PathParams = n.pathParams
	ctx.MatchedRoute = n.n.route
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

func (h *HTTPServer) flashResp(ctx *Context) {
	if ctx.RespStatusCode > 0 {
		ctx.Resp.WriteHeader(ctx.RespStatusCode)
	}
	_, err := ctx.Resp.Write(ctx.RespData)
	if err != nil {
		log.Fatalln("回写响应失败", err)
	}
}






