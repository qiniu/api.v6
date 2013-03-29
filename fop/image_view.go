package fop

type ImageView struct {
	Mode int    // 缩略模式
	Width int    // Width = 0 表示不限定宽度
	Height int    // Height = 0 表示不限定高度
	Quality int    // 质量, 1-100
	Format string  // 输出格式，如jpg, gif, png, tif等等
}

func (this ImageView) MakeRequest(url string) (reqUrl string, err error) {
	query := "?imageView/" + itoa(this.Mode)
	
	if this.Width > 0 {
		query += "/w/" + itoa(this.Width)
	}
	
	if this.Height > 0 {
		query += "/h/" + itoa(this.Height)
	}
	
	if this.Quality > 0 {
		query += "/q/" + itoa(this.Quality)
	}
	
	if this.Format != "" {
		query += "/format/" + this.Format
	}
	
	reqUrl = url + query
	return
}
