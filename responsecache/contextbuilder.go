package responsecache

import (
	"time"

	"github.com/herb-go/deprecated/cache"
	"github.com/herb-go/herb/middleware/httpinfo"
)

type ContextBuilder interface {
	BuildContext(*Context)
}

type PlainContextBuilder struct {
	ID          string
	Validator   httpinfo.Validator
	TTL         time.Duration
	Cache       cache.Cacheable
	CachePrefix string
	CachePath   []string
}

func (b *PlainContextBuilder) BuildContext(ctx *Context) {
	ctx.ID = b.ID
	ctx.Validator = b.Validator
	ctx.TTL = b.TTL
	ctx.Cache = b.Cache
	ctx.CachePrefix = b.CachePrefix
	ctx.CachePath = b.CachePath
}

func NewPlainContextBuilder() *PlainContextBuilder {
	return &PlainContextBuilder{}
}
