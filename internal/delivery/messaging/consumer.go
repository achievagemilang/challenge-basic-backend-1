package messaging

import (
	"context"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type ConsumerHandler func(message *sarama.ConsumerMessage) error

type ConsumerGroupHandler struct {
	Handler ConsumerHandler
	Log     *zap.SugaredLogger
}

func (h *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			if message == nil {
				return nil
			}

			err := h.Handler(message)
			if err != nil {
				h.Log.Errorw("Failed to process message", "error", err)
			} else {
				session.MarkMessage(message, "")
			}

		case <-session.Context().Done():
			return nil
		}
	}
}

func ConsumeTopic(ctx context.Context, consumerGroup sarama.ConsumerGroup, topic string, log *zap.SugaredLogger, handler ConsumerHandler) {
	consumerHandler := &ConsumerGroupHandler{
		Handler: handler,
		Log:     log,
	}

	go func() {
		for {
			if err := consumerGroup.Consume(ctx, []string{topic}, consumerHandler); err != nil {
				log.Errorw("Error from consumer", "error", err)
			}

			if ctx.Err() != nil {
				log.Info("Context cancelled, stopping consumer")
				return
			}
		}
	}()

	go func() {
		for err := range consumerGroup.Errors() {
			log.Errorw("Consumer group error", "error", err)
		}
	}()

	<-ctx.Done()
	log.Infof("Closing consumer group for topic: %s", topic)
	if err := consumerGroup.Close(); err != nil {
		log.Errorw("Error closing consumer group", "error", err)
	}
}
