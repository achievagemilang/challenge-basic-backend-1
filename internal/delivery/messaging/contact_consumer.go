package messaging

import (
	"encoding/json"

	"challenge-backend-1/internal/model"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type ContactConsumer struct {
	Log *zap.SugaredLogger
}

func NewContactConsumer(log *zap.SugaredLogger) *ContactConsumer {
	return &ContactConsumer{
		Log: log,
	}
}

func (c ContactConsumer) Consume(message *sarama.ConsumerMessage) error {
	ContactEvent := new(model.ContactEvent)
	if err := json.Unmarshal(message.Value, ContactEvent); err != nil {
		c.Log.Errorw("error unmarshalling Contact event", "error", err)
		return err
	}

	// TODO process event
	c.Log.Infof("Received topic contacts with event: %v from partition %d", ContactEvent, message.Partition)
	return nil
}
