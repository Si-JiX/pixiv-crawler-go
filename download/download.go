package download

import (
	"fmt"
	"os"
	"pixiv-cil/config"
	"pixiv-cil/utils"
	"strconv"
)

func CurrentDownloader(IllustID interface{}) {
	var err error
	switch IllustID.(type) {
	case string:
		if config.FindImageFile(IllustID.(string)) {
			fmt.Println(IllustID, "is exist, skip")
		} else {
			_, err = config.App.Download(config.INT(IllustID.(string)), "imageFile")
		}
	case int:
		if config.FindImageFile(strconv.Itoa(IllustID.(int))) {
			fmt.Println(IllustID, "is exist, skip")
		} else {
			_, err = config.App.Download(IllustID.(int), "imageFile")
		}
	default:
		fmt.Println("you input is not a number,please check", IllustID)
	}
	if err != nil {
		fmt.Println("download fail", IllustID, err)
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
