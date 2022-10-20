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
	"time"
)

var CommandLines = struct {
	URL  string
	Name string
}{}

var commandLines = []cli.Flag{
	cli.StringFlag{
		Name:        "d, download",
		Value:       "",
		Usage:       "url of image",
		Destination: &CommandLines.URL,
	},
	cli.StringFlag{
		Name:        "n, name",
		Value:       "",
		Usage:       "author name",
		Destination: &CommandLines.Name,
	},
}
var implement = func(c *cli.Context) error {
	if CommandLines.URL != "" {
		FindID := regexp.MustCompile(`(\d+)`).FindAllString(CommandLines.URL, -1)
		if FindID != nil {
			download.GET_IMAGE_INFO(FindID[0])
		} else {
			fmt.Println("url error", CommandLines.URL)
		}
	} else if CommandLines.Name != "" {
		download.GET_AUTHOR(CommandLines.Name, 1)
		fmt.Println(len(config.ImageList), "images found")
		for index, value := range config.ImageList {
			fmt.Println("index:", index, "\t\ttitle:", value.Title, "\t\tid:", value.ID)
		}
		for _, value := range config.ImageList {
			download.ImageDownloader(value.ImageUrls.Large, value.Title)
		}
	} else {
		_ = cli.ShowAppHelp(c)
	}
	return nil
}

func init_command() {
	app := cli.NewApp()
	app.Name = "image downloader"
	app.Version = "V.1.0.0"
	app.Usage = ""
	app.Flags = commandLines
	app.Action = implement
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}
func init() {
	var PIXAPI_TOKEN_KEY = ""
	var PIXAPI_RE_TOKEN_KEY = ""
	var PIXAPI_TOKEN_TIME_KEY = time.Now()
	f, err := os.ReadFile("author_key.txt")
	if err == nil {
		PIXAPI_TOKEN_KEY = strings.Split(string(f), "\n")[0]
		PIXAPI_RE_TOKEN_KEY = strings.Split(string(f), "\n")[1]
		fmt.Println(PIXAPI_TOKEN_KEY)
		fmt.Println(PIXAPI_RE_TOKEN_KEY)
		account, ok := pixiv.LoadAuth(PIXAPI_TOKEN_KEY, PIXAPI_RE_TOKEN_KEY, PIXAPI_TOKEN_TIME_KEY)
		if ok != nil {
			fmt.Println(ok)
		} else {
			fmt.Println("account:", account.Account)
		}
	} else {
		panic(err)
	}

}

func main() {
	//if !config.IsExist("imageFile") {
	//	_ = os.Mkdir("imageFile", 0777)
	//}
	//init_command()

	app := pixiv.NewApp()
	user, err := app.UserDetail(36966416) // print user detail information, exclude illusts collections
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
	//illusts, next, err := app.UserIllusts(uid, "illust", 0)
	//illusts, next, err := app.UserBookmarksIllust(uid, "public", 0, "")
	//illusts, next, err := app.IllustFollow("public", 0)
}
