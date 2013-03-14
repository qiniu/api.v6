package fop

import (
	"net/http"
	
	. "github.com/qiniu/api/conf"
	"github.com/qiniu/rpc"
	"github.com/qiniu/api/rs"
	"github.com/qiniu/api/encodeuri"
	"github.com/qiniu/api/auth/digest"
)

// ----------------------------------------------------------

type Client struct {
	Conn rpc.Client
	DownloadToken string
}

func New(downloadToken string) Client {
	t := digest.NewTransport(ACCESS_KEY, SECRET_KEY, nil)
	client := &http.Client{Transport: t}
	return Client{rpc.Client{client}, downloadToken}
}

func NewEx(downloadToken string, t http.RoundTripper) Client {
	client := &http.Client{Transport: t}
	return Client{rpc.Client{client}, downloadToken}
}

// ----------------------------------------------------------

type InfoImage struct {
	Format string
	Width int
	Height int
	ColorModel string
}

type ValType struct {
	Val string
	Type int
}

type ExifImage map[string] ValType

func ImageInfo(l rpc.Logger, imageURL string) (info InfoImage, err error) {
	err = rpc.DefaultClient.Call(l, &info, imageURL + "?imageInfo")
	return
}

func ImageExif(l rpc.Logger, imageURL string) (exif ExifImage, err error) {
	err = rpc.DefaultClient.Call(l, &exif, imageURL + "?exif")
	return
}

type ViewOption struct {
	Mode int `uri:""`
	Width int `uri:"w"`
	Height int `uri:"h"`
	Quality int `uri:"q"`
	Format string `uri:"format"`
}

type Gravity string
const (
	NorthWest = "NorthWest"
	North = "North"
	NorthEast = "NorthEast"
	West = "West"
	Center = "Center"
	East = "East"
	SouthWest = "SouthWest"
	South = "South"
	SouthEast = "SouthEast"
)

type MogrifyOption struct {
	AutoOrient bool `uri:"auto-orient"`
	Thumbnail string `uri:"thumbnail"`
	Gravity Gravity `uri:"gravity"`
	Crop string `uri:"crop"`
	Quality uint `uri:"quality"`
	Rotate uint `uri:"rotate"`
	Format string `uri:"format"`
	SaveAs string `uri:"save-as,encoded"`
}

func ImageMogrifyForPreview(
	imageURL string, mogrifyOption MogrifyOption) (previewURL string) {
	opts, _ := encodeuri.Marshal(mogrifyOption)
	return imageURL + "?imageMogr" + opts
}

// -----------------------------------------------------------------------------

type FileProfile struct {
	Expires int
	Url string
	MimeType string
	Hash string
	Fsize int64
}

func (fop *Client) get(entryURI string) (url string, err error) {
	profile := new(FileProfile)
	err = fop.Conn.Call(nil, &profile, RS_HOST + "/get/" + rs.EncodeURI(entryURI))
	if err == nil {
		url = profile.Url
	}
	return
}

func (fop *Client) SaveAs(l rpc.Logger,
	entryURISrc, entryURIDest, opStr string) (ret rs.Entry, err error) {
	
	encodedEntryURIDest := rs.EncodeURI(entryURIDest)
    saveAsString := "/save-as/" + encodedEntryURIDest
    sourceUrl, err := fop.get(entryURISrc)
    if err != nil {
    	return
    }
    newUrl := sourceUrl + "?imageMogr" + opStr + saveAsString
    err = fop.Conn.Call(nil, &ret, newUrl)
    return
}

func (fop *Client) ImageMogrifySaveAs(l rpc.Logger, 
	entryURISrc, entryURIDest string, opts MogrifyOption) (ret rs.Entry, err error) {

	opStr, err := encodeuri.Marshal(opts)
	if err != nil {
		return
	}
	
	return fop.SaveAs(l, entryURISrc, entryURIDest, opStr)
}

