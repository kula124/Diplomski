package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	//SETUP
	testString := "This is a test"
	key := "645267556B58703273357638792F423F4528472B4B6250655368566D59713374"
	keyBytes, err := hex.DecodeString(key)
	if err != nil {
		panic("Hex decode failed")
	}
	ct := Encrypt([]byte(testString), keyBytes)
	pt := Decrypt(ct)

	if string(pt) != testString {
		t.Error("Test string mismatch!")
	}
}

func TestFileEncryptionFileName(t *testing.T) {
	// setup
	testDir := t.TempDir()
	fileName := testDir + "/_testFile.txt"
	testString := "This is a test string"
	ioutil.WriteFile(fileName, []byte(testString), 0777)
	if _, err := os.Stat(fileName); err != nil {
		t.Error("file not created correctly")
		//	t.FailNow()
	}
	// TEST
	key := "645267556B58703273357638792F423F4528472B4B6250655368566D59713374"
	keyBytes, err := hex.DecodeString(key)
	if err != nil {
		panic("Hex decode failed")
	}
	newFileName, _ := EncryptFile(fileName, "", keyBytes)
	if newFileName != testDir+"/_testFile.txt.wc" {
		t.Error("newFileName is not correct")
	}

	ct, _ := ioutil.ReadFile(newFileName)
	pt := Decrypt(ct)
	if string(pt) != testString {
		t.Error("Decryption process failed")
	}
}

func TestGetFilesInCurrentDir(t *testing.T) {
	testDir := t.TempDir()
	fileName := testDir + "/_testFile"
	testString := "This is a test string"
	c := 3
	for i := 0; i < c; i++ {
		fullFileName := fileName + fmt.Sprint(i) + ".txt"
		ioutil.WriteFile(fullFileName, []byte(testString), 0777)
		if _, err := os.Stat(fullFileName); err != nil {
			t.Error("file not created correctly")
		}
	}

	for i := 0; i < c+1; i++ {
		fullFileName := fileName + fmt.Sprint(i) + ".jpeg"
		ioutil.WriteFile(fullFileName, []byte(testString), 0777)
		if _, err := os.Stat(fullFileName); err != nil {
			t.Error("file not created correctly")
		}
	}

	for i := 0; i < c-1; i++ {
		fullFileName := fileName + fmt.Sprint(i) + ".doc"
		ioutil.WriteFile(fullFileName, []byte(testString), 0777)
		if _, err := os.Stat(fullFileName); err != nil {
			t.Error("file not created correctly")
		}
	}

	// test
	files := GetFilesInCurrentDir("txt,doc", testDir)
	if len(files) != 5 {
		t.Error("File count mismatch")
	}

	files = GetFilesInCurrentDir("jpeg", testDir)
	if len(files) != c+1 {
		t.Error("File count mismatch")
	}

	files = GetFilesInCurrentDir("jpg", testDir)
	if len(files) != 0 {
		t.Error("File count mismatch")
	}
}
