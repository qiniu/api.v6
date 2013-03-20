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

