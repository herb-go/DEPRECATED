package responsecache

import (
	"net/http"

	"github.com/herb-go/herb/middleware/httpinfo"
)

var DefaultValidator = httpinfo.ValidatorFunc(func(r *http.Request, resp *httpinfo.Response) (bool, error) {
	if !resp.Written {
		return true, nil
	}
	return resp.StatusCode >= 200 && resp.StatusCode < 500, nil
})
