package rsf

import (
	"errors"
	"net/http"
	"strconv"
	"io"
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
// 1. 首次请求 marker = ""
// 2. 无论 err 值如何，均应该先看 entries 是否有内容
// 3. 如果后续没有更多数据，err 返回 EOF，markerOut 返回 ""（但不通过该特征来判断是否结束）
func (rsf Client) ListPrefix(l rpc.Logger, bucket, prefix, marker string, limit int) (entries []ListItem, markerOut string, err error) {

	if bucket == "" {
		err = errors.New("bucket could not be nil")
		return 
	}

	params := map[string][]string{
		"bucket": {bucket},
	}
	if prefix != "" {
		params["prefix"] = []string{prefix}
	}
	if marker != "" {
		params["marker"] = []string{marker}
	}
	if limit > 0 {
		params["limit"] = []string{strconv.FormatInt(int64(limit), 10)}
	}
	listRet := ListRet{}
	err = rsf.Conn.CallWithForm(l, &listRet, RSF_HOST + "/list", params)
	if err != nil {
		return 
	}
	if listRet.Marker == "" {
		return listRet.Items, "", io.EOF
	}
	return listRet.Items, listRet.Marker, err
}
