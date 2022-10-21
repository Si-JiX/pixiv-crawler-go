package config

import (
	"fmt"
	"os"
	"strings"
)

// ShowFileList  show file list
func ShowFileList(path string) []string {
	var FileNameArray []string
	files, _ := os.ReadDir(path)
	for _, file := range files {
		FileNameArray = append(FileNameArray, file.Name())
	}
	if len(FileNameArray) == 0 {
		return nil
	} else {
		return FileNameArray
	}
}

// FindImageFile Find local image file
func FindImageFile(name string) bool {
	for _, file := range ShowFileList("./imageFile") {
		if strings.Contains(file, name) {
			return true
		}
	}
	return false
}

// PrintMap 美化打印map数据
func PrintMap(m map[string]interface{}) {
	for k, v := range m {
		fmt.Println(k, ":", v)
	}
}

// IsExist is exist file
func IsExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// NewFile IsNotExist is not exist file
func NewFile(filePath string) {
	if !IsExist(filePath) {
		_ = os.Mkdir(filePath, 0777)
	}
}
