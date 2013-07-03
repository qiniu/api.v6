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
var upFile string
var extra =  []*PutExtra {
	&PutExtra {
		MimeType: "text/plain",
		CheckCrc: 0,
	},
	&PutExtra {
		MimeType: "text/plain",
		CheckCrc: 1,
	},
	&PutExtra {
		MimeType: "text/plain",
		CheckCrc: 2,
	},
	nil,
}

func init() {
	ACCESS_KEY = os.Getenv("QINIU_ACCESS_KEY")
	SECRET_KEY = os.Getenv("QINIU_SECRET_KEY")

	// create a temp file
	upFile = randomBoundary()
	f, _ := os.Create(upFile)
	defer f.Close()
	f.Write([]byte("this is a temp file"))
}

//---------------------------------------

func randomBoundary() string {

	var buf [30]byte
	_, err := io.ReadFull(rand.Reader, buf[:])
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", buf[:])
}

func crc32File(file string) uint32 {
	//it is so simple that do not check any err!!

	f, _ := os.Open(file)
	defer f.Close()
	info, _ := f.Stat()
	defer f.Close()
	h := crc32.NewIEEE()
	buf := make([]byte, info.Size())
	io.ReadFull(f, buf)
	h.Write(buf)
	return h.Sum32()
}

func crc32String(s string) uint32 {

	h := crc32.NewIEEE()
	h.Write([]byte(s))
	return h.Sum32()
}

//---------------------------------------

func TestAll(t *testing.T) {
	testPut(t)
	testPutWithoutKey(t)
	testPutFile(t)
	testPutWithoutKey(t)

	// remove the temp file
	os.Remove(upFile)
}

func testPut(t *testing.T) {

	buf := bytes.NewBuffer(nil)
	ret := new(PutRet)
	for _, v := range extra {
		key := "test_put_" + randomBoundary()
		buf.WriteString(upString)
		if v != nil {
			v.Crc32 = crc32String(upString)
		}

		err := Put(nil, ret, policy.Token(nil), key, buf, v)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(ret.Hash, ret.Key)
	}
}

func testPutWithoutKey(t *testing.T) {

	buf := bytes.NewBuffer(nil)
	ret := new(PutRet)
	for _, v := range extra {
		buf.WriteString(upString)
		if v != nil {
			v.Crc32 = crc32String(upString)
		}

		err := PutWithoutKey(nil, ret, policy.Token(nil), buf, v)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(ret.Hash, ret.Key)
	}
}

func testPutFile(t *testing.T) {

	ret := new(PutRet)
	for _, v := range extra {
		if v != nil {
			v.Crc32 = crc32File(upFile)
		}

		key := "test_put_" + randomBoundary()
		err := PutFile(nil, ret, policy.Token(nil), key, upFile, v)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(ret.Hash, ret.Key)
	}
}

func testPutFileWithoutKey(t *testing.T) {

	ret := new(PutRet)
	for _, v := range extra {
		if v != nil {
			v.Crc32 = crc32File(upFile)
		}

		err := PutFileWithoutKey(nil, ret, policy.Token(nil),  upFile, v)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(ret.Hash, ret.Key)
	}
}
