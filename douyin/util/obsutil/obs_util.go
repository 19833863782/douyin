// Package obsutil
// @Author shaofan
// @Date 2022/5/16
package obsutil

import (
	"douyin/config"
	"fmt"
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	uuid "github.com/satori/go.uuid"
	"log"
	"strings"
)

// Upload 上传文件
// filePath 文件路径
// bucketName 桶名
// @return 文件云端url
func Upload(filePath string, bucketName string) (string, error) {
	client, err := obs.New(
		config.Config.Obs.AccessKey,
		config.Config.Obs.SecretKey,
		config.Config.Obs.EndPoint,
	)
	if err != nil {
		return "", err
	}
	defer client.Close()
	exist, err := bucketIsExist(bucketName, client)
	if err != nil {
		return "", err
	}
	if !exist {
		err = createBucket(bucketName, client)
		if err != nil {
			return "", err
		}
	}
	var objectName = ParseFileName(filePath)
	err = doUpload(filePath, bucketName, client, objectName)
	if err != nil {
		return "", err
	}
	objectName = strings.Replace(config.Config.Obs.EndPoint,
		"https://obs",
		"https://"+bucketName+".obs", 1) + "/" + objectName
	return objectName, nil
}

// 判断桶是否存在
func bucketIsExist(bucketName string, client *obs.ObsClient) (bool, error) {
	_, err := client.HeadBucket(bucketName)
	if err == nil {
		return true, err
	} else {
		if obsError, ok := err.(obs.ObsError); ok {
			if obsError.StatusCode == 404 {
				return false, nil
			} else {
				log.Printf("StatusCode:%d\n", obsError.StatusCode)
				return false, nil
			}
		} else {
			log.Println(err)
		}
	}
	return false, err
}

// 创建桶
func createBucket(bucketName string, client *obs.ObsClient) error {
	input := &obs.CreateBucketInput{}
	input.Bucket = bucketName
	input.ACL = obs.AclPublicRead
	input.StorageClass = obs.StorageClassStandard
	input.Location = config.Config.Obs.Location
	output, err := client.CreateBucket(input)
	if err == nil {
		log.Printf("RequestId:%s\n", output.RequestId)
	} else {
		if obsError, ok := err.(obs.ObsError); ok {
			log.Println(obsError.Code)
			log.Println(obsError.Message)
		} else {
			log.Println(err)
		}
	}
	return err
}

// 分段上传文件
// filePath 文件本地路径
// bucketName 桶名
// client obs客户端
// objectName 对象名
func doUpload(filePath string, bucketName string, client *obs.ObsClient, objectName string) error {
	input := &obs.UploadFileInput{}
	input.Bucket = bucketName
	input.Key = objectName
	input.UploadFile = filePath      // localfile为待上传的本地文件路径，需要指定到具体的文件名
	input.EnableCheckpoint = true    // 开启断点续传模式
	input.PartSize = 9 * 1024 * 1024 // 指定分段大小为9MB
	input.TaskNum = 5                // 指定分段上传时的最大并发数
	output, err := client.UploadFile(input)
	if err == nil {
		fmt.Printf("RequestId:%s\n", output.RequestId)
		fmt.Printf("ETag:%s\n", output.ETag)
	} else if obsError, ok := err.(obs.ObsError); ok {
		fmt.Printf("Code:%s\n", obsError.Code)
		fmt.Printf("Message:%s\n", obsError.Message)
	}
	return err
}

// 获取uuid
func getUUID() string {
	UUID := uuid.NewV4()
	return UUID.String()
}

// ParseFileName 解析文件名，创建一个唯一文件名
// originName 源文件名
func ParseFileName(originName string) string {
	return strings.ReplaceAll(getUUID(), "-", "") + strings.ToLower(originName[strings.LastIndex(originName, "."):])
}

// IsVideo 判断是否是时评
func IsVideo(fileName string) bool {
	var suffix = strings.ToLower(fileName[strings.LastIndex(fileName, "."):])
	var videoSuffix = []string{".avi", ".mp4", ".rmvp"}
	for _, v := range videoSuffix {
		if v == suffix {
			return true
		}
	}
	return false
}
