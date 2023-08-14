package calendar

import (
	"context"
	"github.com/evgen1067/hw12_13_14_15_calendar/internal/server/grpc"
	"github.com/evgen1067/hw12_13_14_15_calendar/internal/storage/memory"
	"os/signal"
	"syscall"

	"github.com/evgen1067/hw12_13_14_15_calendar/internal/config"
	"github.com/evgen1067/hw12_13_14_15_calendar/internal/logger"
	"github.com/evgen1067/hw12_13_14_15_calendar/internal/server/rest"
	"github.com/evgen1067/hw12_13_14_15_calendar/internal/services"
	"github.com/evgen1067/hw12_13_14_15_calendar/internal/storage"
	psql "github.com/evgen1067/hw12_13_14_15_calendar/internal/storage/sql"
)

func Run(cfg *config.Config, logg *logger.Logger) error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	errs := make(chan error)

	var db storage.EventsStorage
	if cfg.SQL {
		db = psql.NewStorage()
	} else {
		db = memory.NewStorage()
	}

	if r, ok := db.(storage.DBStorage); ok {
		err := r.Connect(ctx, cfg)
		if err != nil {
			return err
		}
		logg.Info("Database started.")
		defer r.Close()
	}

	service := services.NewServices(ctx, db, logg)

	restAPI := rest.NewServer(service, cfg)

	go func() {
		logg.Info("HTTP server started.")
		err := restAPI.ListenAndServe()
		if err != nil {
			errs <- err
		}
	}()

	grpcAPI := grpc.NewGRPC(service, cfg)

	go func() {
		logg.Info("GRPC server started.")
		err := grpcAPI.ListenAndServe()
		if err != nil {
			errs <- err
		}
	}()

	// Выползаем при ошибке или завершении программы
	select {
	case err := <-errs:
		return err
	case <-ctx.Done():
	}

	grpcAPI.Stop()

	if err := restAPI.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
