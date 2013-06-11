package gist

// @gist import
import "github.com/qiniu/api/resumable/io"
// @endgist

import "github.com/qiniu/rpc"
import "github.com/qiniu/api/rs"
import "bytes"

func resumableIoDemo() {
	var logger rpc.Logger
	var bucketName = "<bucketName>"
	var key = "<key>"

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
	buf := bytes.NewReader([]byte("data"))
	fsize := int64(buf.Len())
	uptoken := putPolicy.Token(nil)
	putErr := io.Put(logger, &ret, uptoken, key, buf, fsize, extra)
	// @endgist
{
	// @gist put_file
	localFile := "<path/to/file>"
	uptoken := putPolicy.Token(nil)
	putFileErr := io.PutFile(logger, &ret, uptoken, key, localFile, extra)
	// @endgist

	_, _ = putFileErr, putErr
}}

