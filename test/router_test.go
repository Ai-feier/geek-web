package test

import (
	"fmt"
	"github.com/Ai-feier/geek-web"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

func Test_router_AddRoute(t *testing.T){
	testRoutes := []struct{
		method string
		path string
	} {
		{
			method: http.MethodGet,
			path: "/",
		},
		{
			method: http.MethodGet,
			path: "/user",
		},
		{
			method: http.MethodGet,
			path: "/user/home",
		},
		{
			method: http.MethodGet,
			path: "/order/detail",
		},
		{
			method: http.MethodPost,
			path: "/order/create",
		},
		{
			method: http.MethodPost,
			path: "/login",
		},
		// 通配符测试用例
		{
			method: http.MethodGet,
			path: "/order/*",
		},
		{
			method: http.MethodGet,
			path: "/*",
		},
		{
			method: http.MethodGet,
			path: "/*/*",
		},
		{
			method: http.MethodGet,
			path: "/*/abc",
		},
		{
			method: http.MethodGet,
			path: "/*/abc/*",
		},
		// 参数路由
		{
			method: http.MethodGet,
			path: "/param/:id",
		},
		{
			method: http.MethodGet,
			path: "/param/:id/detail",
		},
		{
			method: http.MethodGet,
			path: "/param/:id/*",
		},
	}

	mockHandler := func(ctx *web.Context) {}
	r := web.newRouter()
	for _, tr := range testRoutes {
		r.addRoute(tr.method, tr.path, mockHandler)
	}

	wantRouter := &web.router{
		trees: map[string]*web.node{
			http.MethodGet: {
				path: "/",
				children: map[string]*web.node{
					"user": {path: "user", children: map[string]*web.node{
						"home": {path: "home", handler: mockHandler},
					}, handler: mockHandler},
					"order": {path: "order", children: map[string]*web.node{
						"detail": {path: "detail", handler: mockHandler},
					}, starChild: &web.node{path: "*", handler: mockHandler}},
					"param": {
						path: "param",
						paramChild: &web.node{
							path: ":id",
							starChild: &web.node{
								path: "*",
								handler: mockHandler,
							},
							children: map[string]*web.node{
								"detail": {path: "detail", handler: mockHandler},
							},
							handler: mockHandler,
						},
					},
				},
				starChild:&web.node{
					path: "*",
					children: map[string]*web.node{
						"abc": {
							path: "abc",
							starChild: &web.node{path: "*", handler: mockHandler},
							handler: mockHandler},
					},
					starChild: &web.node{path: "*", handler: mockHandler},
					handler: mockHandler},
				handler: mockHandler},
			http.MethodPost: { path: "/", children: map[string]*web.node{
				"order": {path: "order", children: map[string]*web.node{
					"create": {path: "create", handler: mockHandler},
				}},
				"login": {path: "login", handler: mockHandler},
			}},
		},
	}
	
	msg, ok := wantRouter.equal(r)
	assert.True(t, ok, msg)

	// 非法用例
	r = web.newRouter()

	// 空字符串
	assert.PanicsWithValue(t, "web: 路由是空字符串", func() {
		r.addRoute(http.MethodGet, "", mockHandler)
	})

	// 前导没有 /
	assert.PanicsWithValue(t, "web: 路由必须以 / 开头", func() {
		r.addRoute(http.MethodGet, "a/b/c", mockHandler)
	})

	// 后缀有 /
	assert.PanicsWithValue(t, "web: 路由不能以 / 结尾", func() {
		r.addRoute(http.MethodGet, "/a/b/c/", mockHandler)
	})

	// 根节点重复注册
	r.addRoute(http.MethodGet, "/", mockHandler)
	assert.PanicsWithValue(t, "web: 路由冲突[/]", func() {
		r.addRoute(http.MethodGet, "/", mockHandler)
	})
	// 普通节点重复注册
	r.addRoute(http.MethodGet, "/a/b/c", mockHandler)
	assert.PanicsWithValue(t, "web: 路由冲突[/a/b/c]", func() {
		r.addRoute(http.MethodGet, "/a/b/c", mockHandler)
	})

	// 多个 /
	assert.PanicsWithValue(t, "web: 非法路由。不允许使用 //a/b, /a//b 之类的路由, [/a//b]", func() {
		r.addRoute(http.MethodGet, "/a//b", mockHandler)
	})
	assert.PanicsWithValue(t, "web: 非法路由。不允许使用 //a/b, /a//b 之类的路由, [//a/b]", func() {
		r.addRoute(http.MethodGet, "//a/b", mockHandler)
	})

	// 同时注册通配符路由和参数路由
	assert.PanicsWithValue(t, "web: 非法路由，已有通配符路由。不允许同时注册通配符路由和参数路由 [:id]", func() {
		r.addRoute(http.MethodGet, "/a/*", mockHandler)
		r.addRoute(http.MethodGet, "/a/:id", mockHandler)
	})

	assert.PanicsWithValue(t, "web: 非法路由，已有路径参数路由。不允许同时注册通配符路由和参数路由 [*]", func() {
		r.addRoute(http.MethodGet, "/a/b/:id", mockHandler)
		r.addRoute(http.MethodGet, "/a/b/*", mockHandler)
	})

	r = web.newRouter()
	assert.PanicsWithValue(t, "web: 非法路由，已有通配符路由。不允许同时注册通配符路由和参数路由 [:id]", func() {
		r.addRoute(http.MethodGet, "/*", mockHandler)
		r.addRoute(http.MethodGet, "/:id", mockHandler)
	})
	r = web.newRouter()
	assert.PanicsWithValue(t, "web: 非法路由，已有路径参数路由。不允许同时注册通配符路由和参数路由 [*]", func() {
		r.addRoute(http.MethodGet, "/:id", mockHandler)
		r.addRoute(http.MethodGet, "/*", mockHandler)
	})

	// 参数冲突
	assert.PanicsWithValue(t, "web: 路由冲突，参数路由冲突，已有 :id，新注册 :name", func() {
		r.addRoute(http.MethodGet, "/a/b/c/:id", mockHandler)
		r.addRoute(http.MethodGet, "/a/b/c/:name", mockHandler)
	})
}

func Test_router_findRouter(t *testing.T){
	testRoutes := []struct{
		method string
		path string
	} {
		{
			method: http.MethodGet,
			path: "/",
		},
		{
			method: http.MethodGet,
			path: "/user",
		},
		{
			method: http.MethodPost,
			path: "/order/create",
		},
		{
			method: http.MethodGet,
			path: "/user/*/home",
		},
		{
			method: http.MethodPost,
			path: "/order/*",
		},
		// 参数路由
		{
			method: http.MethodGet,
			path: "/param/:id",
		},
		{
			method: http.MethodGet,
			path: "/param/:id/detail",
		},
		{
			method: http.MethodGet,
			path: "/param/:id/*",
		},
	}
	
	mockHandler := func(ctx *web.Context) {}
	
	r := web.newRouter()
	for _, tr := range testRoutes {
		r.addRoute(tr.method, tr.path, mockHandler)
	}

	testCases := []struct {
		name string
		method string
		path string
		found bool
		mi *web.matchInfo
	}{
		{
			name: "method not found",
			method: http.MethodHead,
		},
		{
			name: "path not found",
			method: http.MethodGet,
			path: "/abc",
		},
		{
			name: "root",
			method: http.MethodGet,
			path: "/",
			found: true,
			mi: &web.matchInfo{
				n: &web.node{
					path: "/",
					handler: mockHandler,
				},
			},
		},
		{
			name: "user",
			method: http.MethodGet,
			path: "/user",
			found: true,
			mi: &web.matchInfo{
				n: &web.node{
					path: "user",
					handler: mockHandler,
				},
			},
		},
		{
			name: "no handler",
			method: http.MethodPost,
			path: "/order",
			found: true,
			mi: &web.matchInfo{
				n: &web.node{
					path: "order",
				},
			},
		},
		{
			name: "two layer",
			method: http.MethodPost,
			path: "/order/create",
			found: true,
			mi: &web.matchInfo{
				n: &web.node{
					path: "create",
					handler: mockHandler,
				},
			},
		},
		// 通配符匹配
		{
			// 命中/order/*
			name: "star match",
			method: http.MethodPost,
			path: "/order/delete",
			found: true,
			mi: &web.matchInfo{
				n:  &web.node{
					path: "*",
					handler: mockHandler,
				},
			},
		},
		{
			// 命中通配符在中间的
			// /user/*/home
			name: "star in middle",
			method: http.MethodGet,
			path: "/user/Tom/home",
			found: true,
			mi: &web.matchInfo{
				n:&web.node{
					path: "home",
					handler: mockHandler,
				},
			},
		},
		{
			// 比 /order/* 多了一段
			name: "overflow",
			method: http.MethodPost,
			path: "/order/delete/123",
		},
		// 参数匹配
		{
			// 命中 /param/:id
			name: ":id",
			method: http.MethodGet,
			path: "/param/123",
			found: true,
			mi: &web.matchInfo{
				n: &web.node{
					path: ":id",
					handler: mockHandler,
				},
				pathParams: map[string]string{"id": "123"},
			},
		},
		{
			// 命中 /param/:id/*
			name: ":id*",
			method: http.MethodGet,
			path: "/param/123/abc",
			found: true,
			mi: &web.matchInfo{
				n: &web.node{
					path: "*",
					handler: mockHandler,
				},
				pathParams: map[string]string{"id": "123"},
			},
		},

		{
			// 命中 /param/:id/detail
			name: ":id*",
			method: http.MethodGet,
			path: "/param/123/detail",
			found: true,
			mi: &web.matchInfo{
				n: &web.node{
					path: "detail",
					handler: mockHandler,
				},
				pathParams: map[string]string{"id": "123"},
			},
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mi, found := r.findRoute(tc.method, tc.path)
			assert.Equal(t, tc.found, found)
			if !found {
				return
			}
			wantHandler := reflect.ValueOf(tc.mi.n.handler)
			nHandler := reflect.ValueOf(mi.n.handler)
			assert.Equal(t, wantHandler, nHandler)

		})
	}
}

