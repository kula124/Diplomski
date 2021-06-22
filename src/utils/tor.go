package utils

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cretz/bine/tor"
)

func SetupTor() (*http.Client, error) {
	fmt.Println("Establishing connection to TOR network... may take a minute")
	t, err := tor.Start(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	//defer t.Close()
	// Wait at most a minute to start network and get
	dialCtx, dialCancel := context.WithTimeout(context.Background(), time.Minute)
	defer dialCancel()
	// Make connection
	dialer, err := t.Dialer(dialCtx, nil)
	if err != nil {
		return nil, err
	}
	c := &http.Client{Transport: &http.Transport{DialContext: dialer.DialContext}}
	return c, nil
}
