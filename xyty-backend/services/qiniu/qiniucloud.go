package qiniu

import (
	"context"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"ini/services/parseyaml"
	"mime/multipart"
)

// 上传图片到七牛云，然后返回状态和图片的url
func UploadToQiNiu(file *multipart.FileHeader, folder string) (int, string) {
	v := parseyaml.GetYaml()
	var AccessKey = v.GetString("qiniu.AccessKey") // 秘钥对
	var SerectKey = v.GetString("qiniu.SerectKey")
	var Bucket = v.GetString("qiniu.Bucket") // 空间名称
	var ImgUrl = v.GetString("qiniu.ImgUrl") // 自定义域名或测试域名

	src, err := file.Open()
	if err != nil {
		return 10011, err.Error()
	}
	defer src.Close()

	putPlicy := storage.PutPolicy{
		Scope: Bucket,
	}
	mac := qbox.NewMac(AccessKey, SerectKey)

	// 获取上传凭证
	upToken := putPlicy.UploadToken(mac)

	// 配置参数
	cfg := storage.Config{
		Zone:          &storage.ZoneHuanan, // 华南区
		UseCdnDomains: false,
		UseHTTPS:      false, // 非https
	}
	formUploader := storage.NewFormUploader(&cfg)

	ret := storage.PutRet{}        // 上传后返回的结果
	putExtra := storage.PutExtra{} // 额外参数

	// 上传 自定义key，可以指定上传目录及文件名和后缀，
	key := "drifting/" + folder + file.Filename // 上传路径，如果当前目录中已存在相同文件，则返回上传失败错误
	err = formUploader.Put(context.Background(), &ret, upToken, key, src, file.Size, &putExtra)

	// 以默认key方式上传
	// err = formUploader.PutWithoutKey(context.Background(), &ret, upToken, src, fileSize, &putExtra)

	// 自定义key，上传指定路径的文件
	// localFilePath = "./aa.jpg"
	// err = formUploader.PutFile(context.Background(), &ret, upToken, key, localFilePath, &putExtra)

	// 默认key，上传指定路径的文件
	// localFilePath = "./aa.jpg"
	// err = formUploader.PutFile(context.Background(), &ret, upToken, key, localFilePath, &putExtra)

	if err != nil {
		code := 501
		return code, err.Error()
	}

	url := ImgUrl + "/" + ret.Key // 返回上传后的文件访问路径
	return 1, url
}

// ListFilesByPrefix 按前缀查询七牛云文件列表
func ListFilesByPrefix(prefix string) ([]string, int, error) {
	v := parseyaml.GetYaml()
	accessKey := v.GetString("qiniu.AccessKey")
	secretKey := v.GetString("qiniu.SerectKey")
	bucket := v.GetString("qiniu.Bucket")
	imgUrl := v.GetString("qiniu.ImgUrl")

	mac := qbox.NewMac(accessKey, secretKey)
	cfg := storage.Config{UseHTTPS: false}
	bucketManager := storage.NewBucketManager(mac, &cfg)

	limit := 1000
	delimiter := ""
	marker := ""

	list, _, nextMarker, hasNext, err := bucketManager.ListFiles(bucket, prefix, delimiter, marker, limit)
	if err != nil {
		return nil, 0, err
	}

	var fileUrls []string
	for _, item := range list {
		fileUrls = append(fileUrls, imgUrl+"/"+item.Key)
	}

	// 处理分页（如果需要）
	for hasNext {
		list, _, nextMarker, hasNext, err = bucketManager.ListFiles(bucket, prefix, delimiter, nextMarker, limit)
		if err != nil {
			return nil, 0, err
		}
		for _, item := range list {
			fileUrls = append(fileUrls, imgUrl+"/"+item.Key)
		}
	}

	return fileUrls, 1, nil
}
