package request

import (
	"fmt"
	"pixiv-cil/pkg/config"
)

func (req *Request) AddHeader(key string, value string) {
	req.requests.Header.Set(key, value)
}

func (req *Request) Headers() {
	// Set default headers for request
	req.AddHeader("User-Agent", "PixivIOSApp/7.6.2 (iOS 12.2; iPhone9,1)")
	req.AddHeader("App-OS", "ios")
	req.AddHeader("App-OS-VERSION", "12.2")
	req.AddHeader("App-Version", "7.6.2")
	if config.Vars.PixivToken != "" {
		req.AddHeader("Authorization", "Bearer "+config.Vars.PixivToken)
	} else {
		fmt.Println("token is empty!")
	}
	// Set headers for request
	//for key, value := range req.Header {
	//	req.requests.Header.Set(key, value)
	//}
}
