package up

import (
	"testing"
	"bytes"
	
	"github.com/qiniu/api/rs"
	. "github.com/qiniu/api/conf"
)

var bucketName = "a"
var key = "cccd"

func init() {
	ACCESS_KEY = "tGf47MBl1LyT9uaNv-NZV4XZe7sKxOIa9RE2Lp8B"
	SECRET_KEY = "zhbiA6gcQMEi22uZ8CBGvmbnD2sR8SO-5S8qlLCG"
}

func Test(t *testing.T) {
	policy := rs.PutPolicy {
		Scope: bucketName,
		Expires: 3600,
	}
	client := New(policy.Token(), "")
	_ = client
	ua := UploadAction {
		EntryURI: bucketName + ":" + key,
		MimeType: "text/plain",
	}
	_, err := client.Put(nil, "a.txt", ua, bytes.NewBufferString("hahaha"))
	if err != nil {
		t.Error(err)
		return
	}
}
