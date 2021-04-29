package main

import (
	"fmt"
	"io/ioutil"
	"log"
	wcc_crypto "main/src/crypto"
	"main/src/utils"
	"os"
	"path"
	"testing"
)

func createRecursiveFileHierarchy(t *testing.T) ([]string, string) {
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

	return wcc_crypto.GetFilesInCurrentDir("txt", testDir, true), testDir
}

func TestEncryptionDecryptionScheme(t *testing.T) {
	files, rootDir := createRecursiveFileHierarchy(t)
	var dQue utils.Queue
	key := "645267556B58703273357638792F423F4528472B4B6250655368566D59713374"
	dQue.Init(files)
	StartEncryption(&dQue, key)
	for _, f := range files {
		// ff := filepath.Join(rootDir, f)
		err := os.Remove(f)
		if err != nil {
			log.Fatal("Failed to delete the file")
			t.Error(err)
		}
	}
	var eQue utils.Queue
	eFiles := wcc_crypto.GetFilesInCurrentDir("wc", rootDir, true)
	eQue.Init(eFiles)
	StartDecryption(&eQue, key)
	finalFiles := wcc_crypto.GetFilesInCurrentDir("txt", rootDir, true)
	for _, ff := range files {
		if (utils.FindStringIndex(finalFiles, ff)) == -1 {
			t.Error("Final array mismatch")
		}
	}
}
