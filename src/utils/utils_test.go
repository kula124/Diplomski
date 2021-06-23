package utils

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
)

func TestQueuePushPop(t *testing.T) {
	var q Queue
	q.Init([]string{"Test1", "Test2"})
	t1 := q.Pop()
	t2 := q.Pop()
	t3 := q.Pop()

	if t2 != "Test2" {
		t.Error("expected Test2")
	}

	if t1 != "Test1" {
		t.Error("expected Test1")
	}

	if len(t3) > 0 {
		t.Error("Expected empty string")
	}

	q.Push("hue")
	if q.queue[0] != "hue" {
		t.Error("Expected hue")
	}

	q.Pop()
	if len(q.queue) != 0 {
		t.Error("Expected empty array")
	}
}

func TestUnmarshalJson(t *testing.T) {
	testString := `{"status":"NOT PAID","message":"Ransom for selected key is not registered as paid at this moment"}`
	res := unmarshalMessage(testString)
	if len(res.Key) > 0 {
		t.Error("key unexpected")
	}
	if res.Status != string(UNPAID) {
		t.Error("Status mismatch")
	}
	if res.Message != "Ransom for selected key is not registered as paid at this moment" {
		t.Error("Message mismatch")
	}
}

func TestUnmarshalJsonPaid(t *testing.T) {
	testString := `{"status":"PAID","key":"0000000000000000000000000000000000000000000000000000000000000000"}`
	res := unmarshalMessage(testString)
	if res.Key != "0000000000000000000000000000000000000000000000000000000000000000" {
		t.Error("key mismatch")
	}
	if res.Status != string(PAID) {
		t.Error("Status mismatch")
	}
	if len(res.Message) > 0 {
		t.Error("Message unexpected")
	}
}

func TestUnmarshalJsonFailed(t *testing.T) {
	testString := `{aaaaaaa}`
	res := unmarshalMessage(testString)
	if len(res.Key) > 0 {
		t.Error("Key unexpected")
	}
	if res.Status != string(ERROR) {
		t.Error("Status mismatch")
	}
	if res.Message != "Failed to parse JSON" {
		t.Error("Error expected")
	}
}

func teardown(files []string) {
	for _, f := range files {
		os.Remove(f)
	}
}

func TestTorClient(t *testing.T) {
	t.Skip()
	const TorLinksIndex = "http://zqktlwiuavvvqqt4ybvgvi7tyo4hjl5xgfuvpdf6otjiycgwqbym2qad.onion/wiki/index.php/Main_Page"
	fileBytes, err := ioutil.ReadFile("./tor.zip")
	if err != nil {
		log.Fatal(err)
	}
	dest, _ := os.Getwd()
	dest = dest + "/tor"
	torFiles, err := Unzip(fileBytes, dest)
	if err != nil {
		log.Fatal(err)
	}
	os.Chdir(dest)
	// log.Println(torFiles)
	tor, client, err := SetupTor(nil)
	defer tor.Close()
	if err != nil {
		t.Fatalf("Failed to run Tor: %v", err)
	}
	resp, err := client.Get(TorLinksIndex)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		log.Fatal("Status is not OK")
	}
	if err != nil || resp.StatusCode != 200 {
		log.Fatal(err)
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	rawBody, _ := ioutil.ReadAll(resp.Body)
	println(string(rawBody))
	if !strings.Contains(string(rawBody), "<!DOCTYPE html PUBLIC") {
		log.Fatal("Unexpected response")
	}
	teardown(torFiles)
}
