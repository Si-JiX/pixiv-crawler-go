package request

import (
	"github.com/VeronicaAlexia/pixiv-crawler-go/pkg/config"
)

func (req *Request) AddHeader(key string, value string) {
	req.Header[key] = value
}

func (req *Request) Headers() {
	// Set default headers for request
	req.AddHeader("User-Agent", "PixivIOSApp/7.6.2 (iOS 12.2; iPhone9,1)")
	req.AddHeader("App-OS", "ios")
	req.AddHeader("App-OS-VERSION", "12.2")
	req.AddHeader("App-Version", "7.6.2")
	req.AddHeader("Authorization", "Bearer "+config.Vars.PixivToken)
	if req.Mode == "POST" {
		req.AddHeader("Content-Type", "application/x-www-form-urlencoded")
	} else {
		req.AddHeader("Content-Type", "application/json")
	}
	for k, v := range req.Header {
		req.requests.Header.Set(k, v)
	}
}
