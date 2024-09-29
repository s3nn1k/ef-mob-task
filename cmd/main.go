package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/s3nn1k/ef-mob-task/internal/app"
	"github.com/s3nn1k/ef-mob-task/internal/config"
	"github.com/s3nn1k/ef-mob-task/internal/dummy"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("[ERROR] Can't load environment: %s", err.Error())
	}
}

func main() {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatalf("[ERROR] Can't load config: %s", err.Error())
	}

	app, err := app.New(cfg)
	if err != nil {
		log.Fatalf("[ERROR] Can't create app: %s", err.Error())
	}

	dummy := dummy.New(cfg.API)

	go func() {
		_ = dummy.Run()
	}()

	go func() {
		if err := app.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Println(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	err = app.Stop()
	if err != nil {
		log.Println(err)
	}

	_ = dummy.Stop()
}
