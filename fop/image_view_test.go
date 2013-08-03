package fop

import (
	"testing"
)

func TestImageViewRequest(t *testing.T) {

	iv := ImageView {
		Mode: 1,
		Width: 250,
	}

	url := iv.MakeRequest("a")

	if url != "a?imageView/1/w/250" {
		t.Error("result not match")
		return
	}

	iv.Mode = 2
	iv.Height = 250
	iv.Quality = 80
	iv.Format = "jpg"
	url = iv.MakeRequest("a")

	if url != "a?imageView/2/w/250/h/250/q/80/format/jpg" {
		t.Error("result not match")
		return
	}
}
