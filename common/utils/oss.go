package utils

import (
	"bytes"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
)

type OssClient interface {
	UploadFile(fileName, filePath string, data []byte) (string, error)
	GetOssFileFullAccessPath(filePath string) string
	DownLoadFile(fileUploadPath string) ([]byte, error)
	DeleteFile(fileFullPath string) error
}

type AliOssClient struct {
	AccessKeyId     string
	AccessKeySecret string
	BucketName      string
	Endpoint        string
	Client          *oss.Client
}

func NewAliOssClient(accessKeyId, accessKeySecret, endpoint, bucketName string) *AliOssClient {
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		panic(fmt.Sprintf("初始化Oss客户端失败: %v", err))
	}
	return &AliOssClient{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
		Endpoint:        endpoint,
		Client:          client,
		BucketName:      bucketName,
	}
}

func (o *AliOssClient) UploadFile(fileName, filePath string, data []byte) (string, error) {
	if len(fileName) <= 0 || len(filePath) <= 0 || data == nil {
		return "", errors.New("上传必须指定文件名称、文件路径、文件内容")
	}
	bucket, err := o.Client.Bucket(o.BucketName)
	if err != nil {
		return "", errors.Wrap(err, "初始化oss链接失败")
	}
	reader := bytes.NewReader(data)
	fileUploadPath := filePath + "/" + fileName
	err = bucket.PutObject(fileUploadPath, reader)
	if err != nil {
		return "", errors.Wrap(err, "文件上传失败")
	}
	return fileUploadPath, nil
}

func (o *AliOssClient) GetOssFileFullAccessPath(filePath string) string {
	return "https://" + o.BucketName + "." + o.Endpoint + "/" + filePath
}

func (o *AliOssClient) DownLoadFile(fileUploadPath string) ([]byte, error) {
	bucket, err := o.Client.Bucket(o.BucketName)
	if err != nil {
		return nil, errors.Wrap(err, "初始化oss链接失败")
	}
	object, err := bucket.GetObject(fileUploadPath)
	if err != nil {
		logx.Errorf("下载OSS文件失败：%v", err)
		return nil, errors.Wrap(err, "下载文件失败")
	}
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, object); err != nil {
		logx.Errorf("下载OSS文件失败：%v", err)
		return nil, err
	}
	defer func(object io.ReadCloser) {
		err := object.Close()
		if err != nil {

		}
	}(object)
	return buf.Bytes(), nil
}

func (o *AliOssClient) DeleteFile(fileFullPath string) error {
	if len(fileFullPath) <= 0 {
		return errors.New("删除文件必须指定文件完整路径")
	}
	bucket, err := o.Client.Bucket(o.BucketName)
	if err != nil {
		return errors.Wrap(err, "删除文件失败")
	}
	err = bucket.DeleteObject(fileFullPath)
	if err != nil {
		return errors.Wrap(err, "删除文件失败")
	}
	return nil
}
