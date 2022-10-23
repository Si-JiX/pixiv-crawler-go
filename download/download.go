package download

import (
	"fmt"
	"pixiv-cil/config"
	"pixiv-cil/utils"
	"strconv"
)

func CurrentDownloader(IllustID interface{}) {
	var illust_id int
	switch IllustID.(type) {
	case string:
		illust_id = utils.INT(IllustID.(string))
	case int:
		illust_id = IllustID.(int)
	default:
		fmt.Println("you input is not a number,please check", IllustID)
	}
	if utils.ListFind(config.ShowFileList("./imageFile"), strconv.Itoa(illust_id)) {
		fmt.Println(IllustID, "is exist, skip")
	} else {
		if _, err := config.App.Download(illust_id, "imageFile"); err != nil {
			fmt.Println("download fail", IllustID, err)
		}
	}
}

func AuthorImageALL(author_id int) {
	GET_AUTHOR_INFO(author_id, 0) // Get all the images of the author and put them in the ImageUrlList
	utils.CurrentImageLength = len(utils.ImageUrlList)
	if utils.CurrentImageLength != 0 {
		fmt.Println("一共", utils.CurrentImageLength, "张图片,开始下载中...")
		for _, url := range utils.ImageUrlList {
			utils.WG.Add(1)
			go config.App.ThreadDownloadImage(url)
		}
		utils.ImageUrlList = nil
		utils.WG.Wait() // Wait for all threads to finish
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
				utils.ImageUrlList = append(utils.ImageUrlList, img.Images.Original)
			}
		} else {
			utils.ImageUrlList = append(utils.ImageUrlList, Illust.MetaSinglePage.OriginalImageURL)
		}
	}
	// If there is a next page, continue to request
	if err == nil && next != 0 {
		GET_AUTHOR_INFO(author_id, next)
	}
}
