/**
 * @author: yangchangjia
 * @email 1320259466@qq.com
 * @date: 2024/5/15 10:04
 * @desc: about the role of class.
 */

package tlogin

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
)

var aliClient *AliOssClient

type AliOssClient struct {
	*oss.Client
}

var (
	// oss bucket名
	ossBucket = "pibigstar"

	// oss endpoint
	//ossEndpoint = "oss-cn-shanghai.aliyuncs.com"

	// oss访问key
	//ossAccessKeyID = "LTAIKFU1CUmLErUw"

	// oss private key secret
	//ossAccessKeySecret = "n0axekSPgKwCqIGyBa1oSZBQpOyzlp"

	// 默认失效时间，30天
	defaultBucketExpireTime = 30
)

// NewAliOSSClient : 创建oss client对象
func NewAliOSSClient(bucketName, endpoint, accessKeyID, accessKeySecret string, bucketExpireTime int) (*AliOssClient, error) {
	if aliClient != nil {
		return aliClient, nil
	}
	ossBucket = bucketName
	defaultBucketExpireTime = bucketExpireTime
	ossCli, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		return nil, err
	}
	aliClient = &AliOssClient{Client: ossCli}
	return aliClient, nil
}

// GetBucket 获取bucket存储空间
func (client *AliOssClient) GetBucket(bucketNames ...string) *oss.Bucket {
	if client != nil {
		bucketName := ossBucket
		if len(bucketNames) != 0 {
			bucketName = bucketNames[0]
		}
		bucket, err := client.Bucket(bucketName)
		if err != nil {
			return nil
		}
		return bucket
	}
	return nil
}

// Put put the object to the oss
func (client *AliOssClient) Put(key string, reader io.Reader) error {
	return client.GetBucket().PutObject(key, reader)
}

// DeleteObject delete the object, actually use the nil replace the object
func (client *AliOssClient) DeleteObject(key string) {
	client.Put(key, nil)
}

// GetDownloadURL get the down url
func (client *AliOssClient) GetDownloadURL(objName string) (string, error) {
	signedURL, err := client.GetBucket().SignURL(objName, oss.HTTPGet, 3600)
	if err != nil {
		return "", err
	}
	return signedURL, nil
}

// BuildLifecycleRule set lifecycle rule for the specified bucket
func (client *AliOssClient) BuildLifecycleRule(bucketName string) {
	// 表示前缀为test的对象(文件)距最后修改时间30天后过期。
	ruleTest1 := oss.BuildLifecycleRuleByDays("rule1", "test/", true, defaultBucketExpireTime)
	rules := []oss.LifecycleRule{ruleTest1}
	client.SetBucketLifecycle(bucketName, rules)
}
