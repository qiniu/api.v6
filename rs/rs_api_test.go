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

func TestEntry(t *testing.T) {
	einfo, err := client.Stat(nil, bucketName, key)
	if err != nil {
		t.Fatal(err)
	}

	newkey := key + "11111"
	err = client.Copy(nil, bucketName, key, bucketName, newkey)
	if err != nil {
		t.Fatal(err)
	}
	enewinfo, err := client.Stat(nil, bucketName, newkey)
	if err != nil {
		t.Fatal(err)
	}
	if einfo.Hash != enewinfo.Hash {
		t.Fatal("invalid entryinfo:", einfo, enewinfo)
	}

	newkey2 := key + "22222"
	err = client.Move(nil, bucketName, newkey, bucketName, newkey2)
	if err != nil {
		t.Fatal(err)
	}
	enewinfo2, err := client.Stat(nil, bucketName, newkey2)
	if err != nil {
		t.Fatal(err)
	}
	if enewinfo.Hash != enewinfo2.Hash {
		t.Fatal("invalid entryinfo:", enewinfo, enewinfo2)
	}

	err = client.Delete(nil, bucketName, newkey2)
	if err != nil {
		t.Fatal(err)
	}
}
