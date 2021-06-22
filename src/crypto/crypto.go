package wcc_crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"main/src/utils"
	"path/filepath"
	"strings"
)

func Encrypt(plaintext []byte, key []byte) []byte {
	// log.Print("File encryption example")

	/*plaintext, err := ioutil.ReadFile("plaintext.txt")
	if err != nil {
		log.Fatal(err)
	}*/

	// The key should be 16 bytes (AES-128), 24 bytes (AES-192) or
	// 32 bytes (AES-256)

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		log.Panic(err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Panic(err)
	}

	// Never use more than 2^32 random nonces with a given key
	// because of the risk of repeat.
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatal(err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext
	// Save back to file
	/*err = ioutil.WriteFile("ciphertext.bin", ciphertext, 0777)
	if err != nil {
		log.Panic(err)
	}*/
}

func Decrypt(ciphertext []byte, key []byte) []byte {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		log.Panic(err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Panic(err)
	}
	nonce := ciphertext[:gcm.NonceSize()]
	ciphertext = ciphertext[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Panic(err)
	}
	return plaintext
}

func EncryptFile(fileName string, newFileName string, clientPublicKey *rsa.PublicKey) (newFleName string, err error) {
	plaintext, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	keyBytes := make([]byte, 32)
	rand.Read(keyBytes)
	/*if len(utils.Settings.SuppliedAESKey) == 0 {
		keyBytes = make([]byte, 32)
		rand.Read(keyBytes)
	} else {
		keyBytes, _ = hex.DecodeString(utils.Settings.SuppliedAESKey)
	}*/

	if len(newFileName) == 0 {
		newFileName = fileName + "." + utils.Settings.EncryptedFileExt
	}
	chipertext := Encrypt(plaintext, keyBytes)
	encryptedAESKey, err := rsa.EncryptOAEP(sha1.New(), rand.Reader, clientPublicKey, keyBytes, nil)
	if err != nil {
		log.Fatal(err)
	}
	chipertextAndKey := append(chipertext, encryptedAESKey...)
	ioutil.WriteFile(newFileName, chipertextAndKey, 0777)
	return newFileName, nil
}

func DecryptFile(encryptedFilename string, clientPrivateKey *rsa.PrivateKey) (filename string, er error) {
	if !strings.HasSuffix(encryptedFilename, "."+utils.Settings.EncryptedFileExt) {
		return "", errors.New("file is not encrypted")
	}
	ct, err := ioutil.ReadFile(encryptedFilename)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	encryptedKey := ct[len(ct)-256:]
	encryptedFile := ct[:len(ct)-256]
	key, err := rsa.DecryptOAEP(sha1.New(), rand.Reader, clientPrivateKey, encryptedKey, nil)
	if err != nil {
		return "", errors.New("failed to decryped AES key")
	}

	plaintext := Decrypt(encryptedFile, key)
	filename = strings.Split(encryptedFilename, "."+utils.Settings.EncryptedFileExt)[0]
	ioutil.WriteFile(filename, plaintext, 0777)
	return filename, nil
}

func EncryptWithRSAPublicKey(plainData []byte, RSAPublicKeyString string) (string, error) {

	// ----- Converting RSA Public key string to Public key object [ start ] -----
	block, _ := pem.Decode([]byte(RSAPublicKeyString))
	if block == nil {
		return "", errors.New("failed to parse PEM block containing the public key")
	}

	RSAPublicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", errors.New("failed to parse DER encoded public key: " + err.Error())
	}
	// ----- Converting RSA Public key string to Public key object [ end ] -----

	// encrypting Data
	randomSource := rand.Reader
	encryptedDataBytes, err := rsa.EncryptOAEP(sha1.New(), randomSource, RSAPublicKey.(*rsa.PublicKey), plainData, nil)
	if err != nil {
		return "", err
	}

	// hex encoding
	encryptedData := hex.EncodeToString(encryptedDataBytes)
	return encryptedData, nil
}

func GetFilesInCurrentDir(fileFormats string, dirPath string, recursive bool) []string {
	filePaths := []string{}
	absDirPath, err := filepath.Abs(dirPath)
	if recursive {
		subdirs := getDirectoriesInPath(dirPath)
		if len(subdirs) != 0 {
			for _, sd := range subdirs {
				filePaths = append(filePaths, GetFilesInCurrentDir(fileFormats, sd, true)...)
			}
		}
	}
	ffs := strings.Split(fileFormats, ",")
	if len(dirPath) == 0 {
		dirPath = "."
	}
	if err != nil {
		log.Fatal("Directory path invalid")
	}
	fmt.Printf("Encryption in %v directory \n", absDirPath)
	allFiles, err := ioutil.ReadDir(dirPath)
	if !strings.HasSuffix(dirPath, "/") {
		dirPath = dirPath + "/"
	}
	if err != nil {
		log.Fatal("failed to read directory", dirPath)
		return []string{}
	}
	for _, file := range allFiles {
		if extInArray(ffs, filepath.Ext(file.Name())) {
			abs, err := filepath.Abs(dirPath + file.Name())
			if err != nil {
				log.Fatal("failed to read full file path")
				return []string{}
			}
			filePaths = append(filePaths, abs)
		}
	}
	return filePaths
}

func extInArray(arr []string, ext string) bool {
	ext = strings.TrimPrefix(ext, ".")
	for _, e := range arr {
		if e == ext {
			return true
		}
	}
	return false
}

func getDirectoriesInPath(path string) []string {
	files, err := ioutil.ReadDir(path)
	dirs := []string{}
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if f.IsDir() {
			fp := filepath.Join(path, f.Name())
			if err != nil {
				log.Fatal(err)
			}
			dirs = append(dirs, fp)
		}
	}
	return dirs
}

func GetHash(data []byte) string {
	sha := sha256.New()
	sha.Write(data)
	keyHash := sha.Sum(nil)
	return hex.EncodeToString(keyHash)
}
