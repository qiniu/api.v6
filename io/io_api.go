package io

import (
	"io"
	"os"
	"mime/multipart"
	"encoding/base64"
	"github.com/qiniu/rpc"
	. "github.com/qiniu/api/conf"
)

type PutExtra struct {
	CallbackParams  string  // 当 uptoken 指定了 CallbackUrl，则 CallbackParams 必须非空
	Bucket          string
	CustomMeta      string  // 可选。用户自定义 Meta，不能超过 256 字节
	MimeType        string  // 可选。在 uptoken 没有指定 DetectMime 时，用户客户端可自己指定 MimeType
}

type PutRet struct {
	Hash string `json:"hash"` // 如果 uptoken 没有指定 ReturnBody，那么返回值是标准的 PutRet 结构
}

func writeMultipart(writer *multipart.Writer, uptoken, key string, data io.Reader, extra *PutExtra) (err error) {

	// auth
	err = writer.WriteField("auth", uptoken)
	if err != nil {
		return
	}

	// action
	action := "/rs-put/" + encodeURI(extra.Bucket + ":" + key)
	if extra.MimeType != "" {
		action += "/mimeType/" + encodeURI(extra.MimeType)
	}
	if extra.CustomMeta != "" {
		action += "/meta/" + encodeURI(extra.CustomMeta)
	}
	err = writer.WriteField("action", action)
	if err != nil {
		return
	}

	// params
	if extra.CallbackParams != "" {
		err = writer.WriteField("params", extra.CallbackParams)
		if err != nil {
			return
		}
	}

	// file
	writerBuf, err := writer.CreateFormFile("file", key)
	if err != nil {
		return
	}
	_, err = io.Copy(writerBuf, data)
	return
}

func Put(l rpc.Logger, ret interface{}, uptoken, key string, data io.Reader, extra *PutExtra) error {

	r, w := io.Pipe()
	defer r.Close()
	writer := multipart.NewWriter(w)

	go func() {
		err := writeMultipart(writer, uptoken, key, data, extra)
		writer.Close()
		w.CloseWithError(err)
	}()

	contentType := writer.FormDataContentType()
	return rpc.DefaultClient.CallWith64(l, ret, UP_HOST + "/upload", contentType, r, 0)
}

func PutFile(l rpc.Logger, ret interface{}, uptoken, key string, localFile string, extra *PutExtra) (err error) {

	f, err := os.Open(localFile)
	if err != nil {
		return
	}
	defer f.Close()

	return Put(l, ret, uptoken, key, f, extra)
}

// ----------------------------------------------------------

func encodeURI(uri string) string {
	return base64.URLEncoding.EncodeToString([]byte(uri))
}

// ----------------------------------------------------------

