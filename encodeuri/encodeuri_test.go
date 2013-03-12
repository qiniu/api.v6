package encodeuri

import (
	"testing"
)

type upload struct {
	EntryURI string `uri:",encoded"`
	MimeType string `uri:"mimeType,encoded"`
	Rotate int `uri:"rotate"`
	A int64 `uri:"-"`
	C string `uri:"-"`
}

func TestMarshal(t *testing.T) {
	u := upload {
		EntryURI: "a:vvv",
		MimeType: "image/png",
		Rotate: 2,
		A: 3,
		C: "not allow to display",
	}
	
	uri, err := Marshal(u)
	if err != nil {
		t.Error(err)
		return
	}
	
	target := "/YTp2dnY=/mimeType/aW1hZ2UvcG5n/rotate/2"
	if string(uri) != target {
		t.Error(uri)
		return
	}
}
