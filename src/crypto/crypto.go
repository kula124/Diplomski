package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
)

func Encrypt(plaintext []byte, key []byte) []byte {
	log.Print("File encryption example")

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

func EncryptFile(fileName string, newFileName string, key []byte) (newFleName string, err error) {
	plaintext, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	if len(newFileName) == 0 {
		newFileName = fileName + ".wc"
	}
	chipertext := Encrypt(plaintext, key)
	ioutil.WriteFile(newFileName, chipertext, 0777)
	return newFileName, nil
}

func GetFilesInCurrentDir(fileFormats string, dirPath string) []string {
	ffs := strings.Split(fileFormats, ",")
	var filePaths []string
	if len(dirPath) == 0 {
		dirPath = "."
	}
	allFiles, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Fatal("failed to read directory")
		return []string{}
	}
	for _, file := range allFiles {
		if extInArray(ffs, filepath.Ext(file.Name())) {
			filePaths = append(filePaths, file.Name())
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