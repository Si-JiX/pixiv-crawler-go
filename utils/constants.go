package utils

import "sync"

var CurrentImageIndex int
var CurrentImageLength int

var ImageUrlList []string
var CH = make(chan struct{}, CHANNEL_CACHE)
var WG sync.WaitGroup

const IntSize = 32 << (^uint(0) >> 63)
const CHANNEL_CACHE = 20
const ClientID = "MOBrBDS8blbauoSck0ZfDbtuzpyT"
const ClientSecret = "lsACyCD94FhDUtGTXi3QzcFE2uU1hqtDaKeqrdwj"
const ClientHashSecret = "28c1fdd170a5204386cb1313c7077b34f83e4aaf4aa829ce78c231e05b0bae2c"
