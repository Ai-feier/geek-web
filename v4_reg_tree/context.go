//go:build reg_tree
package web

import "net/http"

type Context struct {
	Req        *http.Request
	Resp       http.ResponseWriter
	PathParams map[string]string
}
