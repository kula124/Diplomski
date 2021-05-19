package utils

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func FindStringIndex(strArr []string, target string) int {
	c := len(strArr)
	for i := 0; i < c; i++ {
		if strings.Compare(target, strArr[i]) == 0 {
			return i
		}
	}
	return -1
}

func RemoveAtIndex(strArrPtr **[]string, index int) []string {
	strArr := **strArrPtr
	newArr := append(strArr[:index], strArr[index+1:]...)
	*strArrPtr = &newArr
	return newArr
}

func GetDesktopLocation() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(home, "Desktop")
}
