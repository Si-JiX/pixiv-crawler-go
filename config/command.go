package config

import (
	"gopkg.in/urfave/cli.v1"
)

var CommandLines = struct {
	URL      string
	Name     string
	AuthorID int
}{}

var CommandLineFlag = []cli.Flag{
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
	cli.IntFlag{
		Name:        "a, author",
		Value:       0,
		Usage:       "author id",
		Destination: &CommandLines.AuthorID,
	},
}
