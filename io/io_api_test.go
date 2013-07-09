package io

import (
	"bytes"
	"hash/crc32"
	"io"
	"os"
	"strconv"
	"testing"
	"github.com/qiniu/api/rs"
	. "github.com/qiniu/api/conf"
)

var (
	bucket = "a"
	upString = "hello qiniu world"
	policy = rs.PutPolicy {
		Scope: bucket,
	}
	keys []string
	extra =  []*PutExtra {
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
)

func init() {

	ACCESS_KEY = os.Getenv("QINIU_ACCESS_KEY")
	SECRET_KEY = os.Getenv("QINIU_SECRET_KEY")

	for i := 0; i < 3; i++ {
		keys = append(keys, "test_key_" + strconv.Itoa(i))
	}

	// create a temp file
	f, _ := os.Create(keys[0])
	defer f.Close()
	f.Write([]byte("this is a temp file"))
}

//---------------------------------------

func crc32File(file string) uint32 {

	//it is so simple that do not check any err!!
	f, _ := os.Open(file)
	defer f.Close()
	info, _ := f.Stat()
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

	testPut(t, keys[1])
	k1 := testPutWithoutKey(t)
	testPutFile(t, keys[0], keys[2])
	k2 := testPutFileWithoutKey(t, keys[0])

	//clear all keys
	keys = append(keys, k1)
	keys = append(keys, k2)
	for i, k := range keys {
		if i == 0 {
			os.Remove(k)
		} else {
			rs.New(nil).Delete(nil, bucket, k)
		}
	}
}

func testPut(t *testing.T, key string) {

	buf := bytes.NewBuffer(nil)
	ret := new(PutRet)
	for _, v := range extra {
		buf.WriteString(upString)
		if v != nil {
			v.Crc32 = crc32String(upString)
		}

		err := Put(nil, ret, policy.Token(nil), key, buf, v)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func testPutWithoutKey(t *testing.T) string{

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
	}
	return ret.Key
}

func testPutFile(t *testing.T, localFile, key string) {

	ret := new(PutRet)
	for _, v := range extra {
		if v != nil {
			v.Crc32 = crc32File(localFile)
		}

		err := PutFile(nil, ret, policy.Token(nil), key, localFile, v)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func testPutFileWithoutKey(t *testing.T, localFile string) string {

	ret := new(PutRet)
	for _, v := range extra {
		if v != nil {
			v.Crc32 = crc32File(localFile)
		}

		err := PutFileWithoutKey(nil, ret, policy.Token(nil),  localFile, v)
		if err != nil {
			t.Fatal(err)
		}
	}
	return ret.Key
}
