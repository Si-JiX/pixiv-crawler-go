package pixiv

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

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
