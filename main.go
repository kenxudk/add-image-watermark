package main

import (
	"context"
	"errors"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	imgtype "github.com/shamsher31/goimgtype"
	"image-watermark/watermark"
	"io"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
)

//s3客服端
var s3Client *s3.Client
var (
	s3Bucket  = ""
	awsAk     = ""
	awsSk     = ""
	cdnDomain = "" //域名
	awsRegion = "" //region
)

//图片下载和生成的存放目录
var savePath = "/tmp"

//水印在s3里面存放的前缀目录，比如以前的key=feed/sss.jpg ; 需要拼接为 key=watermark/feed/sss.jpg
var watermarkPrefixPath = "watermark/"

// RequestData 请求的json数据
type RequestData struct {
	Channel string `json:"channel"`
	Name    string `json:"name"`
	Key     string `json:"key"`
}

// ResponseData 返回数据
type ResponseData struct {
	Body BodyData `json:"body"`
}
type BodyData struct {
	Data string `json:"data"`
}

// HandleLambdaEvent lambda
func HandleLambdaEvent(event RequestData) (ResponseData, error) {
	key := event.Key
	log.Println("start key=" + key)
	//如果是webp的图片，改为jpg的后缀，因为app上传有jpg的图片生成
	key = strings.Replace(key, ".webp", ".jpg", 1)
	//通过key,下载图片,返回下载存储的位置
	downloadImagePath := downloadImageFromS3(key)
	if downloadImagePath == "" {
		return ResponseData{Body: BodyData{Data: ""}}, errors.New("下载失败")
	}
	log.Println("download path=" + downloadImagePath)
	//添加水印，返回添加后的水印地址
	waterPath := addLog(downloadImagePath)
	if waterPath == "" {
		return ResponseData{Body: BodyData{Data: ""}}, errors.New("水印失败")
	}
	log.Println("watermark success, path=" + waterPath)
	//添加logo后的水印地址,上传到s3
	url, err := uploadToS3(waterPath, key)
	if err != nil {
		return ResponseData{Body: BodyData{Data: ""}}, err
	}
	log.Println("upload s3 success, key=" + url)
	return ResponseData{Body: BodyData{Data: url}}, nil
}

//初始化
//os.Getenv("Env") 获取各个环境变量
func init() {
	//获取s3桶名
	s3Bucket = os.Getenv("Bucket")
	if s3Bucket == "" {
		log.Println("你还没有配置S3的桶名")
		return
	}
	awsAk = os.Getenv("AwsAccessKey")
	if awsAk == "" {
		log.Println("你还没有aws ak")
		return
	}
	awsSk = os.Getenv("AwsSecretKey")
	if awsSk == "" {
		log.Println("你还没有aws sk")
		return
	}

	awsRegion = os.Getenv("AwsRegion")
	if awsRegion == "" {
		log.Println("你还没有配置aws region")
		return
	}

	s3Client = getAWSS3Client()
}

//获取s3客户端
func getAWSS3Client() *s3.Client {
	options := s3.Options{
		Region:      awsRegion,
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(awsAk, awsSk, "")),
	}
	client := s3.New(options, func(o *s3.Options) {
		o.Region = awsRegion
	})
	return client
}

//从s3下载图片
func downloadImageFromS3(key string) (originImg string) {
	// Get the first page of results for ListObjectsV2 for a bucket
	out, err := s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		log.Fatal("下载失败: ", err)
		return ""
	}
	defer out.Body.Close()
	downloadImagePath := savePath + "/" + path.Base(key)
	downAndSave, err3 := os.OpenFile(downloadImagePath, os.O_CREATE|os.O_RDWR, 0666)
	if err3 != nil {
		log.Fatal("downAndSave error : ", err3)
		return ""
	}
	_, err2 := io.Copy(downAndSave, out.Body)
	if err2 != nil {
		log.Fatal("download Copy error: ", err2)
		return ""
	}
	return downloadImagePath
}

//上传添加logo的图片到s3
func uploadToS3(newImagePath string, oldKey string) (url string, err error) {
	newKey := watermarkPrefixPath + oldKey
	fb, err := os.Open(newImagePath)
	if err != nil {
		return "", err
	}
	defer fb.Close()
	_, err2 := s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(newKey),
		Body:   fb,
	})
	if err2 != nil {
		return "", err2
	}
	url = newKey
	return url, nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}

func addLog(sourceImage string) (waterPath string) {
	//获取os环境变量中logo在图片中的X,Y坐标偏移
	osOffsetX, osOffsetY := "80", "80"
	if os.Getenv("offsetX") != "" {
		osOffsetX = os.Getenv("offsetX")
	}
	offsetX, _ := strconv.Atoi(osOffsetX)
	if os.Getenv("offsetY") != "" {
		osOffsetY = os.Getenv("offsetY")
	}
	offsetY, _ := strconv.Atoi(osOffsetY)

	//var sourceImage string = "/Users/mac/Desktop/1669018251853198000-2522543887406335640.jpeg"
	//var sourceImage string = "/Users/mac/Desktop/WechatIMG723.jpeg"
	//var sourceImage string = "/Users/mac/Desktop/dlq1f11q42430ou5tjjagrttvp-16721597133461212263597.gif"
	//var sourceImage string = "/Users/mac/Desktop/GzauTpqIUQ8KvLRlWBMlMrrVoWWeOhGG.gif"
	//var sourceImage string = "/Users/mac/Desktop/default_audio_live.webp"
	imgSource, err := os.Open(sourceImage)
	if err != nil {
		log.Println(err)
		return ""
	}
	defer imgSource.Close()

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
	if imgType == "gif" {
		//gif加水印
		newImagePath, waterError = watermark.GifWaterMark(imgSource, imageBaseName)
	} else if imgType == "webp" {
		//log.Println("暂不支持的图片类型:", imgType)
		newImagePath, waterError = watermark.WebpWatermark(offsetX, offsetY, imgType, imgSource, imageBaseName)
	} else {
		//png,jpg加水印
		newImagePath, waterError = watermark.PngJpgWaterMark(offsetX, offsetY, imgType, imgSource, imageBaseName)
	}
	if waterError != nil {
		log.Println("水印添加失败:", waterError)
		return ""
	}
	return newImagePath
}
