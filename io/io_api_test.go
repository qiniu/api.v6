package io

import (
	"testing"
	"bytes"
	"os"
	"io"
	"crypto/rand"
	"fmt"

	"github.com/qiniu/api/rs"
	. "github.com/qiniu/api/conf"
)

var bucket = "a"
var policy = rs.PutPolicy {
	Scope: bucket,
}
var extra = &PutExtra {
	MimeType: "text/plain",
	CallbackParams: "hello=yes",
}

func init() {
	ACCESS_KEY = os.Getenv("QINIU_ACCESS_KEY")
	SECRET_KEY = os.Getenv("QINIU_SECRET_KEY")
}

func randomBoundary() string {
	var buf [30]byte
	_, err := io.ReadFull(rand.Reader, buf[:])
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", buf[:])
}

func TestPut(t *testing.T) {
	key := "test_put_" + randomBoundary()
	buf := bytes.NewBuffer(nil)
	ret := new(interface{})

	buf.WriteString("hello! new Put")
	err := Put(nil, ret, 
		policy.Token(), bucket, key, buf, int64(buf.Len()), extra)
	if err != nil {
		t.Error(err)
	}
}

func TestPutFile(t *testing.T) {
	var ret interface{}
	localFile := "./io_api_test.go"
	key := "test_put_file_" + randomBoundary()

	err := PutFile(nil, &ret, policy.Token(), bucket, key, localFile, extra)
	if err != nil {
		t.Error(err)
	}
}

func TestGetUrl(t *testing.T) {
	domain := "http://cheneya.qiniudn.com"
	key := "hello_jpg"
	dntoken := "<dnToken>"
	downloadUrl := GetUrl(domain, key, dntoken)
	if downloadUrl != "http://cheneya.qiniudn.com/hello_jpg?token=<dnToken>" {
		t.Error("result not match")
	}
}
