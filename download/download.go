package download

import (
	"fmt"
	"os"
	"pixiv-cil/config"
	"pixiv-cil/utils"
)

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
