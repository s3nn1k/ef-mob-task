package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/s3nn1k/ef-mob-task/docs"
	"github.com/s3nn1k/ef-mob-task/internal/client"
	"github.com/s3nn1k/ef-mob-task/internal/config"
	"github.com/s3nn1k/ef-mob-task/internal/delivery"
	"github.com/s3nn1k/ef-mob-task/internal/delivery/middleware"
	"github.com/s3nn1k/ef-mob-task/internal/service"
	"github.com/s3nn1k/ef-mob-task/internal/storage/postgres"
	"github.com/s3nn1k/ef-mob-task/pkg/logger"
	httpSwagger "github.com/swaggo/http-swagger"
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

// New creates new instance of application, sets the dependencies and applies migrations
func New(cfg *config.Config) (*App, error) {
	log := logger.NewTextLogger(cfg.Level)

	log.Info("Created logger", slog.String("level", cfg.Level))

	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", cfg.DB.User, cfg.DB.Pass, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name)

	m, err := migrate.New("file:///migrations", connStr)
	if err != nil {
		return nil, err
	}

	if err = m.Up(); err != nil {
		return nil, err
	}

	db, err := postgres.ConnectDB(connStr)
	if err != nil {
		return nil, err
	}

	log.Info("Connected to database and applied migrations", "config", cfg.DB.AsLogValue())

	strg := postgres.NewStorage(db)

	clnt := client.New(cfg.API.Host, cfg.API.Port)

	log.Info("Setup API client", "config", cfg.API.AsLogValue())

	srvc := service.New(strg, clnt)

	hndlr := delivery.NewHandler(log, srvc)

	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)

	r := initRoutes(hndlr, log)

	app := &App{
		db: db,
		server: &http.Server{
			Addr:           addr,
			MaxHeaderBytes: 1 << 20,
			Handler:        r,
			WriteTimeout:   cfg.Server.Timeout,
			ReadTimeout:    cfg.Server.Timeout,
			IdleTimeout:    cfg.Server.IdleTimeout,
		},
	}

	log.Info("Created app with server", "config", cfg.Server.AsLogValue())

	log.Info("Is test api service in use", slog.Bool("value", cfg.UseTestApi))

	log.Info("Visit " + fmt.Sprintf("http://localhost:%s/swagger/index.html", cfg.Server.Port) + " to test the API")

	return app, nil
}

func initRoutes(h *delivery.Handler, log *slog.Logger) *http.ServeMux {
	router := http.NewServeMux()

	router.Handle("POST /songs", middleware.WithLogging(log, http.HandlerFunc(h.Create)))
	router.Handle("PUT /songs/{id}", middleware.WithLogging(log, http.HandlerFunc(h.Update)))
	router.Handle("GET /songs", middleware.WithLogging(log, http.HandlerFunc(h.GetAll)))
	router.Handle("GET /songs/{id}", middleware.WithLogging(log, http.HandlerFunc(h.GetVerses)))
	router.Handle("DELETE /songs/{id}", middleware.WithLogging(log, http.HandlerFunc(h.Delete)))

	router.Handle("GET /swagger/", httpSwagger.WrapHandler)

	log.Info("Available routes", slog.Group("route",
		slog.String("Create", "POST /songs"),
		slog.String("Update", "PUT /songs/{id}"),
		slog.String("GetAll", "GET /songs"),
		slog.String("GetVerses", "GET /songs/{id}"),
		slog.String("Delete", "DELETE /songs/{id}"),
		slog.String("Swagger", "GET /swagger/")))

	return router
}
