package io

import (
	"testing"
	"bytes"
	"os"

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

func TestPut(t *testing.T) {
	key := "test_1"
	buf := bytes.NewBuffer(nil)
	ret := new(interface{})

	buf.WriteString("hello! new Put")
	rs.New().Delete(nil, bucket, key)
	err := Put(nil, ret, 
		policy.Token(), bucket, key, buf, int64(buf.Len()), extra)
	if err != nil {
		t.Error(err)
	}
}

func TestPutFile(t *testing.T) {
	var ret interface{}
	localFile := "./io_api_test.go"
	key := "test_put_file"

	rs.New().Delete(nil, bucket, key)
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
