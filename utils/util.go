package utils

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func SETHeaders(url string, req *http.Request) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	if strings.Contains(url, "pixiv") || strings.Contains(url, "pximg") {
		req.Header.Set("Referer", "https://www.pixiv.net/")
	}

}

func Request(url string) ([]byte, error) {
	if response, err := http.NewRequest("GET", url, nil); err == nil {
		SETHeaders(url, response)
		resp, ok := http.DefaultClient.Do(response)
		if ok == nil && resp.StatusCode == 200 {
			return io.ReadAll(resp.Body)
		}
	} else {
		fmt.Println("error", err)
	}
	return nil, nil
}
