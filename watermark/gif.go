package watermark

import (
	"errors"
	"image"
	"image/draw"
	"image/gif"
	"os"
)

func GifWaterMark(imgSource *os.File, imageBaseName string) (newImagePath string, err error) {
	gifImgs, _ := gif.DecodeAll(imgSource)
	var newGifImgs = make([]*image.Paletted, 0)
	x0, y0, old := 0, 0, 0
	//读取水印图片
	imgWatermark, _, err := GetLogoImage(gifImgs.Image[0].Bounds().Dx())
	//imgWatermark, isScale, err := GetLogoImage(gifImgs.Image[0].Bounds().Dx())
	if err != nil {
		return "", err
	}
	logoX, logoY := gifImgs.Image[0].Bounds().Dx()-imgWatermark.Bounds().Dx(), gifImgs.Image[0].Bounds().Dy()-imgWatermark.Bounds().Dy()

	offset := image.Pt(logoX, logoY)
	//fmt.Println(logoX, logoY, imgWatermark.Bounds().Dx(), imgWatermark.Bounds().Dy(), isScale)
	for k, gifImg := range gifImgs.Image {
		img := image.NewNRGBA(gifImg.Bounds())
		if k == 0 {
			x0 = img.Bounds().Dx()
			y0 = img.Bounds().Dy()
		}
		//fmt.Printf("%v, %v\n", img.Bounds().Dx(), img.Bounds().Dy())
		if k == 0 && gifImgs.Image[k+1].Bounds().Dx() > x0 && gifImgs.Image[k+1].Bounds().Dy() > y0 {
			old = 1
			break
		}
		if x0 == img.Bounds().Dx() && y0 == img.Bounds().Dy() {

			p1 := image.NewPaletted(gifImg.Bounds(), gifImg.Palette)
			//把logo添加到新的图片调色板上
			draw.Draw(p1, gifImg.Bounds(), gifImg, image.Point{}, draw.Src)

			draw.Draw(p1, imgWatermark.Bounds().Add(offset), imgWatermark, image.Point{}, draw.Over)
			//把添加过文字的新调色板放入调色板slice
			newGifImgs = append(newGifImgs, p1)
		} else {
			newGifImgs = append(newGifImgs, gifImg)
		}

	}
	if old == 1 {
		return "", errors.New("gif: image block is out of bounds")
	} else {
		//保存到新文件中
		imageNewPath := NewImageName(imageBaseName)
		//imageNewPath := "/Users/mac/Desktop/water-" + imageBaseName
		//fmt.Println(imageNewPath)
		newFile, err2 := os.Create(imageNewPath)
		if err2 != nil {
			return "", err2
		}
		defer newFile.Close()
		//fmt.Println(gifImgs.Delay, gifImgs.LoopCount)
		g1 := &gif.GIF{
			Image:     newGifImgs,
			Delay:     gifImgs.Delay,
			LoopCount: gifImgs.LoopCount,
		}
		err = gif.EncodeAll(newFile, g1)
		if err != nil {
			return "", err
		}
		return imageNewPath, nil
	}
}
