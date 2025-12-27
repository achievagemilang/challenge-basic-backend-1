package messaging

import (
	"encoding/json"

	"challenge-backend-1/internal/model"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type AddressConsumer struct {
	Log *zap.SugaredLogger
}

func NewAddressConsumer(log *zap.SugaredLogger) *AddressConsumer {
	return &AddressConsumer{
		Log: log,
	}
}

func (c AddressConsumer) Consume(message *sarama.ConsumerMessage) error {
	AddressEvent := new(model.AddressEvent)
	if err := json.Unmarshal(message.Value, AddressEvent); err != nil {
		c.Log.Errorw("error unmarshalling Address event", "error", err)
		return err
	}

	// TODO process event
	c.Log.Infof("Received topic addresses with event: %v from partition %d", AddressEvent, message.Partition)
	return nil
}
