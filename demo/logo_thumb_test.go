// Package main
package main_test

import (
	"github.com/nfnt/resize"
	"image/jpeg"
	"image/png"
	"log"
	"os"
)

func logo() {
	watermark, err := os.Open("/tmp/logo-new.png")
	if err != nil {
		log.Fatalln(err)
	}
	defer watermark.Close()
	imgWatermark, err := png.Decode(watermark)
	scaleH := int(float64(imgWatermark.Bounds().Dy()) * 0.75)
	log.Println(imgWatermark.Bounds().Dy(), scaleH)
	thumbImage := resize.Thumbnail(uint(imgWatermark.Bounds().Dx()), uint(scaleH), imgWatermark, resize.Bicubic)
	imgW, _ := os.Create("/tmp/logo_thumb2.png")
	_ = jpeg.Encode(imgW, thumbImage, &jpeg.Options{Quality: 100})
	log.Println("success, path=/tmp/logo_thumb2.png")
}
