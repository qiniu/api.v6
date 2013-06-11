package gist

// @gist import
import "github.com/qiniu/api/io"
// @endgist

import "bytes"
import "github.com/qiniu/api/rs"
import "github.com/qiniu/rpc"


func ioDemo() {
	var logger rpc.Logger
	var bucketName = "<bucketName>"

	// @gist put_policy
	putPolicy := rs.PutPolicy {
		Scope: bucketName,
	}
	// @endgist

	// @gist put_extra
	extra := &io.PutExtra {
		Bucket: bucketName,
	}
	// @endgist

	var ret interface{}

	// @gist put
	buf := bytes.NewBufferString("data")
	uptoken := putPolicy.Token(nil)
	putErr := io.Put(logger, &ret, uptoken, "<key>", buf, extra)
	// @endgist

{	// @gist put_file
	localFile := "<path/to/file>"
	uptoken := putPolicy.Token(nil)
	putFileErr := io.PutFile(logger, &ret, uptoken, "<key>", localFile, extra)
	// @endgist

	_, _ = putErr, putFileErr
}}

func download() {

	// @gist download
	baseUrl := rs.MakeBaseUrl("<domain>", "<key>")
	policy := rs.GetPolicy{}
	downloadUrl := policy.MakeRequest(baseUrl, nil)
	// @endgist

	_ = downloadUrl
}

