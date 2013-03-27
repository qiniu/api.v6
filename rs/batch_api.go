package rs

import (
	"github.com/qiniu/rpc"
	. "github.com/qiniu/api/conf"
)

// ----------------------------------------------------------

func (rs Client) Batch(l rpc.Logger, ret interface{}, op []string) (err error) {
	return rs.Conn.CallWithForm(l, ret, RS_HOST+"/batch", map[string][]string{"op": op})
}

// ----------------------------------------------------------

type BatchStatItemRet struct {
	Data  Entry       `json:"data"`
	Error string      `json:"error"`
	Code  int         `json:"code"`
}

type EntryPath struct {
	Bucket string
	Key string
}

func (rs Client) BatchStat(l rpc.Logger, entries []EntryPath) (ret []BatchStatItemRet, err error) {
	b := make([]string, len(entries))
	for i, e := range entries {
		b[i] = URIStat(e.Bucket, e.Key)
	}
	err = rs.Batch(l, &ret, b)
	return
}

// ----------------------------------------------------------

type BatchItemRet struct {
	Error string      `json:"error"`
	Code  int         `json:"code"`
}

func (rs Client) BatchDelete(l rpc.Logger, entries []EntryPath) (ret []BatchItemRet, err error) {
	b := make([]string, len(entries))
	for i, e := range entries {
		b[i] = URIDelete(e.Bucket, e.Key)
	}
	err = rs.Batch(l, &ret, b)
	return
}

// ----------------------------------------------------------

type EntryPathPair struct {
	Src EntryPath
	Dest EntryPath
}

func (rs Client) BatchMove(l rpc.Logger, entries []EntryPathPair) (ret []BatchItemRet, err error) {
	b := make([]string, len(entries))
	for i, e := range entries {
		b[i] = URIMove(e.Src.Bucket, e.Src.Key, e.Dest.Bucket, e.Dest.Key)
	}
	err = rs.Batch(l, &ret, b)
	return
}

func (rs Client) BatchCopy(l rpc.Logger, entries []EntryPathPair) (ret []BatchItemRet, err error) {
	b := make([]string, len(entries))
	for i, e := range entries {
		b[i] = URICopy(e.Src.Bucket, e.Src.Key, e.Dest.Bucket, e.Dest.Key)
	}
	err = rs.Batch(l, &ret, b)
	return
}

// ----------------------------------------------------------

