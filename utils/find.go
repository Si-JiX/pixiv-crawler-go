package utils

import (
	"regexp"
	"strings"
)

func FindList(s string, find_string any) bool {
	switch find_string.(type) {
	case string:
		return strings.Contains(s, find_string.(string))
	case []string:
		for _, v := range find_string.([]string) {
			if strings.Contains(s, v) {
				return true
			}
		}
	}
	return false
}
func ListFind(s any, find_string string) bool {
	switch s.(type) {
	case string:
		return strings.Contains(s.(string), find_string)
	case []string:
		for _, v := range s.([]string) {
			if strings.Contains(v, find_string) {
				return true
			}
		}
	}
	return false
}

func GetAll(compile_str string, key_words string) []string {
	return regexp.MustCompile(compile_str).FindAllString(key_words, -1)
}

func GetInt(key_words string) string {
	FindInfo := GetAll(`(\d+)`, key_words)
	if FindInfo != nil {
		return FindInfo[0]
	}
	return ""
}
