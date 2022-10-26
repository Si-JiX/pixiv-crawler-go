package main

import (
	"fmt"
	"gopkg.in/urfave/cli.v1"
	"log"
	"os"
	"pixiv-cil/pkg/command"
	"pixiv-cil/pkg/config"
	"pixiv-cil/pkg/file"
	"pixiv-cil/src"
	pixiv2 "pixiv-cil/src/pixiv"
	"pixiv-cil/utils"
)

func init() {
	config.VarsConfigInit()
	if config.Vars.PixivRefreshToken == "" {
		PixivRefreshToken, err := pixiv2.ChromeDriverLogin()
		if err != nil {
			panic(err)
		}
		if token, ok := pixiv2.InitAuth(PixivRefreshToken); ok != nil {
			fmt.Println("refresh token is invalid,please login again:", ok)
		} else {
			config.VarsFile.Vipers.Set("PIXIV_REFRESH_TOKEN", PixivRefreshToken)
			config.VarsFile.Vipers.Set("PIXIV_TOKEN", token)
			config.VarsFile.SaveConfig()
		}
	}
	pixiv2.TokenVariable = config.Vars.PixivToken
	pixiv2.RefreshTokenVariable = config.Vars.PixivRefreshToken
}

func main() {
	file.NewFile("imageFile")
	cli_app := cli.NewApp()
	cli_app.Name = "image downloader"
	cli_app.Version = "V.1.0.9"
	cli_app.Usage = "download image from pixiv "
	cli_app.Flags = command.CommandLineFlag
	cli_app.Action = func(c *cli.Context) error {
		if command.CommandLines.IllustID != 0 {
			src.CurrentDownloader(command.CommandLines.IllustID)
		} else if command.CommandLines.AuthorID != 0 {
			src.AuthorImageALL(command.CommandLines.AuthorID)
		} else if command.CommandLines.URL != "" {
			src.CurrentDownloader(utils.GetInt(command.CommandLines.URL))
		} else {
			_ = cli.ShowAppHelp(c)
		}
		return nil
	}
	if err := cli_app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
