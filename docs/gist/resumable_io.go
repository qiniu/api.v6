package gist

// @gist import
import resumable_io "github.com/qiniu/api/resumable/io"
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
	extra := &resumable_io.PutExtra {
		Bucket: bucketName,
	}
	// @endgist
	
	var ret interface{}
	
	// @gist put
	buf := bytes.NewReader([]byte("data"))
	var fsize int64 = 4
	putErr := resumable_io.Put(logger, &ret, putPolicy.Token(), key, buf, fsize, extra)
	// @endgist
	
	// @gist put_file
	localFile := "<path/to/file>"
	putFileErr := resumable_io.PutFile(logger, &ret, putPolicy.Token(), key, localFile, extra)
	// @endgist
	
	_, _ = putFileErr, putErr
}
