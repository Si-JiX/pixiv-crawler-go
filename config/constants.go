package config

import (
	pixiv "pixiv-cil/pixiv_api"
)

var ImageList []pixiv.Illust
var App *pixiv.AppPixivAPI

var PIXAPI_TOKEN_KEY = ""
var PIXAPI_RE_TOKEN_KEY = ""
var PIXAPI_TOKEN_TIME_KEY = ""

const IntSize = 32 << (^uint(0) >> 63)
