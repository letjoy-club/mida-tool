package pulsarutil

import (
	"context"
	"fmt"
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/fatih/color"
	"log"
	"time"
)

type PulsarClient struct {
	User      string `yaml:"user"`
	Token     string `yaml:"token"`
	Endpoint  string `yaml:"endpoint"`
	Namespace string `yaml:"namespace"`
	Tenant    string `yaml:"tenant"`

	client pulsar.Client
}

// CreateProducer 创建pulsar生产者
func (m *PulsarClient) CreateProducer(topic string) pulsar.Producer {
	if m.client == nil {
		m.Init()
	}
	if m.client == nil {
		return nil
	}
	producer, err := m.client.CreateProducer(pulsar.ProducerOptions{
		Topic: m.fullTopic(topic),
	})
	if err != nil {
		log.Panicln(err)
	}
	return producer
}

// CreateConsumer 创建pulsar消费者
func (m *PulsarClient) CreateConsumer(topic string, subscriber string) pulsar.Consumer {
	if m.client == nil {
		m.Init()
	}
	if m.client == nil {
		return nil
	}
	consumer, err := m.client.Subscribe(pulsar.ConsumerOptions{
		SubscriptionName:            subscriber,
		Topic:                       m.fullTopic(topic),
		Type:                        pulsar.Shared,
		SubscriptionInitialPosition: pulsar.SubscriptionPositionLatest,
	})
	if err != nil {
		log.Panicln(err)
	}
	return consumer
}

func (m *PulsarClient) fullTopic(topic string) string {
	return fmt.Sprintf("persistent://%s/%s/%s", m.Tenant, m.Namespace, topic)
}

func (m *PulsarClient) Init() {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		// 服务接入地址
		URL: m.Endpoint,
		// 授权角色密钥
		Authentication:    pulsar.NewAuthenticationToken(m.Token),
		OperationTimeout:  30 * time.Second,
		ConnectionTimeout: 30 * time.Second,
	})
	if err != nil {
		log.Println(color.RedString("failed to connect to pulsar: %s", err))
	}
	m.client = client
}

type mqKey struct{}

func WithMQ[MQConfig any](ctx context.Context, config MQConfig) context.Context {
	return context.WithValue(ctx, mqKey{}, config)
}

func GetMQ[MQConfig any](ctx context.Context) MQConfig {
	return ctx.Value(mqKey{}).(MQConfig)
}
