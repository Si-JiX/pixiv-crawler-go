package config

import (
	"fmt"
	"gopkg.in/urfave/cli.v1"
	pixiv "pixiv-cil/pixiv_api"
	"pixiv-cil/utils"
	"strings"
)

var CommandLines = struct {
	IllustID int
	AuthorID int
	Name     string
	URL      string
}{}

var CommandLineFlag = []cli.Flag{
	cli.IntFlag{
		Name:        "d, download",
		Value:       0,
		Usage:       "input IllustID to download",
		Destination: &CommandLines.IllustID,
	},
	cli.StringFlag{
		Name:        "u, url",
		Value:       "",
		Usage:       "input pixiv url to download",
		Destination: &CommandLines.URL,
	},
	cli.StringFlag{
		Name:        "n, name",
		Value:       "",
		Usage:       "author name",
		Destination: &CommandLines.Name,
	},
	cli.IntFlag{
		Name:        "a, author",
		Value:       0,
		Usage:       "author id",
		Destination: &CommandLines.AuthorID,
	},
}

func INIT_PIXIV_AUTH() *pixiv.AppPixivAPI {
	if f := Open("author_key.txt", "r"); f != nil {
		utils.PIXAPI_TOKEN_KEY = strings.Split(string(f), "\n")[0]
		utils.PIXAPI_RE_TOKEN_KEY = strings.Split(string(f), "\n")[1]
		account, state := pixiv.LoadAuth(utils.PIXAPI_TOKEN_KEY, utils.PIXAPI_RE_TOKEN_KEY, utils.PIXAPI_TOKEN_TIME_KEY)
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
