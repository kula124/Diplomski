package utils

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cretz/bine/tor"
)

//go:embed tor.zip
var torZipBytes []byte

var tClient *http.Client = http.DefaultClient

func unzipTor(paramBytes []byte) {
	var fileBytes []byte
	if paramBytes != nil {
		fileBytes = paramBytes
	} else {
		fileBytes = torZipBytes
	}
	dest, _ := os.Getwd()
	dest = dest + "/tor"
	_, err := Unzip(fileBytes, dest)
	if err != nil {
		log.Fatal(err)
	}
	os.Chdir(dest)
}

func cleanup() {
	os.RemoveAll("./tor")
}

func SetupTor(paramBytes []byte) (*tor.Tor, *http.Client, error) {
	fmt.Println("Establishing connection to TOR network... may take a minute")
	defer os.Chdir("..")
	defer cleanup()
	if _, err := os.Stat("./tor.exe"); err != nil {
		unzipTor(paramBytes)
	}

	t, err := tor.Start(context.TODO(), nil)
	if err != nil {
		return nil, nil, err
	}
	dialCtx, dialCancel := context.WithTimeout(context.Background(), time.Minute)
	defer dialCancel()
	// Make connection
	dialer, err := t.Dialer(dialCtx, nil)
	if err != nil {
		return nil, nil, err
	}
	c := &http.Client{Transport: &http.Transport{DialContext: dialer.DialContext}}
	tClient = c
	return t, c, nil
}
