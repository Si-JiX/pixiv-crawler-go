package config

import (
	"fmt"
	"strconv"
)

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

// INT string converted to type int.
func INT(s string) int {
	const intSize = 32 << (^uint(0) >> 63)
	sLen := len(s)
	if intSize == 32 && (0 < sLen && sLen < 10) ||
		intSize == 64 && (0 < sLen && sLen < 19) {
		// Fast path for small integers that fit int type.
		s0 := s
		if s[0] == '-' || s[0] == '+' {
			s = s[1:]
			if len(s) < 1 {
				return 0
			}
		}

		n := 0
		for _, ch := range []byte(s) {
			ch -= '0'
			if ch > 9 {
				return 0
			}
			n = n*10 + int(ch)
		}
		if s0[0] == '-' {
			n = -n
		}
		return n
	}

	// Slow path for invalid, big, or underscored integers.
	i64, _ := strconv.ParseInt(s, 10, 0)
	return int(i64)
}
