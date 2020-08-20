package fetch

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
)

//Server api server struct
type Server struct {
	//Host server base host
	Host string
	//Headers headers which will be sent every request.
	Headers http.Header
}

//EndPoint create a new api server endpoint with given method and path
func (s *Server) EndPoint(method string, path string) *EndPoint {
	return &EndPoint{
		Server: s,
		Method: method,
		Path:   path,
	}
}

//BuildRequest builde given request and return any error raised
func (s *Server) BuildRequest(req *http.Request) error {
	if s.Headers != nil {
		for k, vs := range s.Headers {
			for _, v := range vs {
				req.Header.Add(k, v)
			}
		}
	}
	return nil
}

//RequestMethod return request method
func (s *Server) RequestMethod() string {
	return ""
}

//RequestURL return request url
func (s *Server) RequestURL() string {
	return s.Host
}

//RequestBody return request body
func (s *Server) RequestBody() io.Reader {
	return nil
}

//RequestBuilders return request builders
func (s *Server) RequestBuilders() []func(*http.Request) error {
	return []func(*http.Request) error{
		s.BuildRequest,
	}
}

//NewRequest create a new http.request with given method,path,params,and body.
func (s *Server) NewRequest(method string, path string, params url.Values, body []byte) (*http.Request, error) {
	u, err := url.Parse(s.Host + path)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	if params != nil {
		for k, vs := range params {
			for _, v := range vs {
				q.Add(k, v)
			}
		}
	}
	u.RawQuery = q.Encode()
	req, err := http.NewRequest(method, u.String(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	err = s.BuildRequest(req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

//NewJSONRequest create a new http.request with given method,path,params,and body encode by JSON.
func (s *Server) NewJSONRequest(method string, path string, params url.Values, v interface{}) (*http.Request, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return s.NewRequest(method, path, params, b)
}

//NewXMLRequest create a new http.request with given method,path,params,and body encode by XML.
func (s *Server) NewXMLRequest(method string, path string, params url.Values, v interface{}) (*http.Request, error) {
	b, err := xml.Marshal(v)
	if err != nil {
		return nil, err
	}
	return s.NewRequest(method, path, params, b)
}

//EndPoint api server endpoint struct
//Endpoint should be created by api server's EndPoint method
type EndPoint struct {
	Server *Server
	Path   string
	Method string
}

func (e *EndPoint) EndPoint(path string) *EndPoint {
	return &EndPoint{
		Server: e.Server,
		Path:   e.Path + path,
		Method: e.Method,
	}
}
func (e *EndPoint) BuildRequest(req *http.Request) error {
	return e.Server.BuildRequest(req)
}

//RequestMethod return request method
func (e *EndPoint) RequestMethod() string {
	return e.Method
}

//RequestURL return request url
func (e *EndPoint) RequestURL() string {
	return e.Server.Host + e.Path
}

//RequestBody return request body
func (e *EndPoint) RequestBody() io.Reader {
	return nil
}

//RequestBuilders return request builders
func (e *EndPoint) RequestBuilders() []func(*http.Request) error {
	return e.Server.RequestBuilders()
}

//NewRequest create a new http.request to end point with given params,and body.
func (e *EndPoint) NewRequest(params url.Values, body []byte) (*http.Request, error) {
	return e.Server.NewRequest(e.Method, e.Path, params, body)

}

//NewJSONRequest create a new http.request to end point with given params,and body encode by JSON.
func (e *EndPoint) NewJSONRequest(params url.Values, v interface{}) (*http.Request, error) {
	return e.Server.NewJSONRequest(e.Method, e.Path, params, v)
}

//NewXMLRequest create a new http.request to end point with given params,and body encode by XML.
func (e *EndPoint) NewXMLRequest(params url.Values, v interface{}) (*http.Request, error) {
	return e.Server.NewXMLRequest(e.Method, e.Path, params, v)
}
