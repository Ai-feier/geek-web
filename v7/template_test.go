//go:build v7
package web_test

import (
	web "github.com/Ai-feier/geek-web/v7"
	"github.com/stretchr/testify/require"
	"html/template"
	"log"
	"testing"
)

func TestTemplateEngine(t *testing.T) {
	// 初始化模版
	tpl, err := template.ParseGlob("testdata/tpls/*.gohtml")
	require.NoError(t, err)
	engine := &web.GoTemplateEngine{
		T: tpl,
	}
	
	// 使用 option 模式初始化服务器
	s := web.NewHTTPServer(web.ServerWithTemplateEngine(engine))
	s.GET("/login", func(ctx *web.Context) {
		err := ctx.Render("login.gohtml", nil)
		if err != nil {
			log.Println(err)
		}
	})
	s.Start(":8081")
}
