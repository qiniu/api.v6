package fop

import (
	"bytes"
	"errors"
)

// ----------------------------------------------------------

type ImageView struct {
	Mode    uint    // 1或2
	Width   uint
	Height  uint
	Quality uint    // 质量, 1-100
	Format  string  // 输出格式, jpg, gif, png, tif 等图片格式
}

func (this ImageView) MakeRequest(url string) (reqUrl string, err error) {
	buf := bytes.NewBuffer(make([]byte, 0, bytes.MinRead))
	if this.Mode != 1 && this.Mode != 2 {
		err = errors.New("mode only can be 1 or 2")
		return
	}
	buf.WriteByte('/')
	buf.Write(itoa(this.Mode))
	
	if this.Width > 0 {
		buf.WriteString("/w/")
		buf.Write(itoa(this.Width))
	}
	
	if this.Height > 0 {
		buf.WriteString("/h/")
		buf.Write(itoa(this.Height))
	}
	
	if this.Quality > 0 {
		buf.WriteString("/q/")
		buf.Write(itoa(this.Quality))
	}
	
	if this.Format != "" {
		buf.WriteString("/format/")
		buf.WriteString(this.Format)
	}
	
	reqUrl = url + "?imageView" + string(buf.Bytes())
	return
}
