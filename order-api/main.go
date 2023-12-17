package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/caarlos0/env/v9"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"

	"github.com/lorenzoranucci/golang-monorepo-project-template/lib/http/param"
)

type config struct {
	APIPort  string    `env:"API_PORT" envDefault:"8080"`
	LogLevel log.Level `env:"LOG_LEVEL" envDefault:"info"`
}

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	cfg := config{}
	if err := env.ParseWithOptions(&cfg, env.Options{RequiredIfNoDef: true}); err != nil {
		log.Fatal("error loading configs:", err)
	}
	log.SetLevel(cfg.LogLevel)

	router := httprouter.New()
	router.Handler(
		http.MethodGet,
		"/orders/:orderUUID",
		&GetOrderHandler{
			ParamFromRequest: param.FromRequestContext,
		})

	log.Info("starting server")
	server := &http.Server{
		Addr:         ":" + cfg.APIPort,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler:      router,
	}

	err := server.ListenAndServe()
	log.Fatal(fmt.Errorf("failed serving http: %w", err))
}
