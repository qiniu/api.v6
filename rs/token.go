package rs

import (
	"time"
	"github.com/qiniu/api/auth/digest"
	. "github.com/qiniu/api/conf"
)

// ----------------------------------------------------------

type GetPolicy struct {
	Scope		string `json:"S"`
	Expires		uint32 `json:"E"`
}

func (r GetPolicy) Token() string {
	if r.Expires == 0 {
		r.Expires = 3600
	}
	r.Expires += uint32(time.Now().Unix())
	return digest.SignJson(ACCESS_KEY, []byte(SECRET_KEY), &r)
}

// ----------------------------------------------------------

type PutPolicy struct {
	Scope            string `json:"scope,omitempty"`
	CallbackUrl      string `json:"callbackUrl,omitempty"`
	CallbackBodyType string `json:"callbackBodyType,omitempty"`
	Customer         string `json:"customer,omitempty"`
	AsyncOps         string `json:"asyncOps,omitempty"`
	Expires          uint32 `json:"deadline"` 			// 截止时间（以秒为单位）
	Escape           uint16 `json:"escape,omitempty"`	// 是否允许存在转义符号
	DetectMime       uint16	`json:"detectMime",omitempty`
}

func (r PutPolicy) Token() string {
	if r.Expires == 0 {
		r.Expires = 3600
	}
	r.Expires += uint32(time.Now().Unix())
	return digest.SignJson(ACCESS_KEY, []byte(SECRET_KEY), &r)
}

// ----------------------------------------------------------

