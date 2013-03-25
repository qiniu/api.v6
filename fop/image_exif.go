package fop

import (
	"github.com/qiniu/rpc"
)

type ValType struct {
	Val  string `json:"val"`
	Type int    `json:"type"`
}

type ImageExifRet struct {
	Model             ValType `json:"Model"`
	ColorSpace        ValType `json:"ColorSpace"`
	ImageLength       ValType `json:"ImageLength"`
	YResolution       ValType `json:"YResolution"`
	ExifVersion       ValType `json:"ExifVersion"`
	ResolutionUnit    ValType `json:"ResolutionUnit"`
	FlashPixVersion   ValType `json:"FlashPixVersion"`
	Software          ValType `json:"Software"`
	Orientation       ValType `json:"Orientation"`
	Make              ValType `json:"Make"`
	DateTimeOriginal  ValType `json:"DateTimeOriginal"`
	UserComment       ValType `json:"UserComment"`
	YCbCrPositioning  ValType `json:"YCbCrPositioning"`
	XResolution       ValType `json:"XResolution"`
	ImageWidth        ValType `json:"ImageWidth"`
	DateTime          ValType `json:"DateTime"`
	DateTimeDigitized ValType `json:"DateTimeDigitized"`
}

type ImageExif struct {
	
}

func (this ImageExif) MakeRequest(url string) string {
	return url + "?exif"
}

func (this ImageExif) Call(l rpc.Logger, url string) (ret ImageExifRet, err error) {
 	err = rpc.DefaultClient.Call(l, &ret, this.MakeRequest(url))
 	return
}
