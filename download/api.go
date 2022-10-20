package download

import (
	"encoding/json"
	"pixiv-cil/structs"
	"pixiv-cil/utils"
)

//	func GET_AUTHOR(author_id string, page int) {
//		result, err := utils.Request(fmt.Sprintf("https://api.obfs.dev/api/pixiv/member_illust?id=%v&page=%v", author_id, page))
//		if err == nil {
//			_ = json.Unmarshal(result, &structs.AuthorStruct)
//			for _, Illust := range structs.AuthorStruct.Illusts {
//				config.ImageList = append(config.ImageList, Illust)
//			}
//			fmt.Println(structs.AuthorStruct.NextURL)
//			if structs.AuthorStruct.NextURL != nil {
//				GET_AUTHOR(author_id, page+1)
//			}
//		} else {
//			GET_AUTHOR(author_id, page)
//		}
//
// }
func GET_IMAGE_INFO(ImageID string) []byte {
	if result, err := utils.Request("https://www.pixiv.net/ajax/illust/" + ImageID); err == nil {
		_ = json.Unmarshal(result, &structs.IllustStruct)
		println("作者名称", structs.IllustStruct.Body.UserName)
		println("作者ID", structs.IllustStruct.Body.UserName)
		ImageDownloader(structs.IllustStruct.Body.Urls.Original, structs.IllustStruct.Body.IllustTitle)
	}
	return nil
}
