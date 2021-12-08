package responsecache

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	_ "github.com/herb-go/deprecated/cache/drivers/syncmapcache"

	"github.com/herb-go/deprecated/cache"
)

func newTestCache(ttl int64) *cache.Cache {
	c := cache.New()
	oc := cache.NewOptionConfig()
	oc.Driver = "syncmapcache"
	oc.TTL = ttl * int64(time.Second)
	oc.Config = nil
	oc.Marshaler = "json"
	err := c.Init(oc)
	if err != nil {
		panic(err)
	}
	err = c.Flush()
	if err != nil {
		panic(err)
	}
	return c

}

var content int

func testAction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("ts", strconv.FormatInt(time.Now().UnixNano(), 10))
	w.WriteHeader(content)
	w.Write([]byte(r.Header.Get("test") + strconv.Itoa(content)))

}

var emptyParam = ParamFunc(func(r *http.Request) (string, bool) {
	return "", true
})

var keyParam = ParamFunc(func(r *http.Request) (string, bool) {
	return r.Header.Get("test"), true
})

var groupParams = NewParams().Append(
	ParamFunc(func(r *http.Request) (string, bool) {
		return r.Header.Get("group"), true
	}),
)

func TestResponseCache(t *testing.T) {
	c := newTestCache(3600)
	var rc = NewParamsContextBuilder().WithCache(c)
	content = 200
	mux := http.NewServeMux()
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("rawts", strconv.FormatInt(time.Now().UnixNano(), 10))
		New(rc.AppendParams(keyParam))(w, r, testAction)
	})
	mux.HandleFunc("/raw", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("rawts", strconv.FormatInt(time.Now().UnixNano(), 10))
		testAction(w, r)
	})
	mux.HandleFunc("/empty", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("rawts", strconv.FormatInt(time.Now().UnixNano(), 10))
		New(rc.AppendParams(emptyParam))(w, r, testAction)
	})
	mux.HandleFunc("/group", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("rawts", strconv.FormatInt(time.Now().UnixNano(), 10))
		New(rc.AppendParams(keyParam).AppendPathParams(groupParams))(w, r, testAction)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	req, err := http.NewRequest("GET", server.URL+"/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("test", "test1")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	ts := resp.Header.Get("ts")
	rawts := resp.Header.Get("rawts")

	resp.Body.Close()

	if string(bs) != "test1"+"200" {
		t.Error(string(bs))
	}
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}

	req, err = http.NewRequest("GET", server.URL+"/raw", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("test", "test1")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	bs, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	if string(bs) != "test1"+"200" {
		t.Error(string(bs))
	}
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}

	req, err = http.NewRequest("GET", server.URL+"/empty", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("test", "test1")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	bs, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	if string(bs) != "test1"+"200" {
		t.Error(string(bs))
	}
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}
	content = 404
	req, err = http.NewRequest("GET", server.URL+"/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("test", "test1")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	bs, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	if string(bs) != "test1"+"200" {
		t.Fatal(string(bs))
	}
	if resp.Header.Get("ts") != ts {
		t.Error(resp.Header.Get("ts"))
	}
	if resp.Header.Get("rawts") == rawts {
		t.Error(resp.Header.Get("rawts"))
	}
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}

	req, err = http.NewRequest("GET", server.URL+"/raw", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("test", "test1")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	bs, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	if string(bs) != "test1"+"404" {
		t.Error(string(bs))
	}
	if resp.Header.Get("ts") == ts {
		t.Error(resp.Header.Get("ts"))
	}
	if resp.Header.Get("rawts") == rawts {
		t.Error(resp.Header.Get("rawts"))
	}
	if resp.StatusCode != 404 {
		t.Error(resp.StatusCode)
	}

	req, err = http.NewRequest("GET", server.URL+"/empty", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("test", "test1")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	bs, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	if string(bs) != "test1"+"404" {
		t.Error(string(bs))
	}
	if resp.Header.Get("ts") == ts {
		t.Error(resp.Header.Get("ts"))
	}
	if resp.Header.Get("rawts") == rawts {
		t.Error(resp.Header.Get("rawts"))
	}
	if resp.StatusCode != 404 {
		t.Error(resp.StatusCode)
	}

	req, err = http.NewRequest("GET", server.URL+"/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("test", "test2")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	bs, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	if string(bs) != "test2"+"404" {
		t.Error(string(bs))
	}
	if resp.Header.Get("ts") == ts {
		t.Error(resp.Header.Get("ts"))
	}
	if resp.Header.Get("rawts") == rawts {
		t.Error(resp.Header.Get("rawts"))
	}
	if resp.StatusCode != 404 {
		t.Error(resp.StatusCode)
	}
	content = 500
	req, err = http.NewRequest("GET", server.URL+"/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("test", "test3")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	bs, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	if string(bs) != "test3"+"500" {
		t.Error(string(bs))
	}
	if resp.StatusCode != 500 {
		t.Error(resp.StatusCode)
	}
	content = 403
	req, err = http.NewRequest("GET", server.URL+"/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("test", "test3")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	bs, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	if string(bs) != "test3"+"403" {
		t.Error(string(bs))
	}
	if resp.StatusCode != 403 {
		t.Error(resp.StatusCode)
	}
	content = 200
	req, err = http.NewRequest("GET", server.URL+"/group", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("test", "grouptest1")
	req.Header.Set("group", "group1")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	bs, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	if string(bs) != "grouptest1"+"200" {
		t.Error(string(bs))
	}
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}
	content = 400
	cache.NewNestedCollection(c, "group2").Flush()
	req, err = http.NewRequest("GET", server.URL+"/group", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("test", "grouptest1")
	req.Header.Set("group", "group1")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	bs, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	if string(bs) != "grouptest1"+"200" {
		t.Error(string(bs))
	}
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}
	cache.NewNestedCollection(c, "group1").Flush()
	req, err = http.NewRequest("GET", server.URL+"/group", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("test", "grouptest1")
	req.Header.Set("group", "group1")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	bs, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	if string(bs) != "grouptest1"+"400" {
		t.Error(string(bs))
	}
	if resp.StatusCode != 400 {
		t.Error(resp.StatusCode)
	}
	content = 401
	cache.NewNestedCollection(c).Flush()
	req, err = http.NewRequest("GET", server.URL+"/group", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("test", "grouptest1")
	req.Header.Set("group", "group1")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	bs, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	if string(bs) != "grouptest1"+"401" {
		t.Error(string(bs))
	}
	if resp.StatusCode != 401 {
		t.Error(resp.StatusCode)
	}
}
