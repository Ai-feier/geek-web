package cookie

import "net/http"

type PropagatorOption func(p *Propagator)

type Propagator struct {
	cookieName string
	cookieOption func(c *http.Cookie)
}

func NewPropagator() *Propagator {
	return &Propagator{
		cookieName: "sessId",
		cookieOption: func(c *http.Cookie) {
		},
	}
}

func WithCookieName(name string) PropagatorOption {
	return func(p *Propagator) {
		p.cookieName = name
	}
}

func (p *Propagator) Inject(id string, writer http.ResponseWriter) error {
	c := &http.Cookie{
		Name: p.cookieName,
		Value: id,  // sess_id
	}
	p.cookieOption(c)
	http.SetCookie(writer, c)
	return nil
}

func (p *Propagator) Extract(req *http.Request) (string, error) {
	cookie, err := req.Cookie(p.cookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func (p *Propagator) Remove(writer http.ResponseWriter) error {
	// 只需将 cookie 的 MaxAge 设为负数, 重新放入进行
	cookie := &http.Cookie{
		Name: p.cookieName,
		MaxAge: -1,
	}
	http.SetCookie(writer, cookie )
	return nil
}



