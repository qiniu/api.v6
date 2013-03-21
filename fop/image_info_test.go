package fop

import (
	"os"
	"testing"
	
	. "github.com/qiniu/api/conf"
)

func init() {
	ACCESS_KEY = os.Getenv("QINIU_ACCESS_KEY")
	SECRET_KEY = os.Getenv("QINIU_SECRET_KEY")
}

func TestImageInfo(t *testing.T) {
	ret, err := ImageInfo(nil, "http://cheneya.qiniudn.com/ffdfd_9")
	if err != nil {
		t.Error(err)
		return
	}
	if ret.Format != "png" || ret.Width != 413 || ret.Height != 232 || ret.ColorModel != "nrgba" {
		t.Error("result not match")
	}
}
