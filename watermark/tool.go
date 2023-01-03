package watermark

import (
	"math/rand"
	"time"
)

// GetRand  获取[0,max)随机数
func GetRand(max int) (i int) {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max)
}

//GetRandomString 生成图片名字
func GetRandomString(lenght int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	bytesLen := len(bytes)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < lenght; i++ {
		result = append(result, bytes[r.Intn(bytesLen)])
	}
	return string(result)
}

func NewImageName(oldImageName string) string {
	return savePath + "/" + oldImageName
}
