package src

import (
	"fmt"
	"github.com/VeronicaAlexia/pixiv-crawler-go/pkg/config"
	"github.com/VeronicaAlexia/pixiv-crawler-go/pkg/file"
	"github.com/VeronicaAlexia/pixiv-crawler-go/pkg/input"
	"github.com/VeronicaAlexia/pixiv-crawler-go/pkg/progressbar"
	"github.com/VeronicaAlexia/pixiv-crawler-go/pkg/threadpool"
	"github.com/VeronicaAlexia/pixiv-crawler-go/src/app"
	"github.com/VeronicaAlexia/pixiv-crawler-go/src/pixiv"
	"github.com/VeronicaAlexia/pixiv-crawler-go/utils"
	"github.com/VeronicaAlexia/pixiv-crawler-go/utils/pixivstruct"
	"github.com/pkg/errors"
	"net/http"
	"path/filepath"
)

type Download struct {
	DownloadArray []string
	ArrayLength   int
	Illusts       []pixivstruct.Illust
	Thread        *threadpool.ThreadStruct
	Progress      *progressbar.Bar
}

func Downloader(Illusts []pixivstruct.Illust) *Download {
	var ImageList []string
	for _, illust := range Illusts {
		if illust.MetaSinglePage.OriginalImageURL == "" {
			for _, img := range illust.MetaPages {
				ImageList = append(ImageList, img.Images.Original)
			}
		} else {
			ImageList = append(ImageList, illust.MetaSinglePage.OriginalImageURL)
		}
	}
	return &Download{
		Illusts:       Illusts,
		Thread:        threadpool.InitThread(),
		DownloadArray: ImageList,
		ArrayLength:   len(ImageList),
		Progress:      progressbar.NewProgress(len(ImageList), ""),
	}
}

func (thread *Download) Images(url string) {
	defer thread.Thread.Done()
	_, e := pixiv.DownloadMain(&http.Client{}, url, "imageFile", filepath.Base(url))
	if e != nil {
		fmt.Println(errors.Wrapf(e, "download url %s failed", url))
	}
	thread.Thread.ProgressCountAdd() // progress count add 1
	thread.Progress.AddProgressCount(thread.Thread.GetProgressCount())
}
func (thread *Download) DownloadImages() {
	if thread.ArrayLength != 0 {
		fmt.Println("\n\n一共", thread.ArrayLength, "张图片,开始下载中...")
		thread.Thread.ProgressLength = thread.ArrayLength
		for _, image_url := range thread.DownloadArray {
			thread.Thread.Add()
			go thread.Images(image_url)
		}
		thread.Progress.ProgressEnd()
		thread.Thread.Close() // Wait for all threads to finish
	} else {
		fmt.Println("add image list fail,please check image list")
	}
	thread.DownloadArray = nil
}

func CurrentDownloader(illust_id string) {
	if utils.ListFind(file.ShowFileList("./imageFile"), illust_id) {
		fmt.Println(illust_id, "is exist, skip")
	} else {
		if illust := app.App.Download(illust_id, "imageFile"); illust != nil {
			fmt.Printf("image name: %s \t  image id: %d", illust.Title, illust.ID)
		}
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
		GET_AUTHOR_INFO(user.User.ID, 0)
	}
}

func ShellRanking() {
	RankingMode := []string{"day", "week", "month", "day_male", "day_female", "week_original", "week_rookie", "day_manga"}
	for index, mode := range RankingMode {
		fmt.Println("index:", index, "\tmodel:", mode)
	}
	illusts, err := app.App.IllustRanking(RankingMode[input.OutputInt(">", ">", len(RankingMode))])
	if err != nil {
		fmt.Println("Ranking request fail,please check network", err)
	} else {
		download_illusts := Downloader(illusts.Illusts)
		download_illusts.DownloadImages()
	}
}

func ShellRecommend(next_url string, auth bool) {
	if recommended, err := app.App.Recommended(next_url, auth); err != nil {
		fmt.Println("Recommended request fail,please check network", err)
	} else {
		download_illusts := Downloader(recommended.Illusts)
		download_illusts.DownloadImages()
		if recommended.NextURL != "" {
			ShellRecommend(recommended.NextURL, auth)
		}
	}
}

func GET_AUTHOR_INFO(author_id int, page int) {
	illusts, next, err := app.App.UserIllusts(author_id, "illust", page)
	if err == nil {
		download_illusts := Downloader(illusts)
		download_illusts.DownloadImages()
		if err == nil && next != 0 {
			// If there is a next page, continue to request
			GET_AUTHOR_INFO(author_id, next)
		}
	} else {
		fmt.Println("Request author info fail,please check network", err)
	}
}
