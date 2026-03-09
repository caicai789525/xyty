package qiniu

import (
	"context"
	"fmt"
	"ini/services/parseyaml"
	"mime/multipart"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
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
		// 确保 URL 格式正确，避免双斜杠
		if len(imgUrl) > 0 && imgUrl[len(imgUrl)-1] == '/' {
			fileUrls = append(fileUrls, imgUrl+item.Key)
		} else {
			fileUrls = append(fileUrls, imgUrl+"/"+item.Key)
		}
	}

	// 处理分页（如果需要）
	for hasNext {
		list, _, nextMarker, hasNext, err = bucketManager.ListFiles(bucket, prefix, delimiter, nextMarker, limit)
		if err != nil {
			return nil, 0, err
		}
		for _, item := range list {
			// 确保 URL 格式正确，避免双斜杠
			if len(imgUrl) > 0 && imgUrl[len(imgUrl)-1] == '/' {
				fileUrls = append(fileUrls, imgUrl+item.Key)
			} else {
				fileUrls = append(fileUrls, imgUrl+"/"+item.Key)
			}
		}
	}

	return fileUrls, 1, nil
}

// GetFileURL 获取指定文件的访问外链
// key: 文件的完整存储路径
// isPrivate: 是否为私有空间
// expire: 外链有效期（秒），仅当 isPrivate 为 true 时有效
func GetFileURL(key string, isPrivate bool, expire int64) (string, error) {
	v := parseyaml.GetYaml()
	accessKey := v.GetString("qiniu.AccessKey")
	secretKey := v.GetString("qiniu.SerectKey")
	bucket := v.GetString("qiniu.Bucket")
	imgUrl := v.GetString("qiniu.ImgUrl")

	// 检查配置是否完整
	if accessKey == "" || secretKey == "" || bucket == "" || imgUrl == "" {
		return "", fmt.Errorf("七牛云配置不完整")
	}

	// 检查文件路径是否为空
	if key == "" {
		return "", fmt.Errorf("文件路径不能为空")
	}

	mac := qbox.NewMac(accessKey, secretKey)

	// 对于公开空间，直接拼接 URL
	if !isPrivate {
		// 确保 URL 格式正确，避免双斜杠
		if len(imgUrl) > 0 && imgUrl[len(imgUrl)-1] == '/' {
			return imgUrl + key, nil
		}
		return imgUrl + "/" + key, nil
	}

	// 对于私有空间，生成带签名的临时外链
	if expire <= 0 {
		expire = 3600 // 默认 1 小时
	}

	domain := imgUrl
	privateURL := storage.MakePrivateURL(mac, domain, key, expire)
	return privateURL, nil
}

// GetFileURLs 批量获取文件的访问外链
// keys: 文件的完整存储路径列表
// isPrivate: 是否为私有空间
// expire: 外链有效期（秒），仅当 isPrivate 为 true 时有效
func GetFileURLs(keys []string, isPrivate bool, expire int64) ([]string, error) {
	v := parseyaml.GetYaml()
	accessKey := v.GetString("qiniu.AccessKey")
	secretKey := v.GetString("qiniu.SerectKey")
	bucket := v.GetString("qiniu.Bucket")
	imgUrl := v.GetString("qiniu.ImgUrl")

	// 检查配置是否完整
	if accessKey == "" || secretKey == "" || bucket == "" || imgUrl == "" {
		return nil, fmt.Errorf("七牛云配置不完整")
	}

	// 检查文件路径列表是否为空
	if len(keys) == 0 {
		return nil, fmt.Errorf("文件路径列表不能为空")
	}

	mac := qbox.NewMac(accessKey, secretKey)
	var urls []string

	// 对于公开空间，直接拼接 URL
	if !isPrivate {
		for _, key := range keys {
			if key == "" {
				continue
			}
			// 确保 URL 格式正确，避免双斜杠
			if len(imgUrl) > 0 && imgUrl[len(imgUrl)-1] == '/' {
				urls = append(urls, imgUrl+key)
			} else {
				urls = append(urls, imgUrl+"/"+key)
			}
		}
	} else {
		// 对于私有空间，生成带签名的临时外链
		if expire <= 0 {
			expire = 3600 // 默认 1 小时
		}
		domain := imgUrl
		for _, key := range keys {
			if key == "" {
				continue
			}
			privateURL := storage.MakePrivateURL(mac, domain, key, expire)
			urls = append(urls, privateURL)
		}
	}

	return urls, nil
}
