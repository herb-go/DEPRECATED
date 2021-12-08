package responsecache

import (
	"context"
	"net/http"
	"time"

	"github.com/herb-go/deprecated/cache"
	"github.com/herb-go/herb/middleware/httpinfo"
)

type ContextField string

func (c ContextField) GetContext(r *http.Request) *Context {
	var ctx *Context
	v := r.Context().Value(c)
	if v == nil {
		ctx = NewContext()
		reqctx := context.WithValue(r.Context(), c, ctx)
		req := r.WithContext(reqctx)
		ctx.Request = req
		*r = *req
	} else {
		ctx = v.(*Context)
	}
	return ctx
}

func (c ContextField) ServeMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	ctx := c.GetContext(r)
	if ctx.ID != "" && ctx.Cache != nil {
		page := &cached{}
		k := ctx.ID
		nc := cache.NewNestedCollection(ctx.Cache, ctx.CachePath...)
		if ctx.CachePrefix != "" {
			k = ctx.CachePrefix + cache.KeyPrefix + k
		}
		err := nc.Load(k, page, ctx.TTL, func(key string) (interface{}, error) {
			ctx.Prepare(w, r)
			next(ctx.Response.WrapWriter(w), r)
			if ctx.Response.Autocommit() {
				return nil, cache.ErrNotCacheable
			}
			page = cacheContext(ctx)
			return page, nil
		})
		if err != nil {
			if err != cache.ErrEntryTooLarge && err != cache.ErrNotCacheable {
				panic(err)
			}
			return
		}
		if ctx.Response.Autocommit() {
			return
		}
		err = page.Output(w)
		if err != nil {
			panic(err)
		}
		return
	}
	next(w, r)
}

var DefaultContextField = ContextField("responsecache")

type Context struct {
	Request     *http.Request
	ID          string
	TTL         time.Duration
	Validator   httpinfo.Validator
	Response    *httpinfo.Response
	Cache       cache.Cacheable
	CachePrefix string
	CachePath   []string
}

func NewContext() *Context {
	resp := httpinfo.NewResponse()
	resp.UpdateAutocommit(false)

	return &Context{
		Response: resp,
	}

}
func (c *Context) Prepare(w http.ResponseWriter, r *http.Request) {
	c.Request = r
	controller := httpinfo.NewCommitController(r, c.Response)
	if c.Validator == nil {
		controller.WithChecker(DefaultValidator)
	} else {
		controller.WithChecker(c.Validator)
	}
	c.Response.UpdateController(controller)
}
