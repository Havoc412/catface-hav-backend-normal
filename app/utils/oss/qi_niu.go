package oss

import (
	_ "catface/bootstrap"
	"context"
	"fmt"
	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/http_client"
	"github.com/qiniu/go-sdk/v7/storagev2/objects"
	"github.com/qiniu/go-sdk/v7/storagev2/uploader"
	"io"
	"log"
	"os"
	"time"
)

// MultiUploadToQiNiu 批量上传文件到七牛云
// images: 文件列表
// 返回值: 文件url列表, 错误信息
func MultiUploadToQiNiu(images []string) ([]string, error) {

	// 并发上传图片
	errCh := make(chan error, len(images))

	// 这里初始化而不扩容切片，因为在循环中会对切片进行赋值，扩容会导致切片地址变化
	resImageURLs := make([]string, len(images)) // 保证图片的顺序

	for idx, image := range images {
		// 闭包捕获循环变量
		go func(index int, img string) {
			// 打开文件
			file, err := os.OpenFile(img, os.O_RDONLY, 0666)
			if err != nil {
				log.Println(err)
				errCh <- err
				return
			}

			defer func(file *os.File) {
				err := file.Close()
				if err != nil {
					log.Println(err)
				}
			}(file)

			imageStr, err := UploadToQiNiu(file)
			if err != nil {
				log.Println(err)
				errCh <- err
				return
			}
			resImageURLs[index] = imageStr
			errCh <- nil
		}(idx, image)
	}

	// 等待所有图片上传完成并检查错误
	for range images {
		if err := <-errCh; err != nil {
			return nil, err
		}
	}

	return resImageURLs, nil
}

// UploadToQiNiu 上传文件到七牛云
// reader: 文件流
// 返回值: 文件url, 错误信息
func UploadToQiNiu(reader io.Reader) (string, error) {
	accessKey := qiNiuAccessKey
	secretKey := qiNiuSecretKey
	bucket := qiNiuBucket

	mac := credentials.NewCredentials(accessKey, secretKey)

	key := fmt.Sprintf("img_%d", time.Now().Unix())

	uploadManager := uploader.NewUploadManager(&uploader.UploadManagerOptions{
		Options: http_client.Options{
			Credentials: mac,
		},
	})
	err := uploadManager.UploadReader(context.Background(), reader, &uploader.ObjectOptions{
		BucketName: bucket,
		ObjectName: &key,
		FileName:   key,
	}, nil)
	if err != nil {
		log.Fatal("上传失败，可能是key和secret配错了：", err)
		return "", err
	}

	return qiNiuDomain + "/" + key, err

}

func DeleteFromQiNiu(key string) error {
	accessKey := qiNiuAccessKey
	secretKey := qiNiuSecretKey
	bucketName := qiNiuBucket

	mac := credentials.NewCredentials(accessKey, secretKey)
	objectsManager := objects.NewObjectsManager(&objects.ObjectsManagerOptions{
		Options: http_client.Options{Credentials: mac},
	})

	bucket := objectsManager.Bucket(bucketName)

	err := bucket.Object(key).Delete().Call(context.Background())
	if err != nil {
		return err
	}

	return nil
}
