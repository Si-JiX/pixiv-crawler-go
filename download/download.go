package download

import (
	"fmt"
	"os"
	"pixiv-cil/config"
	"pixiv-cil/utils"
)

func CurrentDownloader(IllustID interface{}) {
	var (
		size []int64
		err  error
	)
	switch IllustID.(type) {
	case string:
		size, err = config.App.Download(config.INT(IllustID.(string)), "imageFile")
	case int:
		size, err = config.App.Download(IllustID.(int), "imageFile")
	default:
		fmt.Println("you input is not a number,please check", IllustID)
	}
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("download success,download size:", size)
	}
}

func ImageDownloader(ImageURL, ImageName string) {
	if config.IsExist("./imageFile/" + ImageName + ".jpg") {
		fmt.Printf("file:%v exist\r", ImageName)
	} else {
		if image, _ := utils.Request(ImageURL); image != nil {
			fmt.Println(ImageName, "download success")
			_ = os.WriteFile("imageFile"+"/"+ImageName+".jpg", image, 0666)
		} else {
			fmt.Println("image download fail")
		}
	}
}
