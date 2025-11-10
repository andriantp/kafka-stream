package kafka

import (
	"context"

	"github.com/segmentio/kafka-go/sasl/plain"
)

type Setting struct {
	Broker  string
	Topic1  string
	Topic2  string
	GroupID string

	Username string
	Password string
}

type repository struct {
	setting Setting
	sasl    plain.Mechanism
}

type RepositoryI interface {
	Producer1(ctx context.Context, msg string) error
	Consumer1(ctx context.Context) error

	Producer2(ctx context.Context, msg string) error
	Consumer2(ctx context.Context) error
}

func NewKafka(setting Setting) (RepositoryI, error) {
	// SASL PLAIN config
	mechanism := plain.Mechanism{
		Username: setting.Username,
		Password: setting.Password,
	}

	return &repository{
		setting: setting,
		sasl:    mechanism,
	}, nil
}
