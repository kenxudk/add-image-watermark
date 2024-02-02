// Package watermark
package watermark

// RequestData 请求的json数据
type RequestData struct {
	Channel      string `json:"channel"`
	Name         string `json:"name"`
	NameYOffset  int    `json:"name_y_offset"` //name 在logo图片的y轴偏移
	NameXOffset  int    `json:"name_x_offset"` //name 在logo图片的x轴偏移
	Key          string `json:"key"`
	FileW        int    `json:"file_w"`        //key 宽
	FileH        int    `json:"file_h"`        //key 高
	FontSize     string `json:"fontsize"`      //字体大小。默认15
	FontColor    string `json:"fontcolor"`     //字体颜色，默认白色
	LogoLocation int    `json:"logo_location"` // 0 随机，1左上角，2右上角，3左下角，4右下角
	LogoYOffset  int    `json:"logo_y_offset"` //logo y轴偏移
	LogoXOffset  int    `json:"logo_x_offset"` //logo x轴偏移
}
