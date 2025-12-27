package messaging

import (
	"challenge-backend-1/internal/model"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type ContactProducer struct {
	Producer[*model.ContactEvent]
}

func NewContactProducer(producer sarama.SyncProducer, log *zap.SugaredLogger) *ContactProducer {
	return &ContactProducer{
		Producer: Producer[*model.ContactEvent]{
			Producer: producer,
			Topic:    "contacts",
			Log:      log,
		},
	}
}
