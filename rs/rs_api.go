package rs

import (
	"net/http"
	"encoding/base64"
	"github.com/qiniu/rpc"
	"github.com/qiniu/api/auth/digest"
	. "github.com/qiniu/api/conf"
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

type Entry struct {
	Hash     string `json:"hash"`
	Fsize    int64  `json:"fsize"`
	PutTime  int64  `json:"putTime"`
	MimeType string `json:"mimeType"`
	Customer string `json:"customer"`
}

func (rs Client) Stat(l rpc.Logger, entryURI string) (entry Entry, err error) {
	err = rs.Conn.Call(l, &entry, RS_HOST+"/stat/"+EncodeURI(entryURI))
	return
}

// ----------------------------------------------------------

func (rs Client) Delete(l rpc.Logger, entryURI string) (err error) {
	return rs.Conn.Call(l, nil, RS_HOST+"/delete/"+EncodeURI(entryURI))
}

func (rs Client) Move(l rpc.Logger, entryURISrc, entryURIDest string) (err error) {
	return rs.Conn.Call(l, nil, RS_HOST+"/move/"+EncodeURI(entryURISrc)+"/"+EncodeURI(entryURIDest))
}

func (rs Client) Copy(l rpc.Logger, entryURISrc, entryURIDest string) (err error) {
	return rs.Conn.Call(l, nil, RS_HOST+"/copy/"+EncodeURI(entryURISrc)+"/"+EncodeURI(entryURIDest))
}

// ----------------------------------------------------------

func (rs Client) Mkbucket(l rpc.Logger, bucketName string) (err error) {
	return rs.Conn.Call(l, nil, RS_HOST+"/mkbucket/"+bucketName)
}

func (rs Client) Drop(l rpc.Logger, bucketName string) (err error) {
	return rs.Conn.Call(l, nil, RS_HOST+"/drop/"+bucketName)
}

func (rs Client) Buckets(l rpc.Logger) (buckets []string, err error) {
	err = rs.Conn.Call(l, &buckets, RS_HOST+"/buckets")
	return
}

// ----------------------------------------------------------

func EncodeURI(uri string) string {
	return base64.URLEncoding.EncodeToString([]byte(uri))
}

// ----------------------------------------------------------

