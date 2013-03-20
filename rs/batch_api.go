package rs

import (
	"github.com/qiniu/rpc"
	. "github.com/qiniu/api/conf"
)

// ----------------------------------------------------------

type BatchItemRet struct {
	Data  interface{} `json:"data"`
	Code  int         `json:"code"`
	Error string      `json:"error"`
}

type Batcher struct {
	op []string
	ret []BatchItemRet
}

func NewBatcher(ncap int) *Batcher {
	return &Batcher{
		op: make([]string, 0, ncap),
		ret: make([]BatchItemRet, 0, ncap),
	}
}

func (b *Batcher) Reset() {
	b.op = nil
	b.ret = nil
}

func (b *Batcher) Len() int {
	return len(b.op)
}

func (b *Batcher) AddItem(ret interface{}, op string) {
	b.op = append(b.op, op)
	b.ret = append(b.ret, BatchItemRet{Data: ret})
}

func (b *Batcher) Stat(bucket, key string) {
	b.AddItem(Entry{}, URIStat(bucket, key))
}

func (b *Batcher) Delete(bucket, key string) {
	b.AddItem(nil, URIDelete(bucket, key))
}

func (b *Batcher) Move(bucketSrc, keySrc, bucketDest, keyDest string) {
	b.AddItem(nil, URIMove(bucketSrc, keySrc, bucketDest, keyDest))
}

func (b *Batcher) Copy(bucketSrc, keySrc, bucketDest, keyDest string) {
	b.AddItem(nil, URICopy(bucketSrc, keySrc, bucketDest, keyDest))
}

func (b *Batcher) Do(l rpc.Logger, rs Client) (ret []BatchItemRet, err error) {
	err = rs.Conn.CallWithForm(l, &b.ret, RS_HOST+"/batch", map[string][]string{"op": b.op})
	ret = b.ret
	return
}

// ----------------------------------------------------------

type EntryPath struct {
	Bucket string
	Key string
}

type EntryPathPair struct {
	Src EntryPath
	Dest EntryPath
}

func (rs Client) BatchStat(l rpc.Logger, entries []EntryPath) (ret []BatchItemRet, err error) {
	b := NewBatcher(len(entries))
	for _, e := range entries {
		b.Stat(e.Bucket, e.Key)
	}
	return b.Do(l, rs)
}

func (rs Client) BatchDelete(l rpc.Logger, entries []EntryPath) (ret []BatchItemRet, err error) {
	b := NewBatcher(len(entries))
	for _, e := range entries {
		b.Delete(e.Bucket, e.Key)
	}
	return b.Do(l, rs)
}

func (rs Client) BatchMove(l rpc.Logger, entries []EntryPathPair) (ret []BatchItemRet, err error) {
	b := NewBatcher(len(entries))
	for _, e := range entries {
		b.Move(e.Src.Bucket, e.Src.Key, e.Dest.Bucket, e.Dest.Key)
	}
	return b.Do(l, rs)
}

func (rs Client) BatchCopy(l rpc.Logger, entries []EntryPathPair) (ret []BatchItemRet, err error) {
	b := NewBatcher(len(entries))
	for _, e := range entries {
		b.Copy(e.Src.Bucket, e.Src.Key, e.Dest.Bucket, e.Dest.Key)
	}
	return b.Do(l, rs)
}

// ----------------------------------------------------------

