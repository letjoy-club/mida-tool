package qcloudutil

import (
	"context"
	"github.com/letjoy-club/mida-tool/logger"
	cls "github.com/tencentcloud/tencentcloud-cls-sdk-go"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"github.com/tencentyun/cos-go-sdk-v5"
	"go.uber.org/zap"
	"net/http"
	"net/url"
)

type QCloudConf struct {
	CDN string   `yaml:"cdn"`
	CLS *CLSConf `yaml:"cls"`
	COS *COSConf `yaml:"cos"`
	TIM *TIMConf `yaml:"tim"`
	SMS *SMSConf `yaml:"sms"`
}

type CLSConf struct {
	Endpoint  string `yaml:"endpoint"`
	SecretID  string `yaml:"secret-id"`
	SecretKey string `yaml:"secret-key"`
	TopicID   string `yaml:"topic-id"`
	Client    *cls.AsyncProducerClient
}

type COSConf struct {
	Endpoint  string `yaml:"endpoint"`
	SecretID  string `yaml:"secret-id"`
	SecretKey string `yaml:"secret-key"`
	Path      string `yaml:"path"`
	Client    *cos.Client
}

type TIMConf struct {
	AppID int    `yaml:"app-id"`
	Key   string `yaml:"key"`
}

type SMSConf struct {
	AppID     string `yaml:"app-id"`
	SecretID  string `yaml:"secret-id"`
	SecretKey string `yaml:"secret-key"`
	SignName  string `yaml:"sign-name"`
	Client    *sms.Client
}

func (q QCloudConf) Init() {
	if q.CLS != nil { // 日志
		if err := q.initClsClient(); err != nil {
			logger.L.Error("Init CLS client error", zap.Error(err))
		}
	}
	if q.COS != nil { // 对象存储
		if err := q.initCosClient(); err != nil {
			logger.L.Error("Init COS client error", zap.Error(err))
		}
	}
	if q.SMS != nil { // 短信
		if err := q.initSmsClient(); err != nil {
			logger.L.Error("Init SMS client error", zap.Error(err))
		}
	}
}

func (q QCloudConf) initClsClient() error {
	c := cls.GetDefaultAsyncProducerClientConfig()
	c.Endpoint = q.CLS.Endpoint
	c.AccessKeyID = q.CLS.SecretID
	c.AccessKeySecret = q.CLS.SecretKey
	c.Retries = 1
	client, err := cls.NewAsyncProducerClient(c)
	if err != nil {
		return err
	}
	q.CLS.Client = client
	return nil
}

func (q QCloudConf) initCosClient() error {
	u, err := url.Parse(q.COS.Endpoint)
	if err != nil {
		return err
	}
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  q.COS.SecretID,
			SecretKey: q.COS.SecretKey,
		},
	})
	q.COS.Client = client
	return nil
}

func (q QCloudConf) initSmsClient() error {
	credential := common.NewCredential(q.SMS.SecretID, q.SMS.SecretKey)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"
	client, _ := sms.NewClient(credential, "ap-nanjing", cpf)
	q.SMS.Client = client
	return nil
}

type qCloudKey struct{}

func WithQCloud(ctx context.Context, conf QCloudConf) context.Context {
	return context.WithValue(ctx, qCloudKey{}, conf)
}

func GetQCloud(ctx context.Context) QCloudConf {
	return ctx.Value(qCloudKey{}).(QCloudConf)
}
