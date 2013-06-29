package rsf 

import (
	"os"
	"testing"
	"io"
	. "github.com/qiniu/api/conf"
)

var bucketName = "a"
var client Client
var maxNum = 100000

func init() {
	ACCESS_KEY = os.Getenv("QINIU_ACCESS_KEY")
	SECRET_KEY = os.Getenv("QINIU_SECRET_KEY")
	if ACCESS_KEY == "" || SECRET_KEY == "" {
		panic("require ACCESS_KEY & SECRET_KEY")
	}
	client = New(nil)
}

func TestList(t *testing.T) {

	ret, marker, err := client.ListPrefix(nil, bucketName, "", "", 5)
	if err != nil {
		t.Fatal(err)
	}
	if len(ret) != 5 && err != io.EOF {
		t.Fatal("TestList failed:", "expect len(ret) 5, but:", len(ret))
	}

	ret, _, err = client.ListPrefix(nil, bucketName, "", marker, 10000)
	if err != nil && err != io.EOF {
		t.Fatal("TestList failed:", "marker failed:", err)
	}

}

func TestEof(t *testing.T) {
	_, _, err := client.ListPrefix(nil, bucketName, "", "", maxNum)

	if err != io.EOF {
		t.Fatal("TestEof failed:", "expect EOF but:", err)
	}

	_, _, err = client.ListPrefix(nil, bucketName, "", "", -1)

	if err != io.EOF {
		t.Fatal("TestEof failed:", "expect EOF but:", err)
	}
}
