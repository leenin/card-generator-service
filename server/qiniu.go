package server

import (
	"bytes"
	"context"
	"os"

	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
)

type Qiniu struct {
	Domain    string `json:"domain"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Bucket    string `json:"bucket"`
}

func (qn *Qiniu) init() {
	qn.Domain = os.Getenv("QINIU_DOMAIN")
	qn.AccessKey = os.Getenv("QINIU_ACCESS_KEY")
	qn.SecretKey = os.Getenv("QINIU_SECRET_KEY")
	qn.Bucket = os.Getenv("QINIU_BUCKET")
}

func (qn *Qiniu) getFileInfoByKey(key string) (fileInfo storage.FileInfo) {
	mac := qbox.NewMac(qn.AccessKey, qn.SecretKey)
	cfg := storage.Config{
		// 是否使用https域名进行资源管理
		UseHTTPS: false,
	}
	// 指定空间所在的区域，如果不指定将自动探测
	// 如果没有特殊需求，默认不需要指定
	//cfg.Zone=&storage.ZoneHuabei
	bucketManager := storage.NewBucketManager(mac, &cfg)

	fileInfo, _ = bucketManager.Stat(qn.Bucket, key)
	return
}

func (qn *Qiniu) uploadFile(data []byte, key string) (ret storage.PutRet, err error) {
	putPolicy := storage.PutPolicy{
		Scope: qn.Bucket,
	}
	mac := qbox.NewMac(qn.AccessKey, qn.SecretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Zone = &storage.ZoneHuanan
	// 是否使用https域名
	cfg.UseHTTPS = false
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	formUploader := storage.NewFormUploader(&cfg)
	putExtra := storage.PutExtra{}
	dataLen := int64(len(data))
	err = formUploader.Put(context.Background(), &ret, upToken, key, bytes.NewReader(data), dataLen, &putExtra)
	return
}
