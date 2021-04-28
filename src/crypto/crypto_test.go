package wcc_crypto

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path"
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
	pt := Decrypt(ct, keyBytes)

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
	pt := Decrypt(ct, keyBytes)
	if string(pt) != testString {
		t.Error("Decryption process failed")
	}
}

func TestDecryptFile(t *testing.T) {
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

	fn, dErr := DecryptFile(newFileName, keyBytes)
	if dErr != nil {
		t.Error(dErr)
	}

	if fileName != fn {
		t.Error("File name mismatch")
	}
	if p, _ := ioutil.ReadFile(fn); string(p) != testString {
		t.Error("Decryption failed")
	}
}

func TestDecryptFileFileNotEncrypted(t *testing.T) {
	testDir := t.TempDir()
	fileName := testDir + "/_testFile.txt"
	testString := "This is a test string"
	err := ioutil.WriteFile(fileName, []byte(testString), 0777)
	if err != nil {
		t.Error(err)
	}
	key := "645267556B58703273357638792F423F4528472B4B6250655368566D59713374"
	keyBytes, err := hex.DecodeString(key)
	if err != nil {
		panic("Hex decode failed")
	}
	_, e := DecryptFile(fileName, keyBytes)
	if e == nil {
		t.Error("Error expected")
	}

	if e.Error() != "file is not encrypted" {
		t.Error("Expected file is not encrypted error")
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
	files := GetFilesInCurrentDir("txt,doc", testDir, false)
	if len(files) != 5 {
		t.Error("File count mismatch")
	}

	files = GetFilesInCurrentDir("jpeg", testDir, false)
	if len(files) != c+1 {
		t.Error("File count mismatch")
	}

	files = GetFilesInCurrentDir("jpg", testDir, false)
	if len(files) != 0 {
		t.Error("File count mismatch")
	}
}

func TestGetFilesInCurrentDirRecursive(t *testing.T) {
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
	recursiveDir1 := testDir + "/testRecursion1"
	recursiveDir2 := testDir + "/testRecursion2"
	recursiveDir3 := testDir + "/wow/this/is/sooo/deep"
	recursiveDir4 := testDir + "/wow/this"
	os.MkdirAll(recursiveDir1, 0777)
	os.MkdirAll(recursiveDir2, 0777)
	os.MkdirAll(recursiveDir3, 0777)
	os.MkdirAll(recursiveDir4, 0777)
	fileName = path.Join(recursiveDir1, "recursive_1_")
	for i := 0; i < c; i++ {
		fullFileName := fileName + fmt.Sprint(i) + ".txt"
		ioutil.WriteFile(fullFileName, []byte(testString), 0777)
		if _, err := os.Stat(fullFileName); err != nil {
			t.Error("file not created correctly")
		}
	}
	fileName = path.Join(recursiveDir2, "recursive_2_")
	for i := 0; i < c; i++ {
		fullFileName := fileName + fmt.Sprint(i) + ".txt"
		ioutil.WriteFile(fullFileName, []byte(testString), 0777)
		if _, err := os.Stat(fullFileName); err != nil {
			t.Error("file not created correctly")
		}
	}

	fileName = path.Join(recursiveDir3, "recursive_3_")
	fullFileName := fileName + fmt.Sprint(1) + ".txt"
	ioutil.WriteFile(fullFileName, []byte(testString), 0777)
	if _, err := os.Stat(fullFileName); err != nil {
		t.Error("file not created correctly")
	}

	fileName = path.Join(recursiveDir4, "recursive_4_")
	fullFileName = fileName + fmt.Sprint(1) + ".txt"
	ioutil.WriteFile(fullFileName, []byte(testString), 0777)
	if _, err := os.Stat(fullFileName); err != nil {
		t.Error("file not created correctly")
	}

	files := GetFilesInCurrentDir("txt", testDir, true)
	if len(files) != 3*c+2 {
		t.Error("File count mismatch")
	}

	for _, f := range files {
		fmt.Println(f)
	}
}
