//go:build v6 
package prometheus

import (
	web "github.com/Ai-feier/geek-web/v6"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

type MiddlewareBuilder struct {
	Name        string
	Subsystem   string
	ConstLabels map[string]string
	Help        string
}

func (m MiddlewareBuilder) Build() web.Middleware {
	// 初始化 prometheus
	summaryVec := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name: m.Name,
		Subsystem: m.Subsystem,
		ConstLabels: m.ConstLabels,
		Help: m.Help, 
	}, []string{"pattern", "method", "status"})
	prometheus.MustRegister(summaryVec)
	
	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			startTime := time.Now()
			next(ctx)
			endTime := time.Now()
			// 开启一个携程记录
			go report(endTime.Sub(startTime), ctx, summaryVec)
		}
	}
}

func report(dur time.Duration, ctx *web.Context, vec prometheus.ObserverVec) {
	status := ctx.RespStatusCode
	route := "unknown"
	if ctx.MatchedRoute != "" {
		route = ctx.MatchedRoute
	}
	ms := dur / time.Millisecond
	vec.WithLabelValues(route, ctx.Req.Method, strconv.Itoa(status)).Observe(float64(ms))
	
}
