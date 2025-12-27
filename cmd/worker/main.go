package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"challenge-backend-1/internal/config"
	"challenge-backend-1/internal/delivery/messaging"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	viperConfig := config.NewViper()
	logger := config.NewLogger(viperConfig)
	logger.Info("Starting worker service")

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	wg.Add(3)
	go RunUserConsumer(logger, viperConfig, ctx, wg)
	go RunContactConsumer(logger, viperConfig, ctx, wg)
	go RunAddressConsumer(logger, viperConfig, ctx, wg)

	terminateSignals := make(chan os.Signal, 1)
	signal.Notify(terminateSignals, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

	stop := false
	for !stop {
		select {
		case s := <-terminateSignals:
			logger.Info("Got one of stop signals, shutting down worker gracefully, SIGNAL NAME :", s)
			cancel()
			stop = true
		}
	}

	wg.Wait()
	logger.Info("Worker exited")
}

func RunAddressConsumer(logger *zap.SugaredLogger, viperConfig *viper.Viper, ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	logger.Info("setup address consumer")
	addressConsumerGroup := config.NewKafkaConsumerGroup(viperConfig, logger)
	addressHandler := messaging.NewAddressConsumer(logger)
	messaging.ConsumeTopic(ctx, addressConsumerGroup, "addresses", logger, addressHandler.Consume)
}

func RunContactConsumer(logger *zap.SugaredLogger, viperConfig *viper.Viper, ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	logger.Info("setup contact consumer")
	contactConsumerGroup := config.NewKafkaConsumerGroup(viperConfig, logger)
	contactHandler := messaging.NewContactConsumer(logger)
	messaging.ConsumeTopic(ctx, contactConsumerGroup, "contacts", logger, contactHandler.Consume)
}

func RunUserConsumer(logger *zap.SugaredLogger, viperConfig *viper.Viper, ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	logger.Info("setup user consumer")
	userConsumerGroup := config.NewKafkaConsumerGroup(viperConfig, logger)
	userHandler := messaging.NewUserConsumer(logger)
	messaging.ConsumeTopic(ctx, userConsumerGroup, "users", logger, userHandler.Consume)
}
