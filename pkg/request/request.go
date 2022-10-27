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
	requests *http.Request
	result   []byte
	Params   url.Values
}

type Response struct {
	Response *http.Response
	content  []byte
	text     string
	Request  *Request
}

func (req *Request) QueryData() io.ReadCloser {
	return io.NopCloser(strings.NewReader(req.EncodeParams()))
}

func (req *Request) EncodeParams(datas ...map[string]string) string {
	for _, data := range datas {
		for key, value := range data {
			req.Params.Add(key, value)
		}
	}
	return req.Params.Encode()
}

func (req *Request) Headers() {
	for k, v := range map[string]string{} {
		req.requests.Header.Set(k, v)
	}
}
func (req *Request) Get(url_api string, params map[string]string) (resp *Response, err error) {
	if params != nil {
		url_api = url_api + "?" + req.EncodeParams(params)
	}
	req.requests, err = http.NewRequest("GET", url_api, nil)
	req.Headers()
	response, err := http.DefaultClient.Do(req.requests)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &Response{Response: response, Request: req}, nil
}

func (req *Request) Post(url_api string) (resp *Response, err error) {
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
	resp.content, _ = io.ReadAll(resp.Request.requests.Body)
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
