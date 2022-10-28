package request

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Request struct {
	Path     string
	Mode     string
	Header   map[string]string
	Query    map[string]string
	Params   url.Values
	requests *http.Request
}

type Response struct {
	Response *http.Response
	Request  *Request
	Body     io.ReadCloser
	content  []byte
	text     string
}

func (req *Request) NewRequest() *Response {
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
	req.Headers()
	if response, err = http.DefaultClient.Do(req.requests); err != nil {
		fmt.Println("http.DefaultClient.Do error:", err)
		return nil
	} else {
		return &Response{Response: response, Request: req, Body: response.Body}
	}
}

func Get(url_api string, params map[string]string) *Response {
	req := &Request{Mode: "GET", Params: url.Values{}, Header: map[string]string{}, Path: url_api, Query: params}
	return req.NewRequest()
}

func Post(url_api string, params map[string]string) *Response {
	req := &Request{Mode: "POST", Params: url.Values{}, Header: map[string]string{}, Path: url_api, Query: params}
	return req.NewRequest()
}

func (resp *Response) Content() []byte {
	resp.content, _ = io.ReadAll(resp.Body)
	return resp.content
}

func (resp *Response) Text() string {
	resp.Content() //	Init resp.content
	resp.text = string(resp.content)
	return resp.text
}

func (resp *Response) Json(value interface{}) interface{} {
	resp.Content() //	Init resp.content
	if strings.Contains("OAuth", string(resp.content)) {
		fmt.Println("Token expired, Refreshing...")
		//RefreshAuth()
	}
	//RefreshAuth()
	if err := json.Unmarshal(resp.content, value); err != nil {
		fmt.Println("json.Unmarshal error:", err)
	}
	return value
}

func (resp *Response) GetCookies() (cookies []*http.Cookie) {
	return resp.Response.Request.Cookies()

}
