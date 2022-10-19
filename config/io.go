package config

import "fmt"

func Input() string {
	for {
		var input string
		fmt.Println("please input image id")
		fmt.Scanln(&input)
		if input != "" {
			return input
		}
	}
}
