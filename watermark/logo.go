package watermark

import (
	"github.com/golang/freetype"
	"github.com/nfnt/resize"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"strconv"
)

// 图片下载和生成的存放目录
var savePath string = "/tmp"

// TextInfo 文字信息
type TextInfo struct {
	Text        string  // 文字内容
	Size        float64 // 文字大小
	Color       string  //颜色，十六进制 #FF0000
	TextXOffset int     // x偏移位置信息
	TextYOffset int     // Y偏移位置信息
}

func GetLogoImage(imageH int, logoUrl string) (image.Image, bool, error) {
	watermark, err := os.Open(logoUrl)
	if err != nil {
		return nil, false, err
	}
	defer watermark.Close()
	imgWatermark, err := png.Decode(watermark)
	if err != nil {
		return nil, false, err
	}
	logoH := imgWatermark.Bounds().Dy()
	if imageH < 200 {
		scaleH := int(float64(logoH) * 0.7)
		return resize.Thumbnail(uint(imgWatermark.Bounds().Dx()), uint(scaleH), imgWatermark, resize.Bicubic), true, nil
	} else {
		return imgWatermark, false, nil
	}
}

// AddTextToLogo : 生成带文字的logo png格式图片
func (t TextInfo) AddTextToLogo(imgSource *os.File) string {
	var originalImg image.Image
	originalImg, _ = png.Decode(imgSource)
	// 字体默认大小为15像素
	theFontSize := 15.0
	if t.Size > 0 {
		theFontSize = t.Size
	}
	// 创建一个新的空白图片，与原始图片大小相同
	// 获取原始图片的宽度和高度
	origWidth := originalImg.Bounds().Dx()
	origHeight := originalImg.Bounds().Dy()

	// 创建一个新的空白图片，宽度和高度比原始图片各增加一些像素
	newWidth := origWidth + t.TextXOffset
	newHeight := origHeight + int(t.Size) + t.TextYOffset + 10
	dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	draw.Draw(dst, dst.Bounds(), originalImg, image.Point{}, draw.Src)
	// 在指定位置添加白色文字
	fontColor := image.White
	if len(t.Color) > 0 {
		// 定义颜色值（十六进制）
		fontColor = image.NewUniform(hexToColorRGBA(t.Color))
	}

	// 加载字体文件
	fontBytes, err := os.ReadFile("./assets/RobotoFlex.ttf")
	if err != nil {
		log.Fatalln("AddTextToLogo fail,ttf fail", err)
		return ""
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Fatalln("AddTextToLogo fail,ParseFont fail", err)
		return ""
	}

	f := freetype.NewContext()

	dstBounds := dst.Bounds()
	f.SetDPI(72)
	f.SetFont(font)            // 加载字体
	f.SetFontSize(theFontSize) // 设置字体尺寸
	f.SetClip(dstBounds)
	f.SetDst(dst)
	f.SetSrc(fontColor) // 设置字体颜色

	// 位置信息
	pt := freetype.Pt(0, newHeight-t.TextYOffset)
	_, err = f.DrawString(t.Text, pt)
	if err != nil {
		log.Fatalln("AddTextToLogo fail,DrawString fail", err)
		return ""
	}

	// 创建一个新的文件来保存生成的图片
	fileNamePath := "/tmp/logo-new.png"
	output, err := os.Create(fileNamePath)
	if err != nil {
		log.Fatalln("AddTextToLogo fail,os.Create fail", err)
		return ""
	}
	defer output.Close()

	// 将生成的图片编码为PNG格式并写入文件
	if err = png.Encode(output, dst); err != nil {
		log.Fatalln("AddTextToLogo fail,png.Encode fail", err)
		return ""
	}
	return fileNamePath
}

// 将十六进制颜色值转换为十进制值
func hexToColorRGBA(hexColor string) color.Color {
	// 将十六进制颜色值转换为十进制值
	r, err := strconv.ParseInt(hexColor[1:3], 16, 64)
	if err != nil {
		return color.RGBA{R: 255, G: 255, B: 255, A: 255}
	}
	g, err := strconv.ParseInt(hexColor[3:5], 16, 64)
	if err != nil {
		return color.RGBA{R: 255, G: 255, B: 255, A: 255}
	}
	b, err := strconv.ParseInt(hexColor[5:7], 16, 64)
	if err != nil {
		return color.RGBA{R: 255, G: 255, B: 255, A: 255}
	}
	a := 255 // Alpha通道的值为255（不透明）
	// 创建一个RGBA颜色值，并将其存储为十六进制字符串以进行打印
	return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
}
