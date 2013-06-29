package gist

// @gist import
import "github.com/qiniu/api/rsf"
// @endgist

var limit = 1000
var prefix = "<prefix>"
var marker = "<marker>"

func listPrefixDemo() {
	// @gist listPrefix
	// 返回 []ListItem, err
	rsf.New(nil).ListPrefix(logger, bucketName, prefix, marker, limit)
	// @endgist
}

