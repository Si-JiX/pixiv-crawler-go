package request

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Request struct {
	Path     string            // API Path
	Mode     string            // GET, POST, PUT
	Header   map[string]string // Request Header
	Query    map[string]string // Query Params
	Params   url.Values        // init in url.Values
	requests *http.Request     // init in http.NewRequest
}

type Response struct {
	Response *http.Response // Response from http.DefaultClient.Do
	Request  *Request       // Request from type Request
	Body     io.ReadCloser  // Body from Response
	content  []byte         // Body -> Content []byte
	text     string         // Content -> string
}

func Get(url_api string, params map[string]string, Head ...map[string]string) *Response {
	req := &Request{Mode: "GET", Params: url.Values{}, Header: map[string]string{}, Path: url_api, Query: params}
	return req.NewRequest(Head...)
}

func Post(url_api string, params map[string]string, Head ...map[string]string) *Response {
	req := &Request{Mode: "POST", Params: url.Values{}, Header: map[string]string{}, Path: url_api, Query: params}
	return req.NewRequest(Head...)
}
func Put(url_api string, params map[string]string, Head ...map[string]string) *Response {
	req := &Request{Mode: "PUT", Params: url.Values{}, Header: map[string]string{}, Path: url_api, Query: params}
	return req.NewRequest(Head...)
}

func (req *Request) NewRequest(Head ...map[string]string) *Response {
	var err error
	var body io.ReadCloser
	var response *http.Response
	if req.Mode == "POST" && req.Query != nil {
		body = req.QueryData()
	} else if req.Mode == "GET" && req.Query != nil {
		req.Path = req.Path + "?" + req.EncodeParams(req.Query)
	}
	req.requests, err = http.NewRequest(req.Mode, req.Path, body)
	if err != nil {
		fmt.Println("http.NewRequest error:", err)
		return nil
	}
	for _, data := range Head {
		for key, value := range data {
			req.Header[key] = value
		}
	}
	req.Headers()
	if response, err = http.DefaultClient.Do(req.requests); err != nil {
		fmt.Println("http.DefaultClient.Do error:", err)
		return nil
	} else {
		return &Response{Response: response, Request: req, Body: response.Body}
	}
}
func (resp *Response) Content() []byte {
	resp.content, _ = io.ReadAll(resp.Body)
	if strings.Contains(string(resp.content), "Token") {
		fmt.Println("Token expired, Refreshing...")
		RefreshAuth()
		resp.content = resp.Request.NewRequest().Content()
		time.Sleep(2 * time.Second)
	}
	return resp.content
}

func (resp *Response) GetBody() io.ReadCloser {
	return resp.Body
}

func (resp *Response) Text() string {
	resp.Content() //	Init resp.content
	resp.text = string(resp.content)
	return resp.text
}

func (resp *Response) Json(value interface{}) interface{} {
	resp.Content() //	Init resp.content
	if err := json.Unmarshal(resp.content, value); err != nil {
		fmt.Println("json.Unmarshal error:", err)
	}
	return value
}

func (resp *Response) GetCookies() (cookies []*http.Cookie) {
	return resp.Response.Request.Cookies()

}
