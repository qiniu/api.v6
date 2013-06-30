package gist

// @gist import
import "github.com/qiniu/api/rsf"
// @endgist

import "io"

var limit = 1000
var prefix = "<prefix>"
var entries []rsf.ListItem
var marker = ""
var err error

func listPrefixDemo() {
	// @gist listPrefix
	// allItems 存储所有文件列表信息
	var allItems []rsf.ListItem
	for err != io.EOF {
		entries, marker, err = rsf.New(nil).ListPrefix(logger, bucketName,
			prefix, marker, limit)
		if err != nil && err != io.EOF {
			break
		}
		allItems = append(allItems, entries...)
	}
	// @endgist
}

