package config

import (
	"fmt"
	"pixiv-cil/pixiv"
	"strings"
	"time"
)

var App *pixiv.AppPixivAPI

func INIT_PIXIV_AUTH() *pixiv.AppPixivAPI {
	var PIXAPI_TOKEN_KEY = ""
	var PIXAPI_RE_TOKEN_KEY = ""
	var PIXAPI_TOKEN_TIME_KEY = time.Now()
	if f := Open("author_key.txt", "r"); f != nil {
		PIXAPI_TOKEN_KEY = strings.Split(string(f), "\n")[0]
		PIXAPI_RE_TOKEN_KEY = strings.Split(string(f), "\n")[1]
		account, state := pixiv.LoadAuth(PIXAPI_TOKEN_KEY, PIXAPI_RE_TOKEN_KEY, PIXAPI_TOKEN_TIME_KEY)
		if state != nil {
			fmt.Println(state)
		} else {
			fmt.Println("you account is valid\taccount name:", account.Account)
		}
	} else {
		fmt.Println("you need login pixiv account first")
	}
	return pixiv.NewApp()
}
