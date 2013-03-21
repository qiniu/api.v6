package fop

import (
	"github.com/qiniu/rpc"
)

type ValType struct {
	Val string
	Type int
}

type ImageExifRet map[string] ValType

func ImageExif(l rpc.Logger, url string) (ret ImageExifRet, err error) {
 	err = New().Conn.Call(l, &ret, url + "?exif")
	return
}
