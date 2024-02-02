// Package demo
package main

import (
	"image-watermark/watermark"
	"log"
)

func main() {
	downloadImagePath := "./demo/test.webp"
	event := watermark.RequestData{
		Name: "@kenxu",
	}
	waterPath := watermark.AddLogoToImage(downloadImagePath, event)
	log.Println("watermark success, path=" + waterPath)
}
