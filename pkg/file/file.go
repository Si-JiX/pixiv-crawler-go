package file

import (
	"fmt"
	"os"
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
