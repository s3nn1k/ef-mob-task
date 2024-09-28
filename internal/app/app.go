package app

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/s3nn1k/ef-mob-task/internal/client"
	"github.com/s3nn1k/ef-mob-task/internal/config"
	"github.com/s3nn1k/ef-mob-task/internal/delivery"
	"github.com/s3nn1k/ef-mob-task/internal/delivery/middleware"
	"github.com/s3nn1k/ef-mob-task/internal/service"
	"github.com/s3nn1k/ef-mob-task/internal/storage/postgres"
	"github.com/s3nn1k/ef-mob-task/pkg/logger"
)

type App struct {
	db     *pgxpool.Pool
	server *http.Server
}

func (a *App) Run() error {
	return a.server.ListenAndServe()
}

func (a *App) Stop() error {
	a.db.Close()

	err := a.server.Shutdown(context.Background())
	if err != nil {
		return err
	}

	return nil
}

// Add migrations
func New(cfg *config.Config) (*App, error) {
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", cfg.DB.User, cfg.DB.Pass, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name)

	db, err := postgres.ConnectDB(connStr)
	if err != nil {
		return nil, err
	}

	schemaPath := "./schema/schema.sql"

	err = initTables(schemaPath, db)
	if err != nil {
		return nil, err
	}

	strg := postgres.NewStorage(db)

	clnt := client.New(cfg.API.Host, cfg.API.Port)

	srvc := service.New(strg, clnt)

	logger := logger.NewTextLogger(cfg.Level)

	hndlr := delivery.NewHandler(logger, srvc)

	r := initRoutes(hndlr, logger)

	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)

	return &App{
		db: db,
		server: &http.Server{
			Addr:           addr,
			MaxHeaderBytes: 1 << 20,
			Handler:        r,
			WriteTimeout:   cfg.Server.Timeout,
			ReadTimeout:    cfg.Server.Timeout,
			IdleTimeout:    cfg.Server.IdleTimeout,
		},
	}, nil
}

func initRoutes(h *delivery.Handler, log *slog.Logger) *http.ServeMux {
	router := http.NewServeMux()

	router.Handle("POST /songs/", middleware.WithLogging(log, http.HandlerFunc(h.Create)))
	router.Handle("PUT /songs/", middleware.WithLogging(log, http.HandlerFunc(h.Update)))
	router.Handle("GET /songs/", middleware.WithLogging(log, http.HandlerFunc(h.Get)))
	router.Handle("DELETE /songs/", middleware.WithLogging(log, http.HandlerFunc(h.Delete)))

	return router
}

func initTables(schemaPath string, db *pgxpool.Pool) error {
	schemaFile, err := os.Open(schemaPath)
	if err != nil {
		return err
	}

	defer schemaFile.Close()

	data, err := io.ReadAll(schemaFile)
	if err != nil {
		return err
	}

	_, err = db.Exec(context.Background(), string(data))
	if err != nil {
		return err
	}

	return nil
}
