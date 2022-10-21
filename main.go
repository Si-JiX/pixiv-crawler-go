package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/urfave/cli.v1"
	"log"
	"os"
	"pixiv-cil/config"
	"pixiv-cil/download"
)

var implement = func(c *cli.Context) error {
	if config.CommandLines.IllustID != 0 {
		siz, err := config.App.Download(uint64(config.CommandLines.IllustID), "imageFile")
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("download success\tsize:", siz)
		}
	} else if config.CommandLines.AuthorID != 0 {
		download.GET_AUTHOR(uint64(config.CommandLines.AuthorID), 0)
	} else if config.CommandLines.IllustID != 0 {
		//FindID := regexp.MustCompile(`(\d+)`).FindAllString(config.CommandLines.URL, -1)
		//if FindID != nil {
		//	download.GET_IMAGE_INFO(FindID[0])
		//} else {
		//	fmt.Println("url error", config.CommandLines.URL)
		//}
	} else {
		_ = cli.ShowAppHelp(c)
	}
	return nil
}

func init() {
	config.INIT_PIXIV_AUTH() // init pixiv auth
	cli_app := cli.NewApp()
	cli_app.Name = "image downloader"
	cli_app.Version = "V.1.0.9"
	cli_app.Usage = "download image from pixiv "
	cli_app.Flags = config.CommandLineFlag
	cli_app.Action = implement
	if err := cli_app.Run(os.Args); err != nil {
		log.Fatal(err)
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
	config.NewFile("imageFile")
	//init_command()
	//for i, illust := range config.ImageList {
	//	fmt.Println(i, illust.Title)
	//}

	//illusts, next, err := app.UserBookmarksIllust(uid, "public", 0, "")
	//illusts, next, err := app.IllustFollow("public", 0)
}
