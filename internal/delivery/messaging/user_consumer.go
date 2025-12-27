package messaging

import (
	"encoding/json"

	"challenge-backend-1/internal/model"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type UserConsumer struct {
	Log *zap.SugaredLogger
}

func NewUserConsumer(log *zap.SugaredLogger) *UserConsumer {
	return &UserConsumer{
		Log: log,
	}
}

func (c UserConsumer) Consume(message *sarama.ConsumerMessage) error {
	UserEvent := new(model.UserEvent)
	if err := json.Unmarshal(message.Value, UserEvent); err != nil {
		c.Log.Errorw("error unmarshalling User event", "error", err)
		return err
	}

	// TODO process event
	c.Log.Infof("Received topic users with event: %v from partition %d", UserEvent, message.Partition)
	return nil
}
