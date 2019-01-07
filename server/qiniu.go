package server

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
)

type QiniuCfg struct {
	Domain    string `json:"domain"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Bucket    string `json:"bucket"`
}

func getQiniuCfg(qcfg *QiniuCfg) (err error) {
	fp, err := filepath.Abs("config/qiniu.json")
	if err != nil {
		return
	}

	qiniuCfgJSONBytes, err := ioutil.ReadFile(fp)
	if err != nil {
		return
	}

	json.Unmarshal(qiniuCfgJSONBytes, &qcfg)
	return
}

func getFileInfoByKey(qcfg *QiniuCfg, key string) (fileInfo storage.FileInfo, err error) {

	mac := qbox.NewMac(qcfg.AccessKey, qcfg.SecretKey)
	cfg := storage.Config{
		// 是否使用https域名进行资源管理
		UseHTTPS: false,
	}
	// 指定空间所在的区域，如果不指定将自动探测
	// 如果没有特殊需求，默认不需要指定
	//cfg.Zone=&storage.ZoneHuabei
	bucketManager := storage.NewBucketManager(mac, &cfg)

	fileInfo, err = bucketManager.Stat(qcfg.Bucket, key)
	if err != nil {
		err = nil
		return
	}
	return
}

func uploadFile(qcfg *QiniuCfg, data []byte, key string) (ret storage.PutRet, err error) {
	putPolicy := storage.PutPolicy{
		Scope: qcfg.Bucket,
	}
	mac := qbox.NewMac(qcfg.AccessKey, qcfg.SecretKey)
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
