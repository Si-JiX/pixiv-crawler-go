package pixiv

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// Download a specific picture from pixiv id
func (a *AppPixivAPI) Download(id int, path string) (sizes []int64, err error) {
	illust, err := a.IllustDetail(id)
	if err != nil {
		err = errors.Wrapf(err, "illust %d detail error", id)
		return
	}
	if illust == nil {
		err = errors.Wrapf(err, "illust %d is nil", id)
		return
	}
	if illust.MetaSinglePage == nil {
		err = errors.Wrapf(err, "illust %d has no single page", id)
		return
	}

	var urls []string
	if illust.MetaSinglePage.OriginalImageURL == "" {
		for _, img := range illust.MetaPages {
			urls = append(urls, img.Images.Original)
		}
	} else {
		urls = append(urls, illust.MetaSinglePage.OriginalImageURL)
	}

	dclient := &http.Client{}
	if a.proxy != nil {
		dclient.Transport = &http.Transport{
			Proxy: http.ProxyURL(a.proxy),
		}
	}
	if a.timeout != 0 {
		dclient.Timeout = a.timeout
	}

	for _, u := range urls {
		size, e := download(dclient, u, path, filepath.Base(u))
		if e != nil {
			err = errors.Wrapf(e, "download url %s failed", u)
			return
		}
		sizes = append(sizes, size)
	}

	return
}

// download image to file (use 6.0 app-api)
func download(client *http.Client, url, path, name string) (int64, error) {
	if path == "" {
		return 0, fmt.Errorf("download path needed")
	}
	if name == "" {
		name = filepath.Base(url)
	}
	fullPath := filepath.Join(path, name)

	if _, err := os.Stat(fullPath); err == nil {
		return 0, nil
	}

	output, err := os.Create(fullPath)
	if err != nil {
		return 0, err
	}
	defer func(output *os.File) {
		err = output.Close()
		if err != nil {
			log.Println(err)
		}
	}(output)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Add("Referer", API_BASE)
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("download failed: %s", resp.Status)
	}

	n, err := io.Copy(output, resp.Body)
	if err != nil {
		return 0, err
	}
	return n, nil
}
