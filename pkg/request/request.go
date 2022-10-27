package request

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Request struct {
	requests *http.Request
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
	req := &Request{}
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

func Post(url_api string) (resp *Response, err error) {
	req := &Request{}
	req.requests, err = http.NewRequest("POST", url_api, nil)
	req.requests.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Headers()
	if response, err := http.DefaultClient.Do(req.requests); err != nil {
		return nil, err
	} else {

		return &Response{Response: response, Request: req}, nil
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
	if err := json.Unmarshal(resp.content, value); err != nil {
		fmt.Println("json.Unmarshal error:", err)
	}
	return value
}

func (resp *Response) GetCookies() (cookies []*http.Cookie) {
	return resp.Response.Request.Cookies()

}