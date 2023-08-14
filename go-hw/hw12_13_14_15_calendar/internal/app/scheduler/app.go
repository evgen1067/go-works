package scheduler

import (
	"context"
	"fmt"
	"github.com/evgen1067/hw12_13_14_15_calendar/internal/config"
	"github.com/evgen1067/hw12_13_14_15_calendar/internal/logger"
	"github.com/evgen1067/hw12_13_14_15_calendar/internal/rabbit/producer"
	psql "github.com/evgen1067/hw12_13_14_15_calendar/internal/storage/sql"
	"os/signal"
	"syscall"
	"time"
)

func Run(cfg *config.Config, logg *logger.Logger) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logg.Info("The scheduler has started working")
	defer logg.Info("The scheduler has finished its work")

	db := psql.NewStorage()
	err := db.Connect(ctx, cfg)
	if err != nil {
		logg.Error("Error when connecting to the database: " + err.Error())
	}
	defer db.Close()

	prod := producer.NewProducer(cfg.AMQP.URI, cfg.AMQP.Queue)
	err = prod.Start()
	if err != nil {
		logg.Error(err.Error())
		return
	}
	defer prod.Stop()

	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		for {
			err := db.ClearOldEvents(ctx)
			if err != nil {
				logg.Error("Error when clearing old events: " + err.Error())
			}
			notices, err := db.SchedulerList(ctx)
			if err != nil {
				logg.Error("Error when receiving notifications from the database: " + err.Error())
			}
			for _, v := range notices {
				n, err := v.MarshalJSON()
				if err != nil {
					logg.Error("Error when marshaling notifications: " + err.Error())
				}
				err = prod.Publish(ctx, n)
				if err != nil {
					logg.Error("Error when publishing a notification by the producer: " + err.Error())
				} else {
					logg.Info(fmt.Sprintf("[x] Sent %s", n))
				}
			}
			select {
			case <-ticker.C:
				continue
			case <-ctx.Done():
				return
			}
		}
	}()

	<-ctx.Done()
}
