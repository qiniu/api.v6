package rs

import (
	"os"
	"fmt"
	"io"
	"testing"
	"crypto/sha1"
	"encoding/base64"
	"net/http"
	. "github.com/qiniu/api/conf"
)

var key = "aa"
var newkey1 = "bbbb"
var newkey2 = "cccc"
var bucketName = "a"
var domain = "aatest.qiniudn.com"
var client Client

func init() {
	ACCESS_KEY = os.Getenv("QINIU_ACCESS_KEY")
	SECRET_KEY = os.Getenv("QINIU_SECRET_KEY")
	if ACCESS_KEY == "" || SECRET_KEY == "" {
		panic("require ACCESS_KEY & SECRET_KEY")
	}
	client = New(nil)

	// 删除 可能存在的 newkey1  newkey2 
	delFile(newkey1)
	delFile(newkey2)
}

func delFile(key string) {
	if _, err := client.Stat(nil,bucketName, key); err == nil {
		client.Delete(nil, bucketName, key)
	}

}

func TestGetPrivateUrl(t *testing.T) {

	baseUrl := MakeBaseUrl(domain, key)

	policy := GetPolicy{}
	privateUrl := policy.MakeRequest(baseUrl, nil)
	fmt.Println("privateUrl:", privateUrl)

	resp, err := http.Get(privateUrl)
	if err != nil {
		t.Fatal("http.Get failed:", err)
	}
	defer resp.Body.Close()

	h := sha1.New()
	io.Copy(h, resp.Body)
	etagExpected := base64.URLEncoding.EncodeToString(h.Sum([]byte{'\x16'}))

	etag := resp.Header.Get("Etag")
	if etag[1:len(etag)-1] != etagExpected {
		t.Fatal("http.Get etag failed:", etag, etagExpected)
	}
}

func TestEntry(t *testing.T) {

	einfo, err := client.Stat(nil, bucketName, key)
	if err != nil {
		t.Fatal(err)
	}

	err = client.Copy(nil, bucketName, key, bucketName, newkey1)
	if err != nil {
		t.Fatal(err)
	}
	enewinfo, err := client.Stat(nil, bucketName, newkey1)
	if err != nil {
		t.Fatal(err)
	}
	if einfo.Hash != enewinfo.Hash {
		t.Fatal("invalid entryinfo:", einfo, enewinfo)
	}

	err = client.Move(nil, bucketName, newkey1, bucketName, newkey2)
	if err != nil {
		t.Fatal(err)
	}
	enewinfo2, err := client.Stat(nil, bucketName, newkey2)
	if err != nil {
		t.Fatal(err)
	}
	if enewinfo.Hash != enewinfo2.Hash {
		t.Fatal("invalid entryinfo:", enewinfo, enewinfo2)
	}

	err = client.Delete(nil, bucketName, newkey2)
	if err != nil {
		t.Fatal(err)
	}
}

func TestBatchStat(t *testing.T) {

	entryPath := EntryPath{
		Bucket: bucketName,
		Key:    key,
	}

	rets, err := client.BatchStat(nil, []EntryPath{entryPath, entryPath, entryPath})
	if err != nil {
		t.Fatal(err)
	}

	if len(rets) != 3 {
		t.Fatal("BatchStat failed: len(result) = ", len(rets))
	}

	stat, _ := client.Stat(nil, bucketName, key)

	if rets[0].Data != stat || rets[1].Data != stat || rets[2].Data != stat {
		t.Error("BatchStat failed : returns err")
	}
}

func TestBatchMove(t *testing.T) {
	stat0, err := client.Stat(nil, bucketName, key)
	entryPair1 := EntryPathPair{
		Src: EntryPath{
			Bucket: bucketName,
			Key:    key,
		},
		Dest: EntryPath{
			Bucket: bucketName,
			Key:    newkey1,
		},
	}

	entryPair2 := EntryPathPair{
		Src: EntryPath{
			Bucket: bucketName,
			Key:    newkey1,
		},
		Dest: EntryPath{
			Bucket: bucketName,
			Key:    newkey2,
		},
	}

	_, err = client.BatchMove(nil, []EntryPathPair{entryPair1, entryPair2})
	if err != nil {
		t.Fatal(err)
	}
	defer client.Move(nil, bucketName, newkey2, bucketName, key)

	stat1, _ := client.Stat(nil, bucketName, newkey2)

	if stat0 != stat1 {
		t.Error("BatchMove failed : Move err")
	}
}

func TestBatchCopy(t *testing.T) {
	entryPair1 := EntryPathPair{
		Src: EntryPath{
			Bucket: bucketName,
			Key:    key,
		},
		Dest: EntryPath{
			Bucket: bucketName,
			Key:    newkey1,
		},
	}

	entryPair2 := EntryPathPair{
		Src: EntryPath{
			Bucket: bucketName,
			Key:    newkey1,
		},
		Dest: EntryPath{
			Bucket: bucketName,
			Key:    newkey2,
		},
	}

	_, err := client.BatchCopy(nil, []EntryPathPair{entryPair1, entryPair2})
	if err != nil {
		t.Fatal(err)
	}
	defer client.Delete(nil, bucketName, newkey1)
	defer client.Delete(nil, bucketName, newkey2)

	stat0, _ := client.Stat(nil, bucketName, key)
	stat1, _ := client.Stat(nil, bucketName, newkey1)
	stat2, _ := client.Stat(nil, bucketName, newkey2)
	if stat0.Hash != stat1.Hash || stat0.Hash != stat2.Hash {
		t.Error("BatchCopy failed : Copy err")
	}
}

func TestBatchDelete(t *testing.T) {
	client.Copy(nil, bucketName, key, bucketName, newkey1)
	client.Copy(nil, bucketName, key, bucketName, newkey2)

	entryPath1 := EntryPath{
		Bucket: bucketName,
		Key:    newkey1,
	}
	entryPath2 := EntryPath{
		Bucket: bucketName,
		Key:    newkey2,
	}

	_, err := client.BatchDelete(nil, []EntryPath{entryPath1, entryPath2})
	if err != nil {
		t.Fatal(err)
	}

	_, err1 := client.Stat(nil, bucketName, newkey1)
	_, err2 := client.Stat(nil, bucketName, newkey2)

	//这里 err1 != nil，否则文件没被成功删除
	if err1 == nil || err2 == nil {
		t.Error("BatchDelete failed : File do not delete")
	}
}

func TestBatch(t *testing.T) {
	ops := []string{
		URICopy(bucketName, key, bucketName, newkey1),
		URIDelete(bucketName, key),
		URIMove(bucketName, newkey1, bucketName, key),
	}
	rets := new([]BatchItemRet)
	err := client.Batch(nil, rets, ops)
	if err != nil {
		t.Fatal(err)
	}
}
