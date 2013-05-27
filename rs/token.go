package rs

import (
	"time"
	"strings"
	"github.com/qiniu/api/auth/digest"
	. "github.com/qiniu/api/conf"
)

// ----------------------------------------------------------

type GetPolicy struct {
	Scope		string `json:"S"`
	Expires		uint32 `json:"E"`
}

func (r GetPolicy) token() string {
	if r.Expires == 0 {
		r.Expires = 3600
	}
	r.Expires += uint32(time.Now().Unix())
	return digest.SignJson(ACCESS_KEY, []byte(SECRET_KEY), &r)
}

func (r GetPolicy) MakeRequest(baseUrl string) (privateUrl string) {

	token := r.token()
	if strings.Contains(baseUrl, "?") {
		return baseUrl + "&token=" + token
	}
	return baseUrl + "?token=" + token
}

// --------------------------------------------------------------------------------

type PutPolicy struct {
	Scope            string `json:"scope,omitempty"`
	CallbackUrl      string `json:"callbackUrl,omitempty"`
	CallbackBody     string `json:"callbackBody,omitempty"`
	ReturnUrl        string `json:"returnUrl,omitempty"`
	ReturnBody       string `json:"returnBody,omitempty"`
	AsyncOps         string `json:"asyncOps,omitempty"`
	EndUser          string `json:"endUser,omitempty"`
	Expires          uint32 `json:"deadline"` 			// 截止时间（以秒为单位）
}

func (r *PutPolicy) Token() string {
	if r.Expires == 0 {
		r.Expires = 3600
	}
	r.Expires += uint32(time.Now().Unix())
	return digest.SignJson(ACCESS_KEY, []byte(SECRET_KEY), &r)
}

// ----------------------------------------------------------

