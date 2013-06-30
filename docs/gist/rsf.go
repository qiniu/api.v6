package gist

import (
	"io"
	"fmt"
	"github.com/qiniu/rpc"
	"github.com/qiniu/api/rsf"
	. "github.com/qiniu/api/conf"
)

func init() {
	// @gist init
	ACCESS_KEY = "<YOUR_APP_ACCESS_KEY>"
	SECRET_KEY = "<YOUR_APP_SECRET_KEY>"
	// @endgist
}

// @gist listPrefix
func list(l rpc.Logger, rs *rsf.Client, bucketName string, prefix string) {

	var entries []rsf.ListItem
	var marker = ""
	var err error
	var limit = 1000

	for err == nil {
		entries, marker, err = rs.ListPrefix(logger, bucketName,
			prefix, marker, limit)
		for _, item := range entries {
			//这里为了演示，依次输出文件列表
			fmt.Println(item)
		}
	}
	if err != io.EOF {
		//非预期的错误
		fmt.Println("err:", err)
	}
}
// @endgist

