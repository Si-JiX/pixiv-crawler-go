package main

import (
	"fmt"
	"gopkg.in/urfave/cli.v1"
	"log"
	"os"
	"regexp"
	"sss/config"
	"sss/download"
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

func main() {
	if !config.IsExist("imageFile") {
		_ = os.Mkdir("imageFile", 0777)
	}
	init_command()
}
