package io

import (
	"io"
	"os"
	"encoding/base64"
	"mime/multipart"
	
	. "github.com/qiniu/api/conf"
	"github.com/qiniu/rpc"
)

type PutExtra struct {
    CallbackParams  string  // 当 uptoken 指定了 CallbackUrl，则 CallbackParams 必须非空
    Bucket          string
    CustomMeta      string  // 可选。用户自定义 Meta，不能超过 256 字节
    MimeType        string  // 可选。在 uptoken 没有指定 DetectMime 时，用户客户端可自己指定 MimeType
}

type PutRet struct {
	// 如果 uptoken 没有指定 ReturnBody，那么返回值是标准的 PutRet 结构 
    Hash string `json:"hash"`
}

func Put(l rpc.Logger, ret interface{}, uptoken, key string,
	data io.Reader, extra *PutExtra) (err error) {

	url := UP_HOST + "/upload"
	r, w := io.Pipe()
	writer := multipart.NewWriter(w)
	
	go func() {
		defer w.Close()
		defer writer.Close()
		
		// auth
		writer.WriteField("auth", uptoken)
		
		// action
		_, err := writer.CreateFormField("action")
		if err != nil {
			return
		}
		io.WriteString(w, "/rs-put/")
		io.WriteString(w, encodeURI(extra.Bucket + ":" + key))
	
		if extra.MimeType != "" {
			io.WriteString(w, "/mimeType/")
			io.WriteString(w, encodeURI(extra.MimeType))
		}
	
		if extra.CustomMeta != "" {
			io.WriteString(w, "/meta/")
			io.WriteString(w, encodeURI(extra.CustomMeta))
		}
	
		// params
		if extra.CallbackParams != "" {
			writer.WriteField("params", extra.CallbackParams)
		}
	
		// file
		_, err = writer.CreateFormFile("file", key)
		if err != nil {
			return
		}
		io.Copy(w, data)
	}()
	
	contentType := writer.FormDataContentType()
	return rpc.DefaultClient.CallWith64(l, ret, url, contentType, r, 0)
}

func PutFile(l rpc.Logger, ret interface{},
	uptoken, key string, localFile string, extra *PutExtra) (err error) {

	f, err := os.Open(localFile)
	if err != nil {
		return
	}
	defer f.Close()
	
	return Put(l, ret, uptoken, key, f, extra)
}

// ----------------------------------------------------------

func GetUrl(domain string, key string, dntoken string) (downloadUrl string) {
	url := domain + "/" + key
	if dntoken == "" {
		return url
	}
	return url + "?token=" + dntoken
}

// ----------------------------------------------------------

func encodeURI(uri string) string {
	return base64.URLEncoding.EncodeToString([]byte(uri))
}
