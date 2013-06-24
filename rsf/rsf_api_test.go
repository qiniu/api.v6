package rsf 

import (
	"os"
	"testing"
	. "github.com/qiniu/api/conf"
)

var bucketName = "a"
var client Client

func init() {
	ACCESS_KEY = os.Getenv("QINIU_ACCESS_KEY")
	SECRET_KEY = os.Getenv("QINIU_SECRET_KEY")
	if ACCESS_KEY == "" || SECRET_KEY == "" {
		panic("require ACCESS_KEY & SECRET_KEY")
	}
	client = New(nil)
}

func TestList(t *testing.T) {

	ret, marker, err := client.ListPrefix(bucketName, "", "", 5)
	if err != nil {
		t.Fatal(err)
	}
	if len(ret) != 5 {
		t.Fatal("ret num err")
	}
	ret, _, err = client.ListPrefix(bucketName, "", marker, 100)
	if err != nil {
		t.Fatal(err)
	}
}


