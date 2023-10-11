//go:build v8
package errhdl

import (
	"bytes"
	web "github.com/Ai-feier/geek-web/v6"
	"testing"
	"text/template"
	"time"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	s := web.NewHTTPServer()
	s.GET("/", func(ctx *web.Context) {
		ctx.Resp.Write([]byte("hello, world"))
	})

	s.GET("/user", func(ctx *web.Context) {
		time.Sleep(time.Second)
	})
	
	page := `
<html>
	<h1>404 NOT FOUND</h1>
</html>
`
	tpl, err := template.New("404").Parse(page)
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	err = tpl.Execute(&buf, nil)
	if err != nil {
		t.Fatal(err)
	}
	s.Use(NewMiddlewareBuilder().RegisterError(404, buf.Bytes()).Build())
	s.Start(":8081")
}