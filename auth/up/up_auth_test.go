package up

import (
	"testing"
	"net/http"
)

func TestTransport(t *testing.T) {
	token := "uptoken"
	client := NewClient(token, nil)
	req, _ := http.NewRequest("GET", "http://qiniutek.com", nil)
	client.Do(req)
	if req.Header.Get("Authorization") != "UpToken " + token {
		t.Error("make Authorization failure")
	}
}
