package dummy

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/s3nn1k/ef-mob-task/internal/config"
	"github.com/s3nn1k/ef-mob-task/internal/models"
)

// Dummy is a mock of api service
type Dummy struct {
	server *http.Server
}

func (d *Dummy) Run() error {
	return d.server.ListenAndServe()
}

func (d *Dummy) Stop() error {
	return d.server.Shutdown(context.Background())
}

func New(cfg config.API) *Dummy {
	mux := http.NewServeMux()

	mux.Handle("GET /info", dummyHandler())

	return &Dummy{
		server: &http.Server{
			Addr:           cfg.Host + ":" + cfg.Port,
			MaxHeaderBytes: 1 << 20,
			Handler:        mux,
			WriteTimeout:   4 * time.Second,
			ReadTimeout:    4 * time.Second,
			IdleTimeout:    60 * time.Second,
		},
	}
}

func dummyHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var song models.Song

		_ = json.NewDecoder(r.Body).Decode(&song)

		song.Text = "first verse\n\nsecond verse\n\nthird verse\n\nfourth verse\n\n"
		song.Link = "https://www.youtube.com/watch?v=HIcSWuKMwOw"
		song.Date = time.Now().Format("02.01.2006")

		data, _ := json.Marshal(song)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})
}
