package rs

import (
	"bytes"
	"net/http"
	"io"
	"strings"
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
	err = rs.Conn.Call(l, &entry, RS_HOST + URIStat(entryURI))
	return
}

// ----------------------------------------------------------

func (rs Client) Delete(l rpc.Logger, entryURI string) (err error) {
	return rs.Conn.Call(l, nil, RS_HOST + URIDelete(entryURI))
}

func (rs Client) Move(l rpc.Logger, entryURISrc, entryURIDest string) (err error) {
	return rs.Conn.Call(l, nil, RS_HOST + URIMove(entryURISrc, entryURIDest))
}

func (rs Client) Copy(l rpc.Logger, entryURISrc, entryURIDest string) (err error) {
	return rs.Conn.Call(l, nil, RS_HOST + URICopy(entryURISrc, entryURIDest))
}

// -----------------------------------------------------------------------------

type BatchResult struct {
	Code int `json:"code"`
	Data interface{} `json:"data"`
	Error string `json:"-"`
}

type BatchStatResult struct {
	Code int `json:"code"`
	Data Entry `json:"data"`
	Error string `json:"-"`
}

type Batcher struct {
	*bytes.Buffer
}

func NewBatcher() *Batcher {
	return &Batcher{bytes.NewBuffer(nil)}
}

func (this *Batcher) Add(uris ...string) {
	for _, uri := range uris {
		if this.Len() > 0 {
			this.WriteByte('&')
		}
		
		io.WriteString(this, "op=" + strings.Replace(uri, "/", "%2F", -1))
	}
}

func (rs Client) DoBatch(l rpc.Logger,
	ret interface{}, b *Batcher) (err error) {

	bodyType := "application/x-www-form-urlencoded"
	err = rs.Conn.CallWith(l, &ret, RS_HOST + "/batch", bodyType, b, b.Len())
	return
}

func (rs Client) BatchStat(l rpc.Logger,
	entryURIs ...string) (ret []BatchStatResult, err error){
	
	b := NewBatcher()
	for _, entryURI := range entryURIs{
		b.Add(URIStat(entryURI))
	}
	err = rs.DoBatch(l, &ret, b)
	return
}

func (rs Client) BatchDelete(l rpc.Logger,
	entryURIs ...string) (ret []BatchResult, err error) {
	
	b := NewBatcher()
	for _, entryURI := range entryURIs{
		b.Add(URIDelete(entryURI))
	}
	err = rs.DoBatch(l, &ret, b)	
	return
}

type EntryURIPath struct {
	Src string
	Dest string
}

func (rs Client) BatchCopy(l rpc.Logger,
	pathes ...EntryURIPath) (ret []BatchResult, err error) {
	
	b := NewBatcher()
	for _, path := range pathes {
		b.Add(URICopy(path.Src, path.Dest))
	}
	err = rs.DoBatch(l, &ret, b)
	return
}

func (rs Client) BatchMove(l rpc.Logger,
	pathes ...EntryURIPath) (ret []BatchResult, err error) { 
	
	b := NewBatcher()
	for _, path := range pathes {
		b.Add(URIMove(path.Src, path.Dest))
	}
	err = rs.DoBatch(l, &ret, b)
	return
}
// -----------------------------------------------------------------------------

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

// -----------------------------------------------------------------------------

func EncodeURI(uri string) string {
	return base64.URLEncoding.EncodeToString([]byte(uri))
}

func URIDelete(entryURI string) string {
	return "/delete/" + EncodeURI(entryURI)
}

func URIStat(entryURI string) string {
	return "/stat/" + EncodeURI(entryURI)
}

func URICopy(entryURISrc, entryURIDest string) string {
	return "/copy/" + EncodeURI(entryURISrc) + "/" + EncodeURI(entryURIDest)
}

func URIMove(entryURISrc, entryURIDest string) string {
	return "/move/" + EncodeURI(entryURISrc) + "/" + EncodeURI(entryURIDest)
}

// -----------------------------------------------------------------------------
