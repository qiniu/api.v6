package io

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"testing"
	"github.com/qiniu/api/rs"
	. "github.com/qiniu/api/conf"
)

var bucket = "a"
var policy = rs.PutPolicy {
	Scope: bucket,
}
var upString = "hello qiniu world"
var extra =  []*PutExtra {
	&PutExtra {
		MimeType: "text/plain",
		CheckCrc: 1,
	},
	&PutExtra {
		MimeType: "text/plain",
		CheckCrc: 0,
	},
	nil,
}

func init() {
	ACCESS_KEY = os.Getenv("QINIU_ACCESS_KEY")
	SECRET_KEY = os.Getenv("QINIU_SECRET_KEY")

	for _, v := range extra {
		if v != nil {
			v.Crc32 = crc32String(upString)
		}
	}
}

func randomBoundary() string {

	var buf [30]byte
	_, err := io.ReadFull(rand.Reader, buf[:])
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", buf[:])
}

func crc32String(s string) uint32 {
	h := crc32.NewIEEE()
	h.Write([]byte(s))
	return h.Sum32()
}

func TestPut(t *testing.T) {

	buf := bytes.NewBuffer(nil)
	ret := new(PutRet)
	for i, v := range extra {
		key := "test_put_" + randomBoundary()
		//test undefined key example
		if i == 1 {
			key = "?"
		}
		buf.WriteString(upString)
		err := Put(nil, ret, policy.Token(nil), key, buf, v)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(ret.Hash, ret.Key)
	}
}

