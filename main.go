package main

import (
	"gopkg.in/urfave/cli.v1"
	"log"
	"os"
	"pixiv-cil/pkg/command"
	"pixiv-cil/pkg/config"
	"pixiv-cil/pkg/file"
	"pixiv-cil/src"
	"pixiv-cil/src/pixiv"
	"pixiv-cil/utils"
)

func init() {
	config.VarsConfigInit()
	if config.Vars.PixivRefreshToken == "" {
		if accessToken, err := pixiv.ChromeDriverLogin(); err != nil {
			panic(err)
		} else {
			config.VarsFile.Vipers.Set("PIXIV_REFRESH_TOKEN", accessToken.RefreshToken)
			config.VarsFile.Vipers.Set("PIXIV_TOKEN", accessToken.AccessToken)
			config.VarsFile.Vipers.Set("PIXIV_USER_ID", accessToken.User.ID)
			config.VarsFile.SaveConfig()
		}
	}
	pixiv.TokenVariable = config.Vars.PixivToken
	pixiv.RefreshTokenVariable = config.Vars.PixivRefreshToken
}

func main() {
	file.NewFile("imageFile")
	cli_app := cli.NewApp()
	cli_app.Name = "image downloader"
	cli_app.Version = config.Vars.VersionName
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
