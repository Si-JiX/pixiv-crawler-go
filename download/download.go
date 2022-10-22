package download

import (
	"fmt"
	"pixiv-cil/config"
	"pixiv-cil/utils"
	"strconv"
	"sync"
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

var tmpChan = make(chan struct{}, utils.CHANNEL_CACHE)
var waitGroup sync.WaitGroup

func AuthorImageALL(author_id int) {
	GET_AUTHOR_INFO(author_id, 0) // Get all the images of the author and put them in the ImageUrlList
	utils.CurrentImageLength = len(config.ImageUrlList)
	if utils.CurrentImageLength != 0 {
		fmt.Println("一共", utils.CurrentImageLength, "张图片,开始下载中...")
		for _, url := range config.ImageUrlList {
			waitGroup.Add(1)
			go config.App.ThreadDownloadImage(url, tmpChan, &waitGroup)
		}
		waitGroup.Wait()
		config.ImageUrlList = nil
		utils.CurrentImageLength = 0
		utils.CurrentImageIndex = 0
	} else {
		fmt.Println("Request author info fail,please check author_id or network")
	}
}

func GET_AUTHOR_INFO(author_id int, page int) {
	illusts, next, err := config.App.UserIllusts(author_id, "illust", page)
	for _, Illust := range illusts {
		// Test if the image is a manga or not
		if Illust.MetaSinglePage.OriginalImageURL == "" {
			for _, img := range Illust.MetaPages {
				config.ImageUrlList = append(config.ImageUrlList, img.Images.Original)
			}
		} else {
			config.ImageUrlList = append(config.ImageUrlList, Illust.MetaSinglePage.OriginalImageURL)
		}
	}
	if err == nil && next != 0 {
		GET_AUTHOR_INFO(author_id, next)
	}
}
