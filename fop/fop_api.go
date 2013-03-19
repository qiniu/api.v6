package fop

import (
	"net/http"
	"bytes"
	"strconv"
	
	. "github.com/qiniu/api/conf"
	"github.com/qiniu/rpc"
	"github.com/qiniu/api/rs"
	"github.com/qiniu/api/auth/digest"
)

// -----------------------------------------------------------------------------

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

// -----------------------------------------------------------------------------

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

type MogrifyOption struct {
	AutoOrient bool
	Thumbnail string
	Gravity string
	Crop string
	Quality uint
	Rotate uint
	Format string
}

func MarshalMogrifyOption(mo MogrifyOption) (uri string) {
	buf := bytes.NewBuffer(make([]byte, 0, bytes.MinRead))
	if mo.AutoOrient {
		buf.WriteString("/auto-orient")
	}
	
	if mo.Thumbnail != "" {
		buf.WriteString("/thumbnail/")
		buf.WriteString(mo.Thumbnail)
	}
	
	if mo.Gravity != "" {
		buf.WriteString("/gravity/")
		buf.WriteString(mo.Gravity)
	}
	
	if mo.Crop != "" {
		buf.WriteString("/crop/")
		buf.WriteString(mo.Crop)
	}
	
	if mo.Rotate != 0 {
		buf.WriteString("/rotate/")
		buf.Write(strconv.AppendInt([]byte{}, int64(mo.Rotate), 10))
	}
	
	if mo.Quality > 0 {
		buf.WriteString("/quality/")
		buf.Write(strconv.AppendInt([]byte{}, int64(mo.Quality), 10))
		buf.WriteByte('/')
	}
	
	if mo.Format != "" {
		buf.WriteString("/format/")
		buf.WriteString(mo.Format)
	}
	return string(buf.Bytes())
}

func ImageMogrifyForPreview(
	imageURL string, mogrifyOption MogrifyOption) (previewURL string) {
	return imageURL + "?imageMogr" + MarshalMogrifyOption(mogrifyOption)
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
	return fop.SaveAs(l, entryURISrc, entryURIDest, MarshalMogrifyOption(opts))
}

