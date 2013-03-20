package rs

import (
	"os"
	"testing"
	. "github.com/qiniu/api/conf"
)

var key = "aa"
var bucketName = "a"
var client Client

func init() {
	ACCESS_KEY = os.Getenv("QINIU_ACCESS_KEY")
	SECRET_KEY = os.Getenv("QINIU_SECRET_KEY")
	if ACCESS_KEY == "" || SECRET_KEY == "" {
		panic("require ACCESS_KEY & SECRET_KEY")
	}
	client = New()
}

func TestBatch(t *testing.T) {
	b, err := client.BatchStat(nil, []EntryPath{
		{bucketName, key},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(b) != 1 {
		t.Fatal("BatchStat failed: len(result) =", len(b))
	}
	c, err := client.Stat(nil, bucketName, key)
	if b[0].Data != c {
		t.Error("result not match")
	}
}

