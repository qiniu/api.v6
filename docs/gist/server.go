package gist

import (
	"github.com/qiniu/api/rs"
	."github.com/qiniu/api/conf"
)

func init() {
	// @gist init
	ACCESS_KEY = "<YOUR_APP_ACCESS_KEY>"
	SECRET_KEY = "<YOUR_APP_SECRET_KEY>"
	// @endgist
}

// @gist uptoken
func uptoken(bucketName string) string {
	putPolicy := rs.PutPolicy {
		Scope: bucketName,
	}
	return  putPolicy.Token(nil)
}
// @endgist

// @gist downloadUrl
func downloadUrl(domain, key string) string {
	baseUrl := rs.MakeBaseUrl(domain, key)
	policy := rs.GetPolicy{}
	return  policy.MakeRequest(baseUrl, nil)
}
// @endgist
