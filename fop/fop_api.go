package fop

import (
	"net/http"
	"strconv"
	"encoding/base64"
	
	. "github.com/qiniu/api/conf"
	"github.com/qiniu/rpc"
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

func encodeURI(uri string) string {
	return base64.URLEncoding.EncodeToString([]byte(uri))
}

func itoa(a uint) []byte {
	return strconv.AppendInt([]byte{}, int64(a), 10)
}
