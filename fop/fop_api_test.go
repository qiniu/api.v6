package fop

import (
	"testing"
	
	. "github.com/qiniu/api/conf"
	"github.com/qiniu/api/rs"
)

var bucketName = "a"
var key = "cccd"
var policy = rs.GetPolicy {
	Scope: "a",
	Expires: 3600,
}
var client Client

func init() {
	ACCESS_KEY = "tGf47MBl1LyT9uaNv-NZV4XZe7sKxOIa9RE2Lp8B"
	SECRET_KEY = "zhbiA6gcQMEi22uZ8CBGvmbnD2sR8SO-5S8qlLCG"
	client = New()
}

func TestImageInfo(t *testing.T) {
	_, err := ImageInfo(nil, "http://cheneya.qiniudn.com/ffdfd_9")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestImageExif(t *testing.T) {
	_, err := ImageExif(nil, "http://qiniuphotos.qiniudn.com/gogopher.jpg")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestImageMogrifyForPreview(t *testing.T) {
	url := "http://qiniuphotos.qiniudn.com/gogopher.jpg"
	opts := MogrifyOptions {
		Mode: 2,
		Height: 200,
	}
	u := ImageMogrifyForPreview(url, opts)
	if url + "?imageView/2/h/200" != u {
		t.Error(u)
	}
}

func TestGet(t *testing.T) {
	_, err := client.get("a:aa")
	if err != nil {
		t.Error(err)
	}
}
