### Description
```
此项目是用Goland语言为aws lambda写的一个处理jpg,jpeg,gif,png图片添加
logo图片水印的函数  
This project is written for aws lambda in Goland to process jpg, jpeg, gif and png image addition
Functions of logo image watermark
``` 
### Path
```
├── demo                       // test go file
├── watermark                  // 各类图片处理文件夹 （Various image processing folders）
│   ├── gif.go                 // gif添加logo实现 （Add logo implementation to gif）
│   ├── jpgpng.go              // jpg、png添加logo实现 （Add logo implementation to jpg、png）
│   └── logo.go                // 根据图片大小生成不同尺寸的logo （Generate logos of different sizes according to the image size）
│   └── tool.go                // 工具
│   └── request.go             // 请求参数结构体
│   └── image.go               // 对不同类型图片加水印的实现
│   └── webp.go                // webp 添加logo实现 （Webp adds logo implementation）
├── main.go 
``` 

### Build
```
GOOS=linux CGO_ENABLED=0 go build main.go 
```

### Zip for lambda
```
zip -r  image-watermark.zip main ./assets
```
### The test of Event json for lambda
```
{
  "channel": "",
  "name": "@test",
  "key": "image/258_1706586088.png",
  "file_h":720,
  "file_w":1280,
  "fontsize":15,
  "fontcolor":"#FFFFFF",
  "name_y_offset":6,//logo下面文字相对logo的y偏移
  "name_x_offset":0,//logo下面文字相对logo的x偏移
  "logo_location":0,//logo位置 0 随机，1左上角，2右上角，3左下角，4右下角
  "logo_y_offset":10,//logo相对图片的y偏移
  "logo_x_offset":10,//logo相对图片的x偏移
}
```