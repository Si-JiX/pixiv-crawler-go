package request

import (
	"io"
	"strings"
)

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
