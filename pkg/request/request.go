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
	requests *http.Request
	Header   map[string]string
	Params   url.Values
}

type Response struct {
	Response *http.Response
	Request  *Request
	Body     io.ReadCloser
	content  []byte
	text     string
}

func Get(url_api string, params map[string]string) *Response {
	req := &Request{Params: url.Values{}, Header: map[string]string{}, Path: url_api}
	if params != nil {
		url_api = url_api + "?" + req.EncodeParams(params)
	}
	req.requests, _ = http.NewRequest("GET", url_api, nil)
	req.Headers()
	if response, err := http.DefaultClient.Do(req.requests); err != nil {
		fmt.Println(err)
	} else {
		return &Response{Response: response, Request: req, Body: response.Body}
	}
	return nil
}

func Post(url_api string, req *Request) *Response {
	req.requests, _ = http.NewRequest("POST", url_api, req.QueryData())
	req.Headers()
	req.requests.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//fmt.Println(req.requests.Header)
	if response, err := http.DefaultClient.Do(req.requests); err != nil {
		return nil
	} else {
		return &Response{Response: response, Request: req, Body: response.Body}
	}
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
