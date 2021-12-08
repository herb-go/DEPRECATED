package responsecache

import (
	"net/http"
)

type cached struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

func mergeHeader(src http.Header, dst *http.Header) {
	for key := range src {
		for k := range src[key] {
			(*dst).Add(key, src[key][k])
		}
	}
}

func (c *cached) Output(w http.ResponseWriter) error {
	h := w.Header()
	mergeHeader(c.Header, &h)
	w.WriteHeader(c.StatusCode)
	_, err := w.Write(c.Body)
	return err
}

func cacheContext(ctx *Context) *cached {
	data := ctx.Response.UncommittedData()
	c := &cached{
		StatusCode: ctx.Response.StatusCode,
		Header:     http.Header{},
		Body:       data,
	}
	mergeHeader(ctx.Response.Header(), &c.Header)
	return c
}
