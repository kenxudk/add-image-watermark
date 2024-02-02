package watermark

import (
	"errors"
	"golang.org/x/image/webp"
	"image"
	"image/draw"
	"image/jpeg"
	"os"
	"path"
	"strconv"
	"strings"
)

func (myNewImage MyNewImage) WebpWatermark() (newImagePath string, err error) {
	imgSource := myNewImage.ImgSource
	offsetX := myNewImage.OffsetX
	offsetY := myNewImage.OffsetY
	imageBaseName := myNewImage.ImageBaseName
	logoUrl := myNewImage.LogoUrl
	imgBInfo, er := webp.Decode(imgSource)
	if er != nil {
		return "", er
	}
	//读取水印图片
	imgWatermark, isScale, err := GetLogoImage(imgBInfo.Bounds().Dy(), logoUrl)
	if err != nil {
		return "", err
	}

	//获取logo放的位置
	randNumber := myNewImage.LogoLocation
	if myNewImage.LogoLocation <= 0 {
		randNumber = GetRand(4)
	}
	if isScale {
		offsetX = offsetX / 8
		offsetY = offsetY / 8
	}
	logoX, logoY := offsetX, offsetY

	switch randNumber {
	case 1:
		logoX = imgBInfo.Bounds().Dx() - imgWatermark.Bounds().Dx() - offsetX
	case 2:
		logoY = imgBInfo.Bounds().Dy() - imgWatermark.Bounds().Dy() - offsetY
	case 3:
		logoX = imgBInfo.Bounds().Dx() - imgWatermark.Bounds().Dx() - offsetX
		logoY = imgBInfo.Bounds().Dy() - imgWatermark.Bounds().Dy() - offsetY
	}

	//如果X，Y<0,就不加
	if logoX <= 0 || logoY <= 0 {
		return "", errors.New("原图宽或者高小于或者等于了最小偏移量,即W=" + strconv.Itoa(offsetX) + ", H=" + strconv.Itoa(offsetY))
	}

	offset := image.Pt(logoX, logoY)
	b := imgBInfo.Bounds()
	m := image.NewNRGBA(b) //按原图生成新图

	//新图写入原图和背景图
	draw.Draw(m, b, imgBInfo, image.Point{}, draw.Src)
	draw.Draw(m, imgWatermark.Bounds().Add(offset), imgWatermark, image.Point{}, draw.Over)

	//输出图像
	//webp存储为jpg
	imageBaseName = strings.TrimSuffix(imageBaseName, path.Ext(imageBaseName)) + ".jpg"
	imageNewPath := NewImageName(imageBaseName)
	imgW, _ := os.Create(imageNewPath)
	cErr := jpeg.Encode(imgW, m, &jpeg.Options{Quality: 100})
	if cErr != nil {
		return "", cErr
	}
	return imageNewPath, nil
}
