package gist

import (
	"github.com/qiniu/rpc"
	"github.com/qiniu/api/fop"
)

// @gist makeImageInfoUrl
func makeImageInfoUrl(imageUrl string) string {
	ii := fop.ImageInfo{}
	return ii.MakeRequest(imageUrl)
}
// @endgist

// @gist getImageInfo
func getImageInfo(l rpc.Logger, imageUrl string) (ret fop.ImageInfoRet, err error) {
	ii := fop.ImageInfo{}
	return ii.Call(l, imageUrl)
}
// @endgist

// @gist makeExifUrl
func makeExifUrl(imageUrl string) string {
	e := fop.Exif{}
	return e.MakeRequest(imageUrl)
}
// @endgist

// @gist getExif
func getExif(l rpc.Logger, imageUrl string) (ret fop.ExifRet, err error) {
	e := fop.Exif{}
	return e.Call(l, imageUrl)
}
// @endgist

