package config

import (
	pixiv "pixiv-cil/pixiv_api"
	"time"
)

var ImageList []pixiv.Illust
var App *pixiv.AppPixivAPI
var PIXAPI_TOKEN_KEY = ""
var PIXAPI_RE_TOKEN_KEY = ""
var PIXAPI_TOKEN_TIME_KEY = time.Now()
