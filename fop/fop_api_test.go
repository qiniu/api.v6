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

func TestMarshalMogrifyOption(t *testing.T) {
	
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
	opts := MogrifyOption {
		AutoOrient: true,
		Gravity: "North",
		Thumbnail: "haha",
	}
	u := ImageMogrifyForPreview(url, opts)
	if url + "?imageMogr/auto-orient/thumbnail/haha/gravity/North" != u {
		t.Error(u)
	}
}

func TestGet(t *testing.T) {
	_, err := client.get("a:aa")
	if err != nil {
		t.Error(err)
	}
}

func TestSaveAs(t *testing.T) {
	ret, err := client.SaveAs(nil, "a:ffdfd_9", "a:cdca1", "/rotate/25")
	if err != nil {
		t.Error(err)
		return
	}
	
	if ret.Hash == "" {
		t.Error("miss hash")
		return
	}
}

func TestImageMogrifySaveAs(t *testing.T) {
	src := "a:ffdfd_9"
	dest := "a:cdca2"
	opts := MogrifyOption {
		Rotate: 90,
	}
	_, err := client.ImageMogrifySaveAs(nil, src, dest, opts)
	if err != nil {
		t.Error(err)
		return
	}
	
}
