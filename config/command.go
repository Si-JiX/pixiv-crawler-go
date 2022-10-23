package config

import (
	"gopkg.in/urfave/cli.v1"
)

var CommandLines = struct {
	IllustID int
	AuthorID int
	Name     string
	URL      string
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
