package fop

import (
	"bytes"
)
 
type ImageMogrify struct {
	AutoOrient bool
	Thumbnail string
	Gravity string
	Crop string
	Quality int
	Rotate int
	Format string
}

func (this ImageMogrify) marshal() (uri string) {
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
	return url + "?imageMogr" + this.marshal()
}