func (r *web.router) equal(y web.router) (string, bool) {
	for k, v := range r.trees {
		yv, ok := y.trees[k]
		if !ok {
			return fmt.Sprintf("目标 router 里面没有方法 %s 的路由树", k), false
		}
		equal, ok := v.equal(yv)
		if !ok {
			return k + "-" + equal, ok
		}
	}
	return "", true
}

func (n *web.node) equal(y *web.node) (string, bool) {
	if y == nil {
		return "目标节点为 nil", false
	}
	if n.path != y.path {
		return fmt.Sprintf("%s 节点 path 不相等 x %s, y %s", n.path, n.path, y.path), false
	}
	nhv := reflect.ValueOf(n.handler)
	yhv := reflect.ValueOf(y.handler)
	if nhv != yhv {
		return fmt.Sprintf("%s 节点 handler 不相等 x %s, y %s", n.path, nhv.Type().String(), yhv.Type().String()), false
	}
	
	if len(n.children) != len(y.children) {
		return fmt.Sprintf("%s 子节点长度不等", n.path), false
	}
	if len(n.children) == 0 {
		return "", true
	}
	
	for k, v := range n.children {
		yv, ok := y.children[k]
		if !ok {
			return fmt.Sprintf("%s 目标节点缺少子节点 %s", n.path, k), false
		}
		// 将当前节点向下匹配
		equal, ok := v.equal(yv)
		if !ok {
			return n.path + "-" + equal, false
		}
	}
	return "", true
}





















