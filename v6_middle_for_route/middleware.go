//go:build middle_for_route
package web

type Middleware func(next HandleFunc) HandleFunc

