package main

// @gist import
import qiniu_io "github.com/qiniu/api/io"
// @endgist

import "github.com/qiniu/api/rs"

import "bytes"
import "github.com/qiniu/rpc"

var logger rpc.Logger
var bucketName = "<bucketName>"
var key = "<key>"

func main() {
	// @gist put_policy
	putPolicy := rs.PutPolicy {
		Scope: bucketName,
	}
	// @endgist
	
	
	// @gist put_extra
	extra := &qiniu_io.PutExtra {
		Bucket: bucketName,
	}
	// @endgist
	
	var ret interface{}
	
	// @gist put
	buf := bytes.NewBufferString("data")
	putErr := qiniu_io.Put(logger, &ret, putPolicy.Token(), "<key>", buf, extra)
	// @endgist
	
	// @gist put_file
	localFile := "<path/to/file>"
	putFileErr := qiniu_io.PutFile(logger, &ret, putPolicy.Token(), "<key>", localFile, extra)
	// @endgist
	_, _ = putErr, putFileErr
}

func download() {
	// @gist download
	policy := rs.GetPolicy {
		Scope: "<bucketName>",
	}
	// 生成下载连接, sourceUrl 为资源原有下载链接
	downloadUrl := qiniu_io.GetUrl("<domain>", "<key>", policy.Token())
	// @endgist
	
	_ = downloadUrl
}
