package qcloudutil

import (
	"context"
	"github.com/letjoy-club/mida-tool/logger"
	cls "github.com/tencentcloud/tencentcloud-cls-sdk-go"
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

func (q QCloudConf) Init() {
	if q.CLS != nil {
		if err := q.initClsClient(); err != nil {
			logger.L.Error("Init CLS client error", zap.Error(err))
		}
	}
	if q.COS != nil {
		if err := q.initCosClient(); err != nil {
			logger.L.Error("Init COS client error", zap.Error(err))
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

type qCloudKey struct{}

func WithQCloud(ctx context.Context, conf QCloudConf) context.Context {
	return context.WithValue(ctx, qCloudKey{}, conf)
}

func GetQCloud(ctx context.Context) QCloudConf {
	return ctx.Value(qCloudKey{}).(QCloudConf)
}
