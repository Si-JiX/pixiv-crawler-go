package download

import (
	"encoding/json"
	"fmt"
	"pixiv-cil/config"
	"pixiv-cil/structs"
	"pixiv-cil/utils"
)

func GET_AUTHOR(author_id uint64, page int) {
	illusts, next, err := config.App.UserIllusts(author_id, "illust", page)
	fmt.Println("next:", next)
	fmt.Println("err:", err)
	for _, Illust := range illusts {
		config.ImageList = append(config.ImageList, Illust)
	}
	if err == nil && next != 0 {
		GET_AUTHOR(author_id, next)
	} else {
		fmt.Println("一共", len(config.ImageList), "张图片")
	}

}
func GET_IMAGE_INFO(ImageID string) []byte {
	if result, err := utils.Request("https://www.pixiv.net/ajax/illust/" + ImageID); err == nil {
		_ = json.Unmarshal(result, &structs.IllustStruct)
		println("作者名称", structs.IllustStruct.Body.UserName)
		println("作者ID", structs.IllustStruct.Body.UserName)
		ImageDownloader(structs.IllustStruct.Body.Urls.Original, structs.IllustStruct.Body.IllustTitle)
	}
	return nil
}
