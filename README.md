### Description
```
此项目是用Goland语言为aws lambda写的一个处理jpg,jpeg,gif,png图片添加
logo图片水印的函数  
This project is written for aws lambda in Goland to process jpg, jpeg, gif and png image addition
Functions of logo image watermark
``` 
### Path
```
├── demo                       // demo
├── watermark                  // 各类图片处理文件夹 （Various image processing folders）
│   ├── gif.go                 // gif添加logo实现 （Add logo implementation to gif）
│   ├── jpgpng.go              // jpg、png添加logo实现 （Add logo implementation to jpg、png）
│   └── logo.go                // 根据图片大小生成不同尺寸的logo （Generate logos of different sizes according to the image size）
│   └── tool.go                // 工具
│   └── webp.go                // webp 添加logo实现 （Webp adds logo implementation）
├── main.go 
``` 

### Build
```
GOOS=linux CGO_ENABLED=0 go build main.go 
```

### Zip for lambda
```
zip -r  image-watermark.zip main logo.png
```