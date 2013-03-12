package up

import (
	"net/http"
	"io"
	"io/ioutil"
	"mime/multipart"
	"bytes"
	"os"
	"path"
	
	. "github.com/qiniu/api/conf"
	"github.com/qiniu/api/auth/up"
	"github.com/qiniu/rpc"
	"github.com/qiniu/api/encodeuri"
)

// ----------------------------------------------------------

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

// ----------------------------------------------------------

type UploadAction struct {
	EntryURI string `uri:",encoded"`
	MimeType string `uri:"mimeType,encoded"`
	CustomMeta string `uri:"meta,encoded"`
	Crc32 int `uri:"crc32"`
	Rotate uint `uri:"rotate"`
}

type PutRet struct {
	Hash string
}

func (up Client) Put(l rpc.Logger,
	fileName string, ua UploadAction, reader io.Reader) (ret PutRet, err error) {
	
	actionUri, err := encodeuri.Marshal(ua)
	if err != nil {
		return
	}
	action := "/rs-put" + actionUri
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
