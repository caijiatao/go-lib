package string_util

import (
	"regexp"
)

var numPattern = regexp.MustCompile(`^\d+$|^\d+[.]\d+$`)
var splitNumPattern = regexp.MustCompile("[0-9]+")

func IsNumber(s string) bool {
	return numPattern.MatchString(s)
}

func FindFirstNumberStr(s string) string {
	numList := splitNumPattern.FindAllString(s, -1)
	if len(numList) == 0 {
		return ""
	}
	return numList[0]
}
