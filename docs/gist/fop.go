package main
// @gist import
import "github.com/qiniu/api/fop"
// @endgist
import "github.com/qiniu/rpc"

var logger rpc.Logger
var bucketName = "<bucketName>"
var key = "<key>"

func main() {
	// @gist imageurl
	imageUrl := "http://domain/key"
	// @endgist
	
	// @gist image_info
	ii := fop.ImageInfo{}
	inforet := ii.MakeRequest(imageUrl) // fop.ImageInfoRet, error
	// @endgist
	
	// @gist exif
	exif := fop.Exif{}
	exifret, err := exif.Call(logger, imageUrl) // fop.ExifRet, error
	// @endgist
	
	// @gist image_view
	iv := fop.ImageView{
		Mode: 1,
		Width: 200,
		Height: 200,
	}
	previewUrl := iv.MakeRequest(imageUrl)
	// @endgist
	
	_, _, _, _ = exifret, inforet, err, previewUrl
}
