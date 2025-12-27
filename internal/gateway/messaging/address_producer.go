package messaging

import (
	"challenge-backend-1/internal/model"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type AddressProducer struct {
	Producer[*model.AddressEvent]
}

func NewAddressProducer(producer sarama.SyncProducer, log *zap.SugaredLogger) *AddressProducer {
	return &AddressProducer{
		Producer: Producer[*model.AddressEvent]{
			Producer: producer,
			Topic:    "addresses",
			Log:      log,
		},
	}
}
