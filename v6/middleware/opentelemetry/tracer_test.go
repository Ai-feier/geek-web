//go:build v6
package opentelemetry

import (
	web "github.com/Ai-feier/geek-web/v6"
	"go.opentelemetry.io/otel"
	"testing"
	"time"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	tracer := otel.GetTracerProvider().Tracer("")
	initZipkin(t)
	
	s := web.NewHTTPServer()
	s.GET("/", func(ctx *web.Context) {
		ctx.Resp.Write([]byte("hello, world"))
	})
	s.GET("/user", func(ctx *web.Context) {
		c, span := tracer.Start(ctx.Req.Context(), "first_layer")
		defer span.End()

		c, second := tracer.Start(c, "second_layer")
		time.Sleep(1*time.Second)

		c, third1 := tracer.Start(c, "third_layer_1")
		time.Sleep(100 * time.Millisecond)
		third1.End()
		c, third2 := tracer.Start(c, "third_layer_2")
		time.Sleep(300 * time.Millisecond)
		third2.End()
		
		second.End()
		
		ctx.RespStatusCode = 200
		ctx.RespData = []byte("/user tracing")
	})
	s.Use((&MiddlewareBuilder{Tracer: tracer}).Build())
	s.Start(":8081")
}
