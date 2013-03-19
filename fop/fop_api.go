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
}

func New() Client {
	t := digest.NewTransport(ACCESS_KEY, SECRET_KEY, nil)
	client := &http.Client{Transport: t}
	return Client{rpc.Client{client}}
}

func NewEx(t http.RoundTripper) Client {
	client := &http.Client{Transport: t}
	return Client{rpc.Client{client}}
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

type MogrifyOptions struct {
	Mode int `uri:""`
	Width int `uri:"w"`
	Height int `uri:"h"`
	Quality int `uri:"q"`
	Format string `uri:"format"`
}

func ImageMogrifyForPreview(
	imageURL string, mogrifyOptions MogrifyOptions) (previewURL string) {
	opts, _ := encodeuri.Marshal(mogrifyOptions)
	return imageURL + "?imageView" + opts
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
