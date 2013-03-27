package fop

import (
	"strconv"
	"encoding/base64"
)

// ----------------------------------------------------------

func encodeURI(uri string) string {
	return base64.URLEncoding.EncodeToString([]byte(uri))
}

func itoa(a int) []byte {
	return strconv.AppendInt([]byte{}, int64(a), 10)
}
