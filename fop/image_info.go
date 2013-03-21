package fop

import (
	"github.com/qiniu/rpc"
)

type ImageInfoRet struct {
	Format string
	Width uint
	Height uint
	ColorModel string
}

func ImageInfo(l rpc.Logger, url string) (ret ImageInfoRet, err error) {
	err = New().Conn.Call(l, &ret, url + "?imageInfo")
	return
}
