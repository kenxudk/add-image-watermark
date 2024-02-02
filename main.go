package main

import (
	"context"
	"errors"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"image-watermark/watermark"
	"io"
	"log"
	"os"
	"path"
	"strings"
)

// s3客服端
var s3Client *s3.Client
var (
	s3Bucket = ""
	awsAk    = ""
	awsSk    = ""
	//cdnDomain = "" //域名
	awsRegion = "" //region
)

// 图片下载和生成的存放目录
var savePath = "/tmp"

// 水印在s3里面存放的前缀目录，比如以前的key=feed/sss.jpg ; 需要拼接为 key=watermark/feed/sss.jpg
var watermarkPrefixPath = "watermark/"

// ResponseData 返回数据
type ResponseData struct {
	Body BodyData `json:"body"`
}
type BodyData struct {
	Data string `json:"data"`
}

// HandleLambdaEvent lambda
func HandleLambdaEvent(event watermark.RequestData) (ResponseData, error) {
	key := event.Key
	username := event.Name
	log.Println("start key=" + key + ",name=" + username)
	//通过key,下载图片,返回下载存储的位置
	downloadImagePath := downloadImageFromS3(key)
	if downloadImagePath == "" {
		return ResponseData{Body: BodyData{Data: ""}}, errors.New("下载失败")
	}
	log.Println("download path=" + downloadImagePath)
	//添加水印，返回添加后的水印地址
	waterPath := watermark.AddLogoToImage(downloadImagePath, event)
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

// 初始化
// os.Getenv("Env") 获取各个环境变量
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

// 获取s3客户端
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

// 从s3下载图片
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

// 上传添加logo的图片到s3
func uploadToS3(newImagePath string, oldKey string) (url string, err error) {
	//防止webp加水印后变jpg,后缀还是webp
	newImagePathExt := path.Ext(newImagePath)
	oldImagePathExt := path.Ext(oldKey)
	if oldImagePathExt != newImagePathExt && len(newImagePathExt) > 0 && len(oldImagePathExt) > 0 {
		oldKey = strings.Replace(oldKey, oldImagePathExt, newImagePathExt, 1)
	}
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
