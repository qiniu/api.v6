package io

import (
	"io"
	"os"
	"mime/multipart"
	"fmt"
	"hash/crc32"
	"net/textproto"
	"strconv"
	"strings"

	"github.com/qiniu/rpc"
	. "github.com/qiniu/api/conf"
)

// ----------------------------------------------------------

const UNDEFINED_KEY = "?"

// ----------------------------------------------------------
type PutExtra struct {
	Params   map[string]string    //可选，用户自定义参数，必须以 "x:" 开头
	                              //若不以x:开头，则忽略
	MimeType string               //可选，当为 "" 时候，服务端自动判断 
	Crc32    uint32
	CheckCrc uint32
	        // CheckCrc == 0: 表示不进行 crc32 校验
	        // CheckCrc == 1: 对于 Put 等同于 CheckCrc = 2；对于 PutFile 会自动计算 crc32 值
	        // CheckCrc == 2: 表示进行 crc32 校验，且 crc32 值就是上面的 Crc32 变量
}

type PutRet struct {
	Hash string `json:"hash"`  // 如果 uptoken 没有指定 ReturnBody，那么返回值是标准的 PutRet 结构
	Key  string `json:"key"`   //如果传入的 key == UNDEFINED_KEY，则服务端返回 key
}

// ----------------------------------------------------------

func Put(l rpc.Logger, ret interface{}, uptoken, key string, data io.Reader, extra *PutExtra) error {

	// CheckCrc == 1: 对于 Put 等同于 CheckCrc == 2
	var extra1 PutExtra
	if extra != nil {
		extra1 = *extra
		if extra1.CheckCrc == 1 {
			extra1.CheckCrc = 2
		}
	}
	return put(l, ret, uptoken, key, data, &extra1)
}

func put(l rpc.Logger, ret interface{}, uptoken, key string, data io.Reader, extra *PutExtra) error {
	r, w := io.Pipe()
	defer r.Close()
	writer := multipart.NewWriter(w)

	go func() {
		err := writeMultipart(writer, uptoken, key, data, extra)
		writer.Close()
		w.CloseWithError(err)
	}()

	contentType := writer.FormDataContentType()
	return rpc.DefaultClient.CallWith64(l, ret, UP_HOST, contentType, r, 0)
}

func PutFile(l rpc.Logger, ret interface{}, uptoken, key string, localFile string, extra *PutExtra) (err error) {

	f, err := os.Open(localFile)
	if err != nil {
		return
	}
	defer f.Close()

	return put(l, ret, uptoken, key, f, extra)
}



/*
 * extra.CheckCrc:
 *      0:     不进行crc32校验
 *      1:     以writeMultipart自动生成crc32的值，进行校验
 *      2:     以extra.Crc32的值，进行校验
 *      other: 和2一样， 以 extra.Crc32的值，进行校验   
 */
func writeMultipart(writer *multipart.Writer, uptoken, key string, data io.Reader, extra *PutExtra) (err error) {

	if extra == nil {
		extra = &PutExtra{}
	}

	//token
	if err = writer.WriteField("token", uptoken); err != nil {
		return
	}

	//key
	if key != UNDEFINED_KEY {
		if err = writer.WriteField("key", key); err != nil {
			return
		}
	}

	// extra.Params
	if extra.Params != nil {
		for k, v := range extra.Params {
			if strings.HasPrefix(k, "x:") {
				err = writer.WriteField(k, v)
				if err != nil {
					return 
				}
			}
		}
	}

	//extra.CheckCrc 
	h := crc32.NewIEEE()
	data1 := data
	if extra.CheckCrc == 1 {
		data1 = io.TeeReader(data, h)
	}

	//file 
	head := make(textproto.MIMEHeader)

	//default the filename is same as key , but  ""
	var fileName string
	if key == "" {
		fileName = "index.html"
	} else {
		fileName = key
	}
	head.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,"file" , escapeQuotes(fileName)))
	if  extra.MimeType != "" {
		head.Set("Content-Type", extra.MimeType)
	}

	writerBuf, err := writer.CreatePart(head)
	if err != nil {
		return
	}
	_, err = io.Copy(writerBuf, data1)

	//extra.CheckCrc 
	if extra.CheckCrc == 1 {
		extra.Crc32 = h.Sum32()
	}
	if extra.CheckCrc != 0 {
		err = writer.WriteField("crc32", strconv.FormatInt(int64(extra.Crc32), 10))
	}
	return 
}

// ----------------------------------------------------------

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

