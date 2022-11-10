package main

import (
	"fmt"
	"github.com/VeronicaAlexia/pixiv-crawler-go/pkg/arguments"
	"github.com/VeronicaAlexia/pixiv-crawler-go/pkg/cli"
	"github.com/VeronicaAlexia/pixiv-crawler-go/pkg/config"
	"github.com/VeronicaAlexia/pixiv-crawler-go/pkg/file"
	"github.com/VeronicaAlexia/pixiv-crawler-go/pkg/request"
	"github.com/VeronicaAlexia/pixiv-crawler-go/src"
	"github.com/VeronicaAlexia/pixiv-crawler-go/src/app"
	"github.com/VeronicaAlexia/pixiv-crawler-go/utils"
	"log"
	"os"
)

func init() {
	config.VarsConfigInit()
	if config.Vars.PixivRefreshToken == "" {
		if accessToken, err := request.ChromeDriverLogin(); err != nil {
			panic(err)
		} else {
			config.VarsFile.Vipers.Set("pixiv_refresh_token", accessToken.RefreshToken)
			config.VarsFile.Vipers.Set("pixiv_token", accessToken.AccessToken)
			config.VarsFile.Vipers.Set("PIXIV_USER_ID", accessToken.User.ID)
			config.VarsFile.SaveConfig()
		}
	} else {
		res, err := app.App.UserDetail(config.Vars.UserID)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("account: ", res.User.Name, "\tid: ", res.User.ID)
		}
	}
}

func main() {
	file.NewFile("imageFile")
	cli_app := cli.NewApp()
	cli_app.Name = "image downloader"
	cli_app.Version = config.Vars.VersionName
	cli_app.Usage = "download image from pixiv "
	cli_app.Flags = arguments.CommandLineFlag
	cli_app.Action = command_line_shell
	if err := cli_app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
func command_line_shell(c *cli.Context) error {
	if arguments.CommandLines.IllustID != "" {
		src.DownloaderSingly(arguments.CommandLines.IllustID)

	} else if arguments.CommandLines.AuthorID != 0 {
		src.ShellAuthor(arguments.CommandLines.AuthorID, 0)

	} else if arguments.CommandLines.URL != "" {
		src.DownloaderSingly(utils.GetInt(arguments.CommandLines.URL))

	} else if arguments.CommandLines.Following {
		src.GET_USER_FOLLOWING(arguments.CommandLines.UserID)

	} else if arguments.CommandLines.Recommend {
		src.ShellRecommend("", true)

	} else if arguments.CommandLines.Ranking {
		src.ShellRanking()

	} else if arguments.CommandLines.Stars {
		src.ShellStars(config.Vars.UserID, "")

	} else {
		if len(os.Args) == 1 || os.Args[1] == "-h" || os.Args[1] == "--help" {
			_ = cli.ShowAppHelp(c)
		}

	}
	return nil
}
