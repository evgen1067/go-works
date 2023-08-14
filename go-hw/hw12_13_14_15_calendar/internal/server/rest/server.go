package rest

import (
	"net"
	"net/http"
	"time"

	"github.com/evgen1067/hw12_13_14_15_calendar/internal/config"
	"github.com/evgen1067/hw12_13_14_15_calendar/internal/services"
	"github.com/gorilla/mux"
)

var s *services.Services

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", HelloWorld).Methods(http.MethodGet)
	router.HandleFunc("/events/new", CreateEvent).Methods(http.MethodPost)
	router.HandleFunc("/events/{id}", UpdateEvent).Methods(http.MethodPut)
	router.HandleFunc("/events/{id}", DeleteEvent).Methods(http.MethodDelete)
	router.HandleFunc("/events/list/{period}", EventList).Methods(http.MethodGet)

	router.NotFoundHandler = router.NewRoute().HandlerFunc(CustomNotFoundHandler).GetHandler()

	router.Use(headersMiddleware)
	router.Use(loggerMiddleware)

	return router
}

func NewServer(_s *services.Services, cfg *config.Config) *http.Server {
	s = _s
	return &http.Server{
		Addr:              net.JoinHostPort(cfg.HTTP.Host, cfg.HTTP.Port),
		Handler:           Router(),
		ReadHeaderTimeout: 3 * time.Second,
	}
}
