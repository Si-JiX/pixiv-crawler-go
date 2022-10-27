package request

func (req *Request) AddHeader(key string, value string) {
	req.Header[key] = value
}

func (req *Request) Headers() {
	for k, v := range req.Header {
		req.requests.Header.Set(k, v)
	}
}
