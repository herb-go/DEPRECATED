package responsecache

import (
	"net/http"
	"strings"
	"time"

	"github.com/herb-go/deprecated/cache"
	"github.com/herb-go/herb/middleware/httpinfo"
)

type ParamFunc func(r *http.Request) (param string, success bool)

func (p ParamFunc) GetParam(r *http.Request) (param string, success bool) {
	return p(r)
}

type Param interface {
	GetParam(r *http.Request) (param string, success bool)
}

type Params []Param

func (p *Params) GetParam(r *http.Request) (string, bool) {
	results := make([]string, len(*p))
	for k := range *p {
		s, ok := (*p)[k].GetParam(r)
		if ok == false {
			return "", false
		}
		results[k] = s
	}
	return strings.Join(results, cache.KeyPrefix), true
}

func (p *Params) Clone() *Params {
	params := make(Params, len(*p))
	copy(params, *p)
	return &params
}

func (p *Params) Append(newparams ...Param) *Params {
	*p = append(*p, newparams...)
	return p
}
func NewParams() *Params {
	return &Params{}
}

type ParamsList []*Params

func (p *ParamsList) Clone() *ParamsList {
	pl := make(ParamsList, len(*p))
	copy(pl, *p)
	return &pl
}

func (p *ParamsList) Append(newparamslist ...*Params) *ParamsList {
	*p = append(*p, newparamslist...)
	return p
}

func (p *ParamsList) GetParams(r *http.Request) ([]string, bool) {
	var ok bool
	results := make([]string, len(*p))
	for k := range *p {
		results[k], ok = (*p)[k].GetParam(r)
		if !ok {
			return nil, false
		}
	}
	return results, true
}
func NewParamsList() *ParamsList {
	return &ParamsList{}
}

type ParamsContextBuilder struct {
	params     *Params
	pathparams *ParamsList
	ttl        time.Duration
	validator  httpinfo.Validator
	cache      cache.Cacheable
	prefix     string
}

func (b *ParamsContextBuilder) BuildContext(ctx *Context) {
	ctx.ID, _ = b.params.GetParam(ctx.Request)
	if ctx.ID != "" {
		p, ok := b.pathparams.GetParams(ctx.Request)
		if ok {
			ctx.CachePath = p
		} else {
			ctx.ID = ""
		}
	}
	ctx.Validator = b.validator
	ctx.TTL = b.ttl
	ctx.Cache = b.cache
	ctx.CachePrefix = b.prefix
}
func (b *ParamsContextBuilder) Clone() *ParamsContextBuilder {
	return &ParamsContextBuilder{
		params:     b.params.Clone(),
		pathparams: b.pathparams.Clone(),
		validator:  b.validator,
		ttl:        b.ttl,
		cache:      b.cache,
		prefix:     b.prefix,
	}
}
func (b *ParamsContextBuilder) WithTTL(ttl time.Duration) *ParamsContextBuilder {
	pcb := b.Clone()
	pcb.ttl = ttl
	return pcb
}
func (b *ParamsContextBuilder) WithCache(c cache.Cacheable) *ParamsContextBuilder {
	pcb := b.Clone()
	pcb.cache = c
	return pcb
}
func (b *ParamsContextBuilder) WithCachePrefix(prefix string) *ParamsContextBuilder {
	pcb := b.Clone()
	pcb.prefix = prefix
	return pcb
}

func (b *ParamsContextBuilder) AppendParams(params ...Param) *ParamsContextBuilder {
	pcb := b.Clone()
	pcb.params.Append(params...)
	return pcb
}

func (b *ParamsContextBuilder) AppendPathParams(params ...*Params) *ParamsContextBuilder {
	pcb := b.Clone()
	pcb.pathparams.Append(params...)
	return pcb
}

func (b *ParamsContextBuilder) WithValidator(v httpinfo.Validator) *ParamsContextBuilder {
	pcb := b.Clone()
	pcb.validator = v
	return pcb
}
func NewParamsContextBuilder() *ParamsContextBuilder {
	return &ParamsContextBuilder{
		params:     NewParams(),
		pathparams: NewParamsList(),
	}
}
