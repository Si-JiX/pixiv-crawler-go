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
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Download struct {
	DownloadArray []string
	ArrayLength   int
	Illusts       []pixivstruct.Illust
	Thread        *threadpool.ThreadStruct
	Progress      *progressbar.Bar
}

func DownloadTask(Illusts []pixivstruct.Illust, start bool) *Download {
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
	illust_struct := &Download{
		Illusts:       Illusts,
		Thread:        threadpool.InitThread(),
		DownloadArray: ImageList,
		ArrayLength:   len(ImageList),
		Progress:      progressbar.NewProgress(len(ImageList), ""),
	}
	if start {
		illust_struct.DownloadImages()
		return nil
	} else {
		return illust_struct
	}
}

func Images(url string, thread *Download) {
	if thread != nil {
		defer thread.Thread.Done()
	}
	name := filepath.Base(url)
	if name == "" {
		name = filepath.Base(url)
	}
	fullPath := filepath.Join("imageFile", name)

	if _, err := os.Stat(fullPath); err == nil {
		return
	}

	output, err := os.Create(fullPath)
	if err != nil {
		return
	}
	defer func(output *os.File) {
		err = output.Close()
		if err != nil {
			log.Println(err)
		}
	}(output) // Close the file when the function returns

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req.Header.Add("Referer", pixiv.API_BASE)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("download failed: %s", resp.Status)
	} else {
		if _, err = io.Copy(output, resp.Body); err != nil {
			fmt.Println("download Copy fail", err)
		}

	}
	if thread != nil {
		thread.Thread.ProgressCountAdd() // progress count add 1
		thread.Progress.AddProgressCount(thread.Thread.GetProgressCount())
	}
}

func (thread *Download) DownloadImages() {
	if thread.ArrayLength != 0 {
		fmt.Println("\n\n一共", thread.ArrayLength, "张图片,开始下载中...")
		thread.Thread.ProgressLength = thread.ArrayLength
		for _, image_url := range thread.DownloadArray {
			thread.Thread.Add()
			go Images(image_url, thread)
		}
		thread.Progress.ProgressEnd()
		thread.Thread.Close() // Wait for all threads to finish
	} else {
		fmt.Println("add image list fail,please check image list")
	}
	thread.DownloadArray = nil
}

func DownloaderSingly(illust_id string) error {
	var urls []string
	if utils.ListFind(file.ShowFileList("./imageFile"), illust_id) {
		fmt.Println(illust_id, "is exist, skip")
	} else {
		illust, err := app.App.IllustDetail(illust_id)
		if err != nil {
			fmt.Println("download fail", err)
			return err
		}
		if illust == nil || illust.MetaSinglePage == nil {
			fmt.Println("download fail,illust is nil")
			return nil
		}
		if illust.MetaSinglePage.OriginalImageURL == "" {
			for _, img := range illust.MetaPages {
				urls = append(urls, img.Images.Original)
			}
		} else {
			urls = append(urls, illust.MetaSinglePage.OriginalImageURL)
		}
		for _, url := range urls {
			Images(url, nil)
		}
	}
	return nil
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
		ShellAuthor(user.User.ID, 0)
	}
}

func ShellStars(user_id int, next_url string) {
	bookmarks, err := app.App.UserBookmarksIllust(user_id, next_url)
	if err != nil {
		fmt.Println("Request user bookmarks illust fail,please check network", err)
	} else {
		DownloadTask(bookmarks.Illusts, true)
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
	illusts, err := app.App.IllustRanking(RankingMode[input.OutputInt(">", ">", len(RankingMode))])
	if err != nil {
		fmt.Println("Ranking request fail,please check network", err)
	} else {
		DownloadTask(illusts.Illusts, true)
	}
}

func ShellRecommend(next_url string, auth bool) {
	if recommended, err := app.App.Recommended(next_url, auth); err != nil {
		fmt.Println("Recommended request fail,please check network", err)
	} else {
		DownloadTask(recommended.Illusts, true)
		if recommended.NextURL != "" {
			ShellRecommend(recommended.NextURL, auth)
		}
	}
}

func ShellAuthor(author_id int, page int) {
	illusts, next, err := app.App.UserIllusts(author_id, "illust", page)
	if err == nil {
		DownloadTask(illusts, true)
		if err == nil && next != 0 { // If there is a next page, continue to request
			ShellAuthor(author_id, next)
		}
	} else {
		fmt.Println("Request author info fail,please check network", err)
	}
}
