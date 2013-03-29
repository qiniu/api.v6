package fop

import (
	"strconv"
	"encoding/base64"
)

// ----------------------------------------------------------

func encodeURI(uri string) string {
	return base64.URLEncoding.EncodeToString([]byte(uri))
}

func itoa(a int) string {
	return string(strconv.AppendInt([]byte{}, int64(a), 10))
}
