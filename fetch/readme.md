# Fetch

一个便于使用和可配置化的网络接口调取库

## 状态

功能冻结，仅保证老项目使用，不推荐新项目继续使用



## 功能
* 易于序列化(JSON,TOML)的配置结构
* 支持代理设置
* 支持自动将参数作为JSON/XML序列化
* 支持自动加载响应正文并作为JSON/XML解析
* 响应正文实现error接口
* 支持将响应正文转化为带code值的错误
* 通过Server和EndPoint结构快速的快速定义api接口

## 范例代码

### JSON配置

config.json

    {
        "Clients":{
            	"TimeoutInSecond":120,
	            "MaxIdleConns":20,
	            "IdleConnTimeoutInSecond":120,
	            "TLSHandshakeTimeoutInSecond":30,
	            "ProxyURL":"socks5://127.0.0.1:2345"
        }
    }

main.go

    package main

    import (
        "encoding/json"
        "fmt"
        "io/ioutil"
        "net/http"

        "github.com/herb-go/deprecated/fetch"
    )

    type Config struct {
        Clients fetch.Clients
    }

    func main() {
        bs, err := ioutil.ReadFile("./config.json")
        if err != nil {
            panic(err)
        }
        config := &Config{}
        err = json.Unmarshal(bs, config)
        if err != nil {
            panic(err)
        }
        req, err := http.NewRequest("GET", "https://www.facebook.com", nil)
        if err != nil {
            panic(err)
        }
        resp, err := config.Clients.Fetch(req)
        if err != nil {
            panic(err)
        }
        fmt.Println(string(resp.BodyContent))
    }

### Server与EndPoint

    package main

    import (
        "fmt"
        "net/url"

        "github.com/herb-go/deprecated/fetch"
    )

    type result struct {
        Visibility int `json:"visibility"`
    }

    func main() {
        clients := &fetch.Clients{}
        server := &fetch.Server{
            Host: "http://samples.openweathermap.org",
        }
        APIWeather := server.EndPoint("GET", "/data/2.5/weather")
        p := url.Values{}
        p.Add("q", "London,uk")
        p.Add("appid", "b6907d289e10d714a6e88b30761fae22")
        req, err := APIWeather.NewRequest(p, nil)
        if err != nil {
            panic(err)
        }
        resp, err := clients.Fetch(req)
        if err != nil {
            panic(err)
        }
        r := &result{}
        err = resp.UnmarshalAsJSON(r)
        if err != nil {
            panic(err)
        }
        fmt.Println(r.Visibility)
    }

### 响应作为error
    package main

    import (
        "fmt"

        "github.com/herb-go/deprecated/fetch"
    )

    func main() {
        clients := &fetch.Clients{}
        server := &fetch.Server{
            Host: "http://www.example.com",
        }
        API := server.EndPoint("GET", "/api")
        req, err := API.NewRequest(nil, nil)
        if err != nil {
            panic(err)
        }
        resp, err := clients.Fetch(req)
        if err != nil {
            panic(err)
        }
        fmt.Println(fetch.GetErrorStatusCode(resp))
        panic(resp)
    }


### 响应作为带code的error

    package main

    import (
        "fmt"

        "github.com/herb-go/deprecated/fetch"
    )

    func main() {
        clients := &fetch.Clients{}
        server := &fetch.Server{
            Host: "http://www.example.com",
        }
        API := server.EndPoint("GET", "/api")
        req, err := API.NewRequest(nil, nil)
        if err != nil {
            panic(err)
        }
        resp, err := clients.Fetch(req)
        if err != nil {
            panic(err)
        }
        e := resp.NewAPICodeErr(resp.StatusCode)
        fmt.Println(fetch.GetAPIErrCode(e))
        fmt.Println(fetch.CompareAPIErrCode(e, 404))

        panic(e)
    }