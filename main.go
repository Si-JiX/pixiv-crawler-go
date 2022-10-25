package main

import (
	"fmt"
	"gopkg.in/urfave/cli.v1"
	"log"
	"os"
	"pixiv-cil/config"
	"pixiv-cil/download"
	"pixiv-cil/pixiv"
	"pixiv-cil/utils"
)

func init() {
	config.VarsConfigInit()
	if config.Vars.PixivRefreshToken == "" {
		PixivRefreshToken, err := pixiv.ChromeDriverLogin()
		if err != nil {
			panic(err)
		}
		if token, ok := pixiv.InitAuth(PixivRefreshToken); ok != nil {
			fmt.Println("refresh token is invalid,please login again:", ok)
		} else {
			config.VarsFile.Vipers.Set("PIXIV_REFRESH_TOKEN", PixivRefreshToken)
			config.VarsFile.Vipers.Set("PIXIV_TOKEN", token)
			config.VarsFile.SaveConfig()
		}
	}
	pixiv.TokenVariable = config.Vars.PixivToken
	pixiv.RefreshTokenVariable = config.Vars.PixivRefreshToken
}

var implement = func(c *cli.Context) error {
	if config.CommandLines.IllustID != 0 {
		download.CurrentDownloader(config.CommandLines.IllustID)
	} else if config.CommandLines.AuthorID != 0 {
		download.AuthorImageALL(config.CommandLines.AuthorID)
	} else if config.CommandLines.URL != "" {
		download.CurrentDownloader(utils.GetInt(config.CommandLines.URL))
	} else {
		_ = cli.ShowAppHelp(c)
	}
	return nil
}

//func ShellUserDetail() {
//	user, err :=  App.UserDetail(36966416) // print user detail information, exclude illusts collections
//	if err != nil {
//		fmt.Println(err)
//	} else {
//		//fmt.Println(user.Profile)
//		//fmt.Println(user.ProfilePublicity)
//		str, err := json.Marshal(user)
//		if err != nil {
//			fmt.Println(err)
//		} else {
//			fmt.Println(string(str))
//		}
//		fmt.Println(user.Workspace)
//		//fmt.Println(user.User.Name)
//		//fmt.Println(user.User.ProfileImages)
//	}
//}

func main() {
	config.NewFile("imageFile")
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
