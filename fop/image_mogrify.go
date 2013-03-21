package fop

import (
	"bytes"
	
	"github.com/qiniu/api/rs"
	"github.com/qiniu/rpc"
)
 
type ImageMogrify struct {
	AutoOrient bool
	Thumbnail string
	Gravity string
	Crop string
	Quality uint
	Rotate uint
	Format string
}

func (this ImageMogrify) Marshal() (uri string) {
	buf := bytes.NewBuffer(make([]byte, 0, bytes.MinRead))
	if this.AutoOrient {
		buf.WriteString("/auto-orient")
	}
	
	if this.Thumbnail != "" {
		buf.WriteString("/thumbnail/")
		buf.WriteString(this.Thumbnail)
	}
	
	if this.Gravity != "" {
		buf.WriteString("/gravity/")
		buf.WriteString(this.Gravity)
	}
	
	if this.Crop != "" {
		buf.WriteString("/crop/")
		buf.WriteString(this.Crop)
	}
	
	if this.Quality > 0 {
		buf.WriteString("/quality/")
		buf.Write(itoa(this.Quality))
	}
	
	if this.Rotate > 0 {
		buf.WriteString("/rotate/")
		buf.Write(itoa(this.Rotate))
	}
	
	if this.Format != "" {
		buf.WriteString("/format/")
		buf.WriteString(this.Format)
	}
	return string(buf.Bytes())
}

func (this ImageMogrify) MakeRequest(url string) string {
	return url + "?imageMogr" + this.Marshal()
}

func (this ImageMogrify) SaveAs(l rpc.Logger, 
	entryURISrc, entryURIDest string) (ret rs.Entry, err error) {

	fop := New()
	profile, err := fop.get(entryURISrc)
	if err != nil {
		return
	}
	encodedURIDest := encodeURI(entryURIDest)
	url := this.MakeRequest(profile.Url) + "/save-as/" + encodedURIDest
	err = fop.Conn.Call(l, &ret, url)
	
	return
}

