package fop

import (
	"github.com/qiniu/rpc"
)

type ValType struct {
	Val string
	Type int
}

type ImageExifRet map[string] ValType

type ImageExif struct {
	
}

func (this ImageExif) Call(l rpc.Logger, url string) (ret ImageExifRet, err error) {
 	err = rpc.DefaultClient.Call(l, &ret, url + "?exif")
 	return
}
