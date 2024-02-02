// Package watermark
package watermark

import (
	imgtype "github.com/shamsher31/goimgtype"
	"log"
	"os"
	"path"
)

func AddLogoToImage(sourceImage string, event RequestData) (waterPath string) {
	username := event.Name
	//获取os环境变量中logo在图片中的X,Y坐标偏移
	offsetX, offsetY := 10, 10

	imgSource, err := os.Open(sourceImage)
	if err != nil {
		log.Println(err)
		return ""
	}
	defer imgSource.Close()

	if event.LogoXOffset > 0 {
		offsetX = event.LogoXOffset
	}
	if event.LogoYOffset > 0 {
		offsetY = event.LogoYOffset
	}
	logoUrl := "./assets/image-logo.png"
	usernameLen := len(username)
	if usernameLen > 0 {
		//logo下面有文字需要生成新的logo图片
		defaultNameYOffset := 6
		if event.NameYOffset > 0 {
			defaultNameYOffset = event.NameYOffset
		}
		defaultNameXOffset := 0
		if event.NameXOffset > 0 {
			defaultNameXOffset = event.NameXOffset
		}
		t := TextInfo{
			Text:        username,
			Size:        15,
			TextYOffset: defaultNameYOffset,
			TextXOffset: usernameLen*6 + defaultNameXOffset,
		}
		logoImgSource, err2 := os.Open(logoUrl)
		if err2 != nil {
			log.Fatalln("assets/image-logo.png open fail", err2)
		}
		defer logoImgSource.Close()
		//获取添加文字后的url路径
		logoUrl = t.AddTextToLogo(logoImgSource)
	}

	//获取图片原来的名称和后缀
	imageBaseName := path.Base(sourceImage)

	datatype, err2 := imgtype.Get(sourceImage)
	var imgType = ""
	if err2 != nil {
		imgType = ""
	} else {
		// 根据文件类型执行响应的操作
		switch datatype {
		case `image/jpeg`:
			imgType = "jpeg"
		case `image/png`:
			imgType = "png"
		case `image/gif`:
			imgType = "gif"
		case `image/webp`:
			imgType = "webp"
		}
	}
	if imgType == "" {
		log.Println("暂不支持的图片类型:", imgType)
		return ""
	}
	waterError := err
	newImagePath := ""
	myImageStruct := MyNewImage{
		OffsetX:       offsetX,
		OffsetY:       offsetY,
		ImgType:       imgType,
		ImgSource:     imgSource,
		ImageBaseName: imageBaseName,
		LogoUrl:       logoUrl,
		LogoLocation:  event.LogoLocation,
	}
	log.Println(imgType)
	if imgType == "gif" {
		//gif加水印
		newImagePath, waterError = myImageStruct.GifWaterMark()
	} else if imgType == "webp" {
		//log.Println("暂不支持的图片类型:", imgType)
		newImagePath, waterError = myImageStruct.WebpWatermark()
	} else {
		//png,jpg加水印
		newImagePath, waterError = myImageStruct.PngJpgWaterMark()
	}
	if waterError != nil {
		log.Println("水印添加失败:", waterError)
		return ""
	}
	return newImagePath
}
