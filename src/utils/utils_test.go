package utils

import "testing"

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
