package main

// @gist import
import "github.com/qiniu/api/rs"
// @endgist

import "github.com/qiniu/rpc"

var logger rpc.Logger
var bucketName = "<bucketName>"
var key = "<key>"

var bucket1 = "<bucketName1>"
var key1 = "<key1>"
var bucket2 = "<bucketName2>"
var key2 = "<key2>"
var bucket3 = "<bucketName3>"
var key3 = "<key3>"
var bucket4 = "<bucketName4>"
var key4 = "<key4>"

var bucketSrc = "<bucketSrc>"
var keySrc = "<keySrc>"
var bucketDest = "<bucketDest>"
var keyDest = "<keyDest>"

func main() {
	// @gist stat
	rs.New().Stat(logger, bucketName, key) // 返回: rs.Entry, error
	// @endgist
	
	// @gist copy
	// 返回值 error, 操作成功时err为nil
	rs.New().Copy(logger, bucketSrc, keySrc, bucketDest, keyDest)
	// @endgist
	
	// @gist move
	// 返回值 error, 操作成功时err为nil
	rs.New().Move(logger, bucketSrc, keySrc, bucketDest, keyDest) 
	// @endgist
	
	// @gist delete
	rs.New().Delete(logger, bucketName, key) // 返回值 error, 操作成功时err为nil
	// @endgist
}

func batch() {
	// @gist entry_pathes
	entryPathes := []rs.EntryPath {
		rs.EntryPath {
			Bucket: bucket1,
			Key: key1,
		},
		rs.EntryPath {
			Bucket: bucket2,
			Key: key2,
		},
	}
	// @endgist
	
	// @gist entry_path_pairs
	// 每个复制操作都含有源文件和目标文件
	entryPairs := []rs.EntryPathPair {
		rs.EntryPathPair {
			Src: rs.EntryPath {
				Bucket: bucket1,
				Key: key1,
			},
			Dest: rs.EntryPath {
				Bucket: bucket2,
				Key: key2,
			},
		}, rs.EntryPathPair {
			Src: rs.EntryPath {
				Bucket: bucket3,
				Key: key3,
			},
			Dest: rs.EntryPath {
				Bucket: bucket4,
				Key: key4,
			},
		},
	}
	// @endgist
	
	// @gist batch_stat
	rs.New().BatchStat(logger, entryPathes) // []rs.BatchStatItemRet, error
	// @endgist
	
	
	// @gist batch_copy
	rs.New().BatchCopy(logger, entryPairs) 
	// []rs.BatchResult, error
	// @endgist
	
	// @gist batch_move
	rs.New().BatchMove(logger, entryPairs)
	// []rs.BatchResult, error
	// @endgist
	
	// @gist batch_delete
	rs.New().BatchDelete(logger, entryPathes)
	// []rs.BatchResult, error
	// @endgist
	
	// @gist batch_adv
	ops := []string {
		rs.URIStat(bucketName, key1),
		rs.URICopy(bucketName, key1, bucketName, key2), // 复制key1到key2
		rs.URIDelete(bucketName, key1), // 删除key1
		rs.URIMove(bucketName, key2, bucketName, key1), //将key2移动到key1
	}
	rets := new([]rs.BatchItemRet)
	rs.New().Batch(logger, rets, ops) // 执行操作, 返回error
	// @endgist
}

func token() {
	// @gist put_policy
	putPolicy := rs.PutPolicy {
		Scope: bucketName,
	}
	putPolicy.Token() // UpToken
	// @endgist
	
	// @gist get_policy
	getPolicy := rs.GetPolicy {
		Scope: bucketName,
	}
	getPolicy.Token() // DnToken
	// @endgist
}


