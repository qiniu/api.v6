package io

import (
	"os"
	"testing"
	"math/rand"
	. "github.com/qiniu/api/conf"
	"github.com/qiniu/api/rs"
)

var (
	bucket = "a"
	testkey = "resumableput_key"
	testfile = "resumable_api_test.go"
)

func init() {
	ACCESS_KEY = os.Getenv("QINIU_ACCESS_KEY")
	SECRET_KEY = os.Getenv("QINIU_SECRET_KEY")
	if ACCESS_KEY == "" || SECRET_KEY == "" {
		panic("require ACCESS_KEY & SECRET_KEY")
	}
}


func TestPut(t *testing.T) {
	policy := rs.PutPolicy {
		Scope: bucket,
	}
	var ret PutRet
	f, err := os.Open(testfile)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	
	fi, err := f.Stat()
	if err != nil {
		t.Fatal(err)
	}

	var mockerr bool
	blockNotify := func (blkIdx int, blkSize int, ret *BlkputRet) {
		if rand.Int()%3 == 0 && mockerr == false {
			if ret.Ctx != "" {
				ret.Ctx = ""
				mockerr = true
			}
		}
	}

	extra := &PutExtra {
		Bucket: bucket,
		ChunkSize: 128,
		MimeType: "text/plain",
		Notify: blockNotify,
	}
	defer rs.New().Delete(nil, bucket, testkey)

	err = Put(nil, &ret, policy.Token(), testkey, f, fi.Size(), extra)
	if err != nil {
		t.Fatal(err)
	}
}
