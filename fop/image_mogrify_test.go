package fop

import (
	"testing"
	"os"
	
	. "github.com/qiniu/api/conf"
)

func init() {
	ACCESS_KEY = os.Getenv("QINIU_ACCESS_KEY")
	SECRET_KEY = os.Getenv("QINIU_SECRET_KEY")
}

func TestMakeRequest(t *testing.T) {
	mogrify := ImageMogrify {
		AutoOrient: true,
		Thumbnail: "!256x256r",
		Gravity: "North",
		Crop: "!256x256",
		Quality: 80,
		Rotate: 1,
		Format: "png",
	}
	uri := mogrify.Marshal()
	if uri != "/auto-orient/thumbnail/!256x256r/gravity/North/crop/!256x256/quality/80/rotate/1/format/png" {
		t.Error("result not match")
	}
	
	mogrify = ImageMogrify {
		Gravity: "South",
		Rotate: 1,
	}
	if mogrify.Marshal() != "/gravity/South/rotate/1" {
		t.Error("result not match")
	}
}

func TestSaveAs(t *testing.T) {
	m := ImageMogrify {
		Format: "jpeg",
	}
	_, err := m.SaveAs(nil, "a:ffdfd_9", "a:ddd")
	if err != nil {
		t.Error(err)
		return
	}
}
