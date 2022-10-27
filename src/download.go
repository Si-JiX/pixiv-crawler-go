package src

import (
	"fmt"
	"pixiv-cil/pkg/config"
	"pixiv-cil/pkg/file"
	"pixiv-cil/pkg/request"
	"pixiv-cil/pkg/threadpool"
	"pixiv-cil/src/app"
	"pixiv-cil/src/pixivstruct"
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
		if _, err := app.App.Download(illust_id, "imageFile"); err != nil {
			fmt.Println("download fail", IllustID, err)
		}
	}
}
func ThreadDownloadImages(image_list []string) {
	if len(image_list) != 0 {
		fmt.Println("一共", len(image_list), "张图片,开始下载中...")
		threadpool.InitThread()
		threadpool.Threading.ProgressLength = len(image_list)
		for i := 0; i < len(image_list); i++ {
			threadpool.Threading.Add()
			go app.App.ThreadDownloadImage(image_list[i])
		}
		utils.ImageUrlList = nil
		threadpool.Threading.Close() // Wait for all threads to finish
	} else {
		fmt.Println("add image list fail,please check image list")
	}
}

func GET_USER_FOLLOWING(UserID int) {
	if UserID == 0 {
		UserID = config.Vars.UserID
	}
	following, err := app.App.UserFollowing(UserID, "public", 0)
	if err != nil {
		fmt.Println("Request user following fail,please check network", err)
	}
	for index, user := range following.UserPreviews {
		fmt.Println("index:", index, "\tuser_id:", user.User.ID, "\tuser_name:", user.User.Name)
	}
	fmt.Println("一共", len(following.UserPreviews), "个关注的用户")
	for _, user := range following.UserPreviews {
		ThreadDownloadImages(GET_AUTHOR_INFO(user.User.ID, 0))
	}

}

func GET_RECOMMEND(next_url string) {
	recommended, err := app.App.Recommended(next_url, true)
	if err != nil {
		fmt.Println("Request recommend fail,please check network", err)
		return
	}
	for _, illust := range recommended.Illusts {
		if illust.MetaSinglePage.OriginalImageURL == "" {
			for _, img := range illust.MetaPages {
				utils.ImageUrlList = append(utils.ImageUrlList, img.Images.Original)
			}
		} else {
			utils.ImageUrlList = append(utils.ImageUrlList, illust.MetaSinglePage.OriginalImageURL)
		}
	}
	//ThreadDownloadImages(utils.ImageUrlList)
	if recommended.NextURL != "" {
		IllustRecommended := &pixivstruct.IllustRecommended{}
		request.Get(recommended.NextURL, nil).Json(IllustRecommended) // Get the next page
		for _, illust := range IllustRecommended.Illusts {
			println(illust.Title)
		}
		return
	}

}

func GET_AUTHOR_INFO(author_id int, page int) []string {
	illusts, next, err := app.App.UserIllusts(author_id, "illust", page)
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
