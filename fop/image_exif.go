package fop

import (
	"github.com/qiniu/rpc"
)

type ExifValType struct {
	Val  string `json:"val"`
	Type int    `json:"type"`
}

type ImageExifRet map[string] ExifValType

type ImageExif struct {}

func (this ImageExif) MakeRequest(url string) string {
	return url + "?exif"
}

func (this ImageExif) Call(l rpc.Logger, url string) (ret ImageExifRet, err error) {
 	err = rpc.DefaultClient.Call(l, &ret, this.MakeRequest(url))
 	return
}
