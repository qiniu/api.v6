package up

import (
	"net/http"
	"io"
	"io/ioutil"
	"mime/multipart"
	"bytes"
	"os"
	"path"
	"strconv"
	
	. "github.com/qiniu/api/conf"
	"github.com/qiniu/api/rs"
	"github.com/qiniu/api/auth/up"
	"github.com/qiniu/rpc"
)

// -----------------------------------------------------------------------------

type Client struct {
	Conn rpc.Client
	Host string
	UpToken string
}

func New(upToken, host string) Client {
	if host == "" {
		host = UP_HOST
	}
	t := up.NewTransport(upToken, nil)
	client := &http.Client{Transport: t}
	return Client{rpc.Client{client}, host, upToken}
}

func NewEx(upToken, host string, t http.RoundTripper) Client {
	if host == "" {
		host = UP_HOST
	}
	client := &http.Client{Transport: t}
	return Client{rpc.Client{client}, host, upToken}
}

// -----------------------------------------------------------------------------

type UploadAction struct {
	EntryURI string
	MimeType string
	CustomMeta string
	Crc32 int
	Rotate uint
}

type PutRet struct {
	Hash string
}

func (up Client) Put(l rpc.Logger,
	fileName string, ua UploadAction, reader io.Reader) (ret PutRet, err error) {
	
	action := "/rs-put/" + MarshalUploadAction(ua)
	url := up.Host + "/upload"
	
	bodyBuf := bytes.NewBuffer(nil)
	writer := multipart.NewWriter(bodyBuf)
	writer.WriteField("auth", up.UpToken)
	writer.WriteField("action", action)
	writer.CreateFormFile("file", fileName)
	fileContent, err := ioutil.ReadAll(reader)
	if err != nil {
		return
	}
	bodyBuf.Write(fileContent)
	writer.Close()
	
	bodyType := writer.FormDataContentType()
	err = up.Conn.CallWith(l, &ret, url, bodyType, bodyBuf, bodyBuf.Len())
	return
}

func (up Client) PutFile(l rpc.Logger,
	filePath string, ua UploadAction) (ret PutRet, err error) {
	
	fileName := path.Base(filePath)
	f, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer f.Close()
	
	return up.Put(l, fileName, ua, f)
}

// -----------------------------------------------------------------------------

func MarshalUploadAction(ua UploadAction) (uri string) {
	buf := bytes.NewBuffer(make([]byte, 0, bytes.MinRead))
	if ua.EntryURI != "" {
		buf.WriteString(rs.EncodeURI(ua.EntryURI))
		buf.WriteByte('/')
	}
	
	if ua.MimeType != "" {
		buf.WriteString("mimeType/")
		buf.WriteString(rs.EncodeURI(ua.MimeType))
		buf.WriteByte('/')
	}
	
	if ua.CustomMeta != "" {
		buf.WriteString("meta/")
		buf.WriteString(rs.EncodeURI(ua.CustomMeta))
		buf.WriteByte('/')
	}
	
	if ua.Crc32 != 0 {
		buf.WriteString("crc32/")
		buf.Write(strconv.AppendInt([]byte{}, int64(ua.Crc32), 10))
		buf.WriteByte('/')
	}
	
	if ua.Rotate != 0 {
		buf.WriteString("rotate/")
		buf.Write(strconv.AppendInt([]byte{}, int64(ua.Rotate), 10))
		buf.WriteByte('/')
	}
	
	return string(buf.Bytes())
}
