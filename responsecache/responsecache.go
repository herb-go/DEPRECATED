package responsecache

import "net/http"

func (c ContextField) NewResponseCache(b ContextBuilder) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		ctx := c.GetContext(r)
		b.BuildContext(ctx)
		c.ServeMiddleware(w, r, next)
	}
}

func New(b ContextBuilder) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return DefaultContextField.NewResponseCache(b)
}
