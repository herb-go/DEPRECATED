package responsecache

import (
	"net/http"
	"testing"
)

func TestContext(t *testing.T) {
	r, _ := http.NewRequest("GET", "127.0.0.1", nil)
	ctx := DefaultContextField.GetContext(r)
	ctx.ID = "test"
	ctx2 := DefaultContextField.GetContext(r)
	if ctx.ID != ctx2.ID {
		t.Fatal(ctx2)
	}
}
