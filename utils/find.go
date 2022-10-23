package utils

import (
	"regexp"
)

func Find(compile_str string, key_words string) []string {
	return regexp.MustCompile(compile_str).FindAllString(key_words, -1)
}

func FindInt(key_words string) string {
	FindInfo := Find(`(\d+)`, key_words)
	if FindInfo != nil {
		return FindInfo[0]
	}
	return ""
}
