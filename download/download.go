package download

import (
	"fmt"
	"pixiv-cil/pkg/config"
	"pixiv-cil/pkg/file"
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
	if utils.ListFind(file.ShowFileList("./imageFile"), strconv.Itoa(illust_id)) {
		fmt.Println(IllustID, "is exist, skip")
	} else {
		if _, err := config.App.Download(illust_id, "imageFile"); err != nil {
			fmt.Println("download fail", IllustID, err)
		}
	}
}

func close_thread() {
	utils.ImageUrlList = nil
	utils.WG.Wait() // Wait for all threads to finish
	utils.CurrentImageIndex = 0
}

func AuthorImageALL(author_id int) {
	if image_list := GET_AUTHOR_INFO(author_id, 0); len(image_list) != 0 {
		fmt.Println("一共", len(image_list), "张图片,开始下载中...")
		for i := 0; i < len(image_list); i++ {
			utils.WG.Add(1)
			go config.App.ThreadDownloadImage(image_list[i], len(image_list))
		}
		close_thread()
	} else {
		fmt.Println("Request author info fail,please check author_id or network")
	}
}

func GET_AUTHOR_INFO(author_id int, page int) []string {
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
	return utils.ImageUrlList
}
