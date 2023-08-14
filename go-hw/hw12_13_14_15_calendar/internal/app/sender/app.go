package sender

import (
	"context"
	"fmt"
	"github.com/evgen1067/hw12_13_14_15_calendar/internal/common"
	"github.com/evgen1067/hw12_13_14_15_calendar/internal/config"
	"github.com/evgen1067/hw12_13_14_15_calendar/internal/logger"
	"github.com/evgen1067/hw12_13_14_15_calendar/internal/rabbit/consumer"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Config, logg *logger.Logger) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logg.Info("The sender has started working")
	defer logg.Info("The sender has finished its work")

	cons := consumer.NewConsumer(cfg.AMQP.URI, cfg.AMQP.Queue)
	err := cons.Start()
	if err != nil {
		logg.Error(err.Error())
		return
	}
	defer cons.Stop()

	messages, err := cons.Consume()
	if err != nil {
		logg.Error(err.Error())
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-messages:
				var notice common.Notice
				err := notice.UnmarshalJSON(msg.Body)
				if err != nil {
					logg.Error("Error when unmarshaling notifications: " + err.Error())
					continue
				}
				logg.Info(fmt.Sprintf("ID: %v, Title: %v, Datetime: %v, OwnerID: %v",
					notice.EventID, notice.Title, notice.Datetime, notice.OwnerID))
			}
		}
	}()

	<-ctx.Done()
}
