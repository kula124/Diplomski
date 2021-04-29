package utils

import "strings"

func FindStringIndex(strArr []string, target string) int {
	c := len(strArr)
	for i := 0; i < c; i++ {
		if strings.Compare(target, strArr[i]) == 0 {
			if i+1 == c {
				return 0
			}
			return i
		}
	}
	return -1
}

func RemoveAtIndex(strArr []string, index int) []string {
	return append(strArr[:index], strArr[index+1:]...)
}
