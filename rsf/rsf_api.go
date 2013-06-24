package rsf

import (
	"net/http"
	"net/url"
	"strconv"
	"github.com/qiniu/rpc"
	"github.com/qiniu/api/auth/digest"
	. "github.com/qiniu/api/conf"
)

// ----------------------------------------------------------

type ListItem struct {
	Key      string `json:"key"`
	Hash     string `json:"hash"`
	Fsize    int64  `json:"fsize"`
	PutTime  int64  `json:"putTime"`
	MimeType string	`json:"mimeType"`
	EndUser  string	`json:"endUser"`
}

type ListRet struct{
	Marker string     `json:"marker"`
	Items  []ListItem `json:"items"`
}

// ----------------------------------------------------------

type Client struct {
	Conn rpc.Client
}

func New(mac *digest.Mac) Client {
	t := digest.NewTransport(mac, nil)
	client := &http.Client{Transport: t}
	return Client{rpc.Client{client}}
}

func NewEx(t http.RoundTripper) Client {
	client := &http.Client{Transport: t}
	return Client{rpc.Client{client}}
}

// ----------------------------------------------------------

func (rsf Client) ListPrefix(l rpc.Logger, bucket, prefix, marker string, limit int) (entries []ListItem, markerOut string, err error) {
	URL := makeListURL(bucket, prefix, marker, limit)
	listRet := ListRet{}
	err = rsf.Conn.Call(l, &listRet, URL)
	return listRet.Items, listRet.Marker, err
}

func makeListURL(bucket, prefix, marker string, limit int) string {
	query := make(url.Values)
	if bucket != "" {
		query.Add("bucket", bucket)
	}
	if prefix != "" {
		query.Add("prefix", prefix)
	}
	if marker != "" {
		query.Add("marker", marker)
	}
	if limit > 0 {
		query.Add("limit", strconv.FormatInt(int64(limit), 10))
	}

	return RSF_HOST + "/list?" + query.Encode()
}
