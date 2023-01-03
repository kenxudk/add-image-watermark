package watermark

import (
	"github.com/nfnt/resize"
	"image"
	"image/png"
	"os"
)

//logo图片名字
var logoImageName string = "logo.png"

//图片下载和生成的存放目录
var savePath string = "/tmp"

func GetLogoImage(imageW int) (image.Image, bool, error) {
	watermark, err := os.Open(logoImageName)
	if err != nil {
		return nil, false, err
	}
	defer watermark.Close()
	imgWatermark, err := png.Decode(watermark)
	if err != nil {
		return nil, false, err
	}
	//按照加logo的图片和logo的图片8：1来缩小
	logoW := imgWatermark.Bounds().Dx()
	//8倍缩小后的logo宽和高
	scaleW := imageW / 8
	if scaleW < logoW {
		minScaleW := 40
		if scaleW < minScaleW {
			scaleW = minScaleW
		}
		return resize.Thumbnail(uint(scaleW), uint(scaleW), imgWatermark, resize.Lanczos3), true, nil
	} else {
		return imgWatermark, false, nil
	}
}
