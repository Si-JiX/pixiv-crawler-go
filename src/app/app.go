package app

import (
	"fmt"
	"github.com/VeronicaAlexia/pixiv-crawler-go/pkg/config"
	"github.com/VeronicaAlexia/pixiv-crawler-go/pkg/file"
	"github.com/VeronicaAlexia/pixiv-crawler-go/pkg/input"
	"github.com/VeronicaAlexia/pixiv-crawler-go/src"
	"github.com/VeronicaAlexia/pixiv-crawler-go/src/pixiv"
	"github.com/VeronicaAlexia/pixiv-crawler-go/utils"
)

var App = pixiv.NewApp()

func DownloaderSingly(illust_id string) {
	var urls []string
	if utils.ListFind(file.ShowFileList("./imageFile"), illust_id) {
		fmt.Println(illust_id, "is exist, skip")
	} else {
		illust, err := App.IllustDetail(illust_id)
		if err != nil {
			fmt.Println("download fail", err)
			return
		}
		if illust == nil || illust.MetaSinglePage == nil {
			fmt.Println("download fail,illust is nil")
			return
		}
		if illust.MetaSinglePage.OriginalImageURL == "" {
			for _, img := range illust.MetaPages {
				urls = append(urls, img.Images.Original)
			}
		} else {
			urls = append(urls, illust.MetaSinglePage.OriginalImageURL)
		}
		for index, url := range urls {
			fmt.Println("download", illust.Title, "\timage", index+1, "of", len(urls))
			src.ImagesSingly(url, nil)
		}
		fmt.Println("\033[2J")
	}
}

func GET_USER_FOLLOWING(UserID int) {
	if UserID == 0 {
		UserID = config.Vars.UserID
	}
	following, err := App.UserFollowing(UserID, "public", 0)
	if err != nil {
		fmt.Println("Request user following fail,please check network", err)
	}
	for index, user := range following.UserPreviews {
		fmt.Println("index:", index, "\tuser_id:", user.User.ID, "\tuser_name:", user.User.Name)
	}
	fmt.Println("一共", len(following.UserPreviews), "个关注的用户")
	for _, user := range following.UserPreviews {
		ShellAuthor(user.User.ID, 0)
	}
	// 刷新屏幕
}

func ShellStars(user_id int, next_url string) {
	bookmarks, err := App.UserBookmarksIllust(user_id, next_url)
	if err != nil {
		fmt.Println("Request user bookmarks illust fail,please check network", err)
	} else {
		src.DownloadTask(bookmarks.Illusts, true)
		if bookmarks.NextURL != "" {
			ShellStars(user_id, bookmarks.NextURL)
		}
	}
}

func ShellRanking() {
	RankingMode := []string{"day", "week", "month", "day_male", "day_female", "week_original", "week_rookie", "day_manga"}
	for index, mode := range RankingMode {
		fmt.Println("index:", index, "\tmodel:", mode)
	}
	illusts, err := App.IllustRanking(RankingMode[input.OutputInt(">", ">", len(RankingMode))])
	if err != nil {
		fmt.Println("Ranking request fail,please check network", err)
	} else {
		src.DownloadTask(illusts.Illusts, true)
	}
}

func ShellRecommend(next_url string, auth bool) {
	if recommended, err := App.Recommended(next_url, auth); err != nil {
		fmt.Println("Recommended request fail,please check network", err)
	} else {
		src.DownloadTask(recommended.Illusts, true)
		if recommended.NextURL != "" {
			ShellRecommend(recommended.NextURL, auth)
		}
	}
}

func ShellAuthor(author_id int, page int) {
	illusts, next, err := App.UserIllusts(author_id, "illust", page)
	if err == nil {
		src.DownloadTask(illusts, true)
		if err == nil && next != 0 { // If there is a next page, continue to request
			ShellAuthor(author_id, next)
		}
	} else {
		fmt.Println("Request author info fail,please check network", err)
	}
}
