package responsecache

import (
	"net/http"
	"testing"

	"github.com/herb-go/herb/middleware/httpinfo"

	"github.com/herb-go/deprecated/cache"
)

type testID struct {
}

func (i *testID) IdentifyRequest(r *http.Request) (string, error) {
	return "testid", nil
}
func TestParamsWith(t *testing.T) {
	cb := NewParamsContextBuilder()
	ctx := NewContext()
	cb.BuildContext(ctx)
	if ctx.TTL != 0 || ctx.Validator != nil || ctx.Cache != nil {
		t.Fatal(ctx)
	}
	ttlcb := cb.WithTTL(1800)
	ctx = NewContext()
	ttlcb.BuildContext(ctx)
	if ctx.TTL != 1800 || ctx.Validator != nil || ctx.Cache != nil {
		t.Fatal(ctx)
	}
	valttlcb := ttlcb.WithValidator(httpinfo.ValidatorAlways)
	ctx = NewContext()
	valttlcb.BuildContext(ctx)
	if ctx.TTL != 1800 || ctx.Validator == nil || ctx.Cache != nil {
		t.Fatal(ctx)
	}
	c := newTestCache(1800)
	cachevalttlcb := valttlcb.WithCache(c)
	ctx = NewContext()
	cachevalttlcb.BuildContext(ctx)
	if ctx.TTL != 1800 || ctx.Validator == nil || ctx.Cache != c {
		t.Fatal(ctx)
	}
	cachevalttlcbprefix := cachevalttlcb.WithCachePrefix("prefix")
	ctx = NewContext()
	cachevalttlcbprefix.BuildContext(ctx)
	if ctx.TTL != 1800 || ctx.Validator == nil || ctx.Cache != c || ctx.CachePrefix != "prefix" {
		t.Fatal(ctx)
	}

	ctx = NewContext()
	cb.BuildContext(ctx)
	if ctx.TTL != 0 || ctx.Validator != nil || ctx.Cache != nil {
		t.Fatal(ctx)
	}
}

func TestParams(t *testing.T) {

	p := NewParams()
	p.Append(ParamFunc(func(r *http.Request) (string, bool) {
		return "test", true
	}))
	s, ok := p.GetParam(nil)
	if s != "test" || ok != true {
		t.Fatal(s)
	}
	p.Append(ParamFunc(func(r *http.Request) (string, bool) {
		return "test2", true
	}))
	s, ok = p.GetParam(nil)
	if s != "test"+cache.KeyPrefix+"test2" || ok != true {
		t.Fatal(s)
	}
	p.Append(ParamFunc(func(r *http.Request) (string, bool) {
		return "", false
	}))
	s, ok = p.GetParam(nil)
	if s != "" || ok != false {
		t.Fatal(s)
	}
}
