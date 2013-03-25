package fop

import (
	"os"
	
	. "github.com/qiniu/api/conf"
)

func init() {
	ACCESS_KEY = os.Getenv("QINIU_ACCESS_KEY")
	SECRET_KEY = os.Getenv("QINIU_SECRET_KEY")
}

