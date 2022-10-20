package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/urfave/cli.v1"
	"log"
	"os"
	"pixiv-cil/config"
	"pixiv-cil/download"
	pixiv "pixiv-cil/pixiv_api"
	"regexp"
	"strings"
)

var implement = func(c *cli.Context) error {
	if config.CommandLines.URL != "" {
		FindID := regexp.MustCompile(`(\d+)`).FindAllString(config.CommandLines.URL, -1)
		if FindID != nil {
			download.GET_IMAGE_INFO(FindID[0])
		} else {
			fmt.Println("url error", config.CommandLines.URL)
		}
	} else if config.CommandLines.AuthorID != 0 {
		download.GET_AUTHOR(uint64(config.CommandLines.AuthorID), 0)
	} else {
		_ = cli.ShowAppHelp(c)
	}
	return nil
}

func init() {
	cli_app := cli.NewApp()
	cli_app.Name = "image downloader"
	cli_app.Version = "V.1.0.1"
	cli_app.Usage = ""
	cli_app.Flags = config.CommandLineFlag
	cli_app.Action = implement
	if err := cli_app.Run(os.Args); err != nil {
		log.Fatal(err)
	} else {
		config.App = pixiv.NewApp()
	}
	f, err := os.ReadFile("author_key.txt")
	if err == nil {
		config.PIXAPI_TOKEN_KEY = strings.Split(string(f), "\n")[0]
		config.PIXAPI_RE_TOKEN_KEY = strings.Split(string(f), "\n")[1]
		fmt.Println(config.PIXAPI_TOKEN_KEY)
		fmt.Println(config.PIXAPI_RE_TOKEN_KEY)
		account, ok := pixiv.LoadAuth(config.PIXAPI_TOKEN_KEY, config.PIXAPI_RE_TOKEN_KEY, config.PIXAPI_TOKEN_TIME_KEY)
		if ok != nil {
			fmt.Println(ok)
		} else {
			fmt.Println("account:", account.Account)
		}
	} else {
		panic(err)
	}

}

func ShellUserDetail() {
	user, err := config.App.UserDetail(36966416) // print user detail information, exclude illusts collections
	if err != nil {
		fmt.Println(err)
	} else {
		//fmt.Println(user.Profile)
		//fmt.Println(user.ProfilePublicity)
		str, err := json.Marshal(user)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(string(str))
		}
		fmt.Println(user.Workspace)
		//fmt.Println(user.User.Name)
		//fmt.Println(user.User.ProfileImages)
	}
}

func main() {
	//if !config.IsExist("imageFile") {
	//	_ = os.Mkdir("imageFile", 0777)
	//}
	//init_command()
	//for i, illust := range config.ImageList {
	//	fmt.Println(i, illust.Title)
	//}

	//illusts, next, err := app.UserBookmarksIllust(uid, "public", 0, "")
	//illusts, next, err := app.IllustFollow("public", 0)
}
