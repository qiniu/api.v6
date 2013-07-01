package gist

import (
	"log"
	gio"io"
	"github.com/qiniu/rpc"
	"github.com/qiniu/api/io"
	rio "github.com/qiniu/api/resumable/io"
)

func uploadFileDemo(localFile, key, uptoken string) {
// @gist uploadFile
	var logger rpc.Logger 
	var err error
	var ret io.PutRet
	var extra = &io.PutExtra {
		//Params:    params,
		//MimeType:  mieType,
		//Crc32:     crc32,
		//CheckCrc:  CheckCrc,
	}

	// logger    为rpc.Logger类型，日志参数,可选
	// ret       变量用于存取返回的信息，详情见 io.PutRet
	// uptoken   为业务服务器生成的上传口令
	// key       为文件存储的标识，当 key == "?"，则服务端自动生成key
	// localFile 为本地文件名
	// extra     为上传文件的额外信息，详情见 io.PutExtra，可选
	err = io.PutFile(logger, &ret, uptoken, key, localFile, extra)

	if err != nil {
	//上传产生错误
		log.Print("io.PutFile failed:", err)
		return
	}

	//上传成功，处理返回值
	log.Print(ret.Hash, ret.Key)
// @endgist
}

func uploadBufDemo( r gio.Reader, key, uptoken string) {
// @gist uploadBuf
	var logger rpc.Logger 
	var err error
	var ret io.PutRet
	var extra = &io.PutExtra {
		//Params:    params,
		//MimeType:  mieType,
		//Crc32:     crc32,
		//CheckCrc:  CheckCrc,
	}

	// logger    为rpc.Logger类型，日志参数,可选
	// ret       变量用于存取返回的信息，详情见 io.PutRet
	// uptoken   为业务服务器端生成的上传口令
	// key       为文件存储的标识，当 key == "?"，则服务端自动生成key
	// r         为io.Reader类型，用于从其读取数据
	// extra     为上传文件的额外信息,可为空， 详情见 io.PutExtra, 可选
	err = io.Put(logger, &ret, uptoken, key, r, extra)

	if err != nil {
	//上传产生错误
		log.Print("io.Put failed:", err)
		return
	}

	//上传成功，处理返回值
	log.Print(ret.Hash, ret.Key)
// @endgist
}

func resumableUploadFileDemo(localFile, key, uptoken string) {
// @gist resumableUploadFile
	var logger rpc.Logger 
	var err error
	var ret rio.PutRet
	var extra = &rio.PutExtra {
		//CallbackParams: callbackParams,
		//Bucket:         bucket,
		//CustomMeta:     customMeta,
		//MimeType:       mieType,
		//ChunkSize:      chunkSize,
		//TryTimes:       tryTimes,	
		//Progresses:     progresses,
		//Notify:         notify,		
		//NotifyErr:      NotifyErr,
	}

	// logger    为rpc.Logger类型，日志参数,可选
	// ret       变量用于存取返回的信息，详情见 resumable.io.PutRet
	// uptoken   为业务服务器生成的上传口令
	// key       为文件存储的标识
	// localFile 为本地文件名
	// extra     为上传文件的额外信息,可为空， 详情见 resumable.io.PutExtra
	err = rio.PutFile(logger, ret, uptoken, key, localFile, extra)

	if err != nil {
	//上传产生错误
		log.Print("resumable.io.Put failed:", err)
		return
	}

	//上传成功，处理返回值
	log.Print(ret.Hash)
// @endgist
}

func resumableUploadBufDemo(r gio.ReaderAt, fsize int64,  key, uptoken string) {
// @gist resumableUploadBuf
	var logger rpc.Logger 
	var err error
	var ret io.PutRet
	var extra = &rio.PutExtra {
		//CallbackParams: callbackParams,
		//Bucket:         bucket,
		//CustomMeta:     customMeta,
		//MimeType:       mieType,
		//ChunkSize:      chunkSize,
		//TryTimes:       tryTimes,	
		//Progresses:     progresses,
		//Notify:         notify,		
		//NotifyErr:      NotifyErr,
	}

	// logger    为rpc.Logger类型，日志参数,可选
	// ret       变量用于存取返回的信息，详情见 resumable.io.PutRet
	// uptoken   为业务服务器生成的上传口令
	// key       为文件存储的标识
	// r         为io.ReaderAt,用于读取数据
	// fsize     数据总字节数
	// extra     为上传文件的额外信息, 详情见 resumable.io.PutExtra
	err = rio.Put(logger, ret, uptoken, key, r, fsize, extra)

	if err != nil {
	//上传产生错误
		log.Print("resumable.io.Put failed:", err)
		return
	}

	//上传成功，处理返回值
	log.Print(ret.Hash)
// @endgist
}
