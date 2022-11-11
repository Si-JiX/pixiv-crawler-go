package download

import (
	"fmt"
	"github.com/VeronicaAlexia/pixiv-crawler-go/pkg/config"
	"github.com/VeronicaAlexia/pixiv-crawler-go/pkg/progressbar"
	"github.com/VeronicaAlexia/pixiv-crawler-go/pkg/request"
	"github.com/VeronicaAlexia/pixiv-crawler-go/pkg/threadpool"
	"github.com/VeronicaAlexia/pixiv-crawler-go/src/pixiv"
	"github.com/VeronicaAlexia/pixiv-crawler-go/utils/pixivstruct"
	"io"
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
	if !start {
		return illust_struct // return illust struct
	}
	// Thread Download Image List
	if illust_struct.ArrayLength == 0 {
		fmt.Println("No image to download in this page")
	} else {
		fmt.Println("\n\nthe number of images is", illust_struct.ArrayLength, ",start download...")
		illust_struct.Thread.ProgressLength = illust_struct.ArrayLength
		for _, image_url := range illust_struct.DownloadArray {
			illust_struct.Thread.Add()
			go ImagesSingly(image_url, illust_struct) // download image by thread
		}
		illust_struct.Progress.ProgressEnd()
		illust_struct.Thread.Close() // Wait for all threads to finish

	}
	illust_struct.DownloadArray = nil
	return nil
}

func ImagesSingly(url string, thread *Download) {
	if thread != nil {
		defer thread.Thread.Done()
	}
	var fullPath string
	if name := filepath.Base(url); name == "" {
		fullPath = filepath.Join(config.Vars.CacheDir, filepath.Base(url))
	} else {
		fullPath = filepath.Join(config.Vars.CacheDir, name)
	}
	if output, err := os.Create(fullPath); err != nil {
		fmt.Println("create file fail", err)
	} else {
		defer output.Close() // Close the file when the function returns
		Body := request.Get(url, nil, map[string]string{"Referer": pixiv.API_BASE}).GetBody()
		if _, err = io.Copy(output, Body); err != nil {
			fmt.Println("download Copy fail", err)
		}
	}
	if thread != nil {
		thread.Thread.ProgressCountAdd() // progress count add 1
		thread.Progress.AddProgressCount(thread.Thread.GetProgressCount())
	}
}

//func (thread *Download) DownloadImages() {
//	if thread.ArrayLength != 0 {
//		fmt.Println("\n\n一共", thread.ArrayLength, "张图片,开始下载中...")
//		thread.Thread.ProgressLength = thread.ArrayLength
//		for _, image_url := range thread.DownloadArray {
//			thread.Thread.Add()
//			go Images(image_url, thread)
//		}
//		thread.Progress.ProgressEnd()
//		thread.Thread.Close() // Wait for all threads to finish
//	} else {
//		fmt.Println("add image list fail,please check image list")
//	}
//	thread.DownloadArray = nil
//}
