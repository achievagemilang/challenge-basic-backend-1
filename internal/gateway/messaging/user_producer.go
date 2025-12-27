package messaging

import (
	"challenge-backend-1/internal/model"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type UserProducer struct {
	Producer[*model.UserEvent]
}

func NewUserProducer(producer sarama.SyncProducer, log *zap.SugaredLogger) *UserProducer {
	return &UserProducer{
		Producer: Producer[*model.UserEvent]{
			Producer: producer,
			Topic:    "users",
			Log:      log,
		},
	}
}
