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

type ImageInfo struct {
	
}

func (this ImageInfo) MakeRequest(url string) string {
	return url + "?imageInfo"	
}

func (this ImageInfo) Call(l rpc.Logger, url string) (ret ImageInfoRet, err error) {
	err = rpc.DefaultClient.Call(l, &ret, this.MakeRequest(url))
	return
}
