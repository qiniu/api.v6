package gist

import (
	gio"io"
	"github.com/qiniu/rpc"
	"github.com/qiniu/api/io"
	rio "github.com/qiniu/api/resumable/io"
)

// @gist uploadFile
func uploadFile(l rpc.Logger, uptoken, key, localFile string) (ret io.PutRet, err error) {
	err = io.PutFile(l, &ret, uptoken, key, localFile, nil)
	return
}
// @endgist

// @gist simpleUploadFile
func simpleUploadFile(l rpc.Logger, uptoken, key, localFile string) error {
	return io.PutFile(l, nil, uptoken, key, localFile, nil)
}
// @endgist

// @gist uploadBuf
func uploadBuf(l rpc.Logger, uptoken, key string, r gio.Reader) (ret io.PutRet, err error) {
	err = io.Put(l, &ret, uptoken, key, r, nil)
	return
}
// @endgist

// @gist resumableUpload
func resumableUpload(l rpc.Logger, uptoken, key, localFile string) error {
	return rio.PutFile(l, nil, uptoken, key, localFile, nil)
}
// @endgist
