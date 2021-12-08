package responsecache

import (
	"testing"

	"github.com/herb-go/herb/middleware/httpinfo"
)

func TestPlainContextBuilder(t *testing.T) {
	b := NewPlainContextBuilder()
	b.ID = "test"
	b.CachePath = []string{"path"}
	b.Validator = httpinfo.ValidatorAlways
	c := newTestCache(3600)
	b.Cache = c
	b.TTL = 1800
	b.CachePrefix = "prefix"
	ctx := NewContext()
	b.BuildContext(ctx)
	ok, err := ctx.Validator.Validate(nil, nil)
	if err != nil {
		panic(err)
	}
	if ctx.ID != "test" || len(ctx.CachePath) != 1 || ctx.CachePath[0] != "path" || ok != true || ctx.Cache != c || ctx.TTL != 1800 || ctx.CachePrefix != "prefix" {
		t.Fatal(ctx)
	}
}
