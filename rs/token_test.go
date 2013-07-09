package rs

import (
	"crypto/sha1"
	"encoding/base64"
	"io"
	"net/http"
	"os"
	"testing"
	. "github.com/qiniu/api/conf"
)

func init() {

	ACCESS_KEY = os.Getenv("QINIU_ACCESS_KEY")
	SECRET_KEY = os.Getenv("QINIU_SECRET_KEY")
	if ACCESS_KEY == "" || SECRET_KEY == "" {
		panic("require ACCESS_KEY & SECRET_KEY")
	}
	client = New(nil)
	client.Delete(nil, bucketName, key)
}

func TestGetPrivateUrl(t *testing.T) {

	//上传一个文件用用于测试
	err := upFile("token.go", bucketName, key)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Delete(nil, bucketName, key)

	baseUrl := MakeBaseUrl(domain, key)

	policy := GetPolicy{}
	privateUrl := policy.MakeRequest(baseUrl, nil)

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
