package arguments

import (
	"gopkg.in/urfave/cli.v1"
)

var CommandLines = struct {
	IllustID  int
	UserID    int
	Following bool
	Recommend bool
	AuthorID  int
	Name      string
	URL       string
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
	cli.IntFlag{
		Name:        "a, author",
		Value:       0,
		Usage:       "author id",
		Destination: &CommandLines.AuthorID,
	},
}
