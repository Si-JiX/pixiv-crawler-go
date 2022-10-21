package download

import (
	"fmt"
	"pixiv-cil/config"
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

var ImageList [][]string

func GET_AUTHOR(author_id int, page int) [][]string {
	illusts, next, err := config.App.UserIllusts(author_id, "illust", page)
	for _, Illust := range illusts {
		ImageList = append(ImageList, []string{strconv.FormatUint(Illust.ID, 10), Illust.Title, Illust.Images.Original})
	}
	if err == nil && next != 0 {
		GET_AUTHOR(author_id, next)
	} else {
		fmt.Println("一共", len(ImageList), "张图片")
	}
	return ImageList
}
