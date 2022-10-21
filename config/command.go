package config

import (
	"fmt"
	"gopkg.in/urfave/cli.v1"
	"os"
	pixiv "pixiv-cil/pixiv_api"
	"strings"
)

var CommandLines = struct {
	IllustID int
	Name     string
	AuthorID int
}{}

var CommandLineFlag = []cli.Flag{
	cli.IntFlag{
		Name:        "d, download",
		Value:       0,
		Usage:       "input IllustID to download",
		Destination: &CommandLines.IllustID,
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

func INIT_PIXIV_AUTH() {
	App = pixiv.NewApp()
	if f, ok := os.ReadFile("author_key.txt"); ok == nil {
		PIXAPI_TOKEN_KEY = strings.Split(string(f), "\n")[0]
		PIXAPI_RE_TOKEN_KEY = strings.Split(string(f), "\n")[1]
		account, state := pixiv.LoadAuth(PIXAPI_TOKEN_KEY, PIXAPI_RE_TOKEN_KEY, PIXAPI_TOKEN_TIME_KEY)
		if state != nil {
			fmt.Println(state)
		} else {
			fmt.Println("you account is valid\taccount name:", account.Account)
		}
	} else {
		panic(ok)
	}
}
