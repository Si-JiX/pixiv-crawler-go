package input

import (
	"fmt"
	"strconv"
)

func Input(info1 string, info string) string {
	fmt.Printf(info1)
	var input string
	_, _ = fmt.Scanln(&input)
	if input != "" {
		return input
	} else if input == "quit" || input == "exit" {
		return ""
	} else {
		return Input(info, info)
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
