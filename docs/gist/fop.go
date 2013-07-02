package gist

import (
	"log"
	"github.com/qiniu/rpc"
	"github.com/qiniu/api/fop"
)

// @gist makeImageInfoUrl
func makeImageInfoUrl(imageUrl string) string {
	ii := fop.ImageInfo{}
	return ii.MakeRequest(imageUrl)
}
// @endgist

func getImageInfo(imageUrl string) {
	// @gist getImageInfo
	var logger rpc.Logger
	var err error
	var ii  = fop.ImageInfo{}
	var infoRet  fop.ImageInfoRet
	infoRet, err = ii.Call(logger, imageUrl)
	if err != nil {
	//产生错误
		log.Println("fop getImageInfo failed:", err)
		return 
	}
	log.Println(infoRet.Height, infoRet.Width, infoRet.ColorModel,
		infoRet.Format)
	// @endgist
}

// @gist makeExifUrl
func makeExifUrl(imageUrl string) string {
	e := fop.Exif{}
	return e.MakeRequest(imageUrl)
}
// @endgist

func getExif(imageUrl string) {
	// @gist getExif
	var logger rpc.Logger
	var err error
	var ie = fop.Exif{}
	var exifRet fop.ExifRet

	exifRet, err = ie.Call(logger, imageUrl)
	if err != nil {
	//产生错误
		log.Println("fop getExif failed:", err)
		return 
	}

	//处理返回结果
	for _, item := range exifRet {
		log.Println(item.Type, item.Val)
	
	}
	// @endgist
}

