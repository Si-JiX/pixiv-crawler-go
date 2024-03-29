package arguments

import (
	"github.com/VeronicaAlexia/pixiv-crawler-go/pkg/cli"
)

var CommandLines = struct {
	IllustID  string
	UserID    int
	Following bool
	Recommend bool
	Ranking   bool
	AuthorID  int
	Stars     bool
	Name      string
	URL       string
}{}

var CommandLineFlag = []cli.Flag{
	cli.StringFlag{
		Name:        "d, download",
		Value:       "",
		Usage:       "input IllustID to download",
		Destination: &CommandLines.IllustID,
	},
	cli.StringFlag{
		Name:        "u, url",
		Value:       "",
		Usage:       "input pixiv url to download",
		Destination: &CommandLines.URL,
	},
	cli.IntFlag{
		Name:        "a, author",
		Value:       0,
		Usage:       "author id",
		Destination: &CommandLines.AuthorID,
	},
	cli.IntFlag{
		Name:        "user, userid",
		Value:       0,
		Usage:       "input user id",
		Destination: &CommandLines.UserID,
	},
	cli.StringFlag{
		Name:        "n, name",
		Value:       "",
		Usage:       "author name",
		Destination: &CommandLines.Name,
	},
	cli.BoolFlag{
		Name:        "f, following",
		Usage:       "following",
		Destination: &CommandLines.Following,
	},
	cli.BoolFlag{
		Name:        "r, recommend",
		Usage:       "recommend illust",
		Destination: &CommandLines.Recommend,
	},
	cli.BoolFlag{
		Name:        "s, stars",
		Usage:       "download stars",
		Destination: &CommandLines.Stars,
	},
	cli.BoolFlag{
		Name:        "rk, ranking",
		Usage:       "ranking illust",
		Destination: &CommandLines.Ranking,
	},
}
