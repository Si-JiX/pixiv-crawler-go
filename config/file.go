package config

import (
	"fmt"
	"os"
)

// 美化打印map数据
func PrintMap(m map[string]interface{}) {
	for k, v := range m {
		fmt.Println(k, ":", v)
	}
}

// is exist file
func IsExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
