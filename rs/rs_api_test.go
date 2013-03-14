package rs

import (
	"testing"
	
	. "github.com/qiniu/api/conf"
)

var key = "aa"
var bucketName = "a"
var client Client

func init() {
	ACCESS_KEY = "tGf47MBl1LyT9uaNv-NZV4XZe7sKxOIa9RE2Lp8B"
	SECRET_KEY = "zhbiA6gcQMEi22uZ8CBGvmbnD2sR8SO-5S8qlLCG"
	client = New()
}

func TestBatch(t *testing.T) {
	batch := NewBatcher()
	batch.Add(URIStat(bucketName + ":" + key))
	a := new([]BatchStatResult)
	err := client.DoBatch(nil, a, batch)
	if err != nil {
		t.Error(err)
		return
	}
	
	b, err := client.BatchStat(nil, bucketName + ":" + key)
	if err != nil {
		t.Error(err)
		return
	}
	
	c, _ := client.Stat(nil, bucketName + ":" + key)
	
	if (*a)[0].Data != b[0].Data {
		t.Error("result not match")
	}
	
	if c != b[0].Data {
		t.Error("result not match")
	}
}
