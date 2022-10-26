package input

import (
	"fmt"
	"strconv"
)

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

func InputInt(start_info string, info string) int {
	fmt.Println(start_info)
	for {
		var input string
		fmt.Println(info)
		_, _ = fmt.Scanln(&input)
		if Atoi, err := strconv.Atoi(input); err != nil {
			fmt.Println("please input int:", err)
		} else {
			return Atoi
		}
	}
}
