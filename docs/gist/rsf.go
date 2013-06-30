package gist

// @gist import
import "github.com/qiniu/api/rsf"
// @endgist

import (
	"io"
	"fmt"
)

func listPrefixDemo() {
	// @gist listPrefix
	var entries []rsf.ListItem
	var marker = ""
	var err error
	var limit = 1000
	var prefix = "<prefix>"

	for err != io.EOF {
		entries, marker, err = rsf.New(nil).ListPrefix(logger, bucketName,
			prefix, marker, limit)
		if err != nil && err != io.EOF {
			break
		}

		for _, v := range entries {
			//这里为了演示，依次输出文件列表
			fmt.Println(v)
		}
	}
	// @endgist
}

