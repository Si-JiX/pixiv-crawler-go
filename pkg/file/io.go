package file

import (
	"fmt"
	"os"
)

func Open(file_path string, mode string) []byte {
	if mode == "r" {
		return ReadFile(file_path)
	} else if mode == "w" {
		WriteFile(file_path, "", 0777)
	} else if mode == "a" {
		WriteFile(file_path, "", 0666)
	}
	return nil
}

func ReadFile(file_path string) []byte {
	if f, ok := os.ReadFile(file_path); ok == nil {
		return f
	} else {
		fmt.Println("read file error:", ok)
	}
	return nil
}
func WriteFile(file_path string, data string, perm os.FileMode) {
	if ok := os.WriteFile(file_path, []byte(data), perm); ok != nil {
		fmt.Println("write file error:", ok)
	}
}

func Input(start_info string, info string) string {
	fmt.Println(start_info)
	for {
		var input string
		fmt.Println(info)
		_, _ = fmt.Scanln(&input)
		if input != "" {
			return input
		}
	}
}
