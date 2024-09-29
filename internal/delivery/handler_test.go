package delivery

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/s3nn1k/ef-mob-task/internal/models"
	"github.com/s3nn1k/ef-mob-task/internal/service/mocks"
	"github.com/s3nn1k/ef-mob-task/pkg/logger"
	"github.com/s3nn1k/ef-mob-task/pkg/test"
)

func TestCreate(t *testing.T) {
	mock := mocks.NewServiceIface(t)

	song := models.Song{
		Id:    1,
		Song:  "TestSong",
		Group: "TestGroup",
		Text:  "TestText\n\nTestText",
		Link:  "TestLink",
		Date:  time.Now().Format("02.01.2006"),
	}

	log := logger.NewTextLogger("")

	mock.On("Create", logger.NewCtxWithLog(context.Background(), log), song.Song, song.Group).
		Return(song, nil)

	testCases := []test.TestCase{
		{
			Name:       "success",
			Body:       `{"song":"TestSong", "group":"TestGroup"}`,
			WantStatus: 200,
			WantRes:    fmt.Sprintf(`{"status":"Ok","result":[{"id":1,"song":"TestSong","group":"TestGroup","text":"TestText\n\nTestText","link":"TestLink","releaseDate":"%s"}]}`, song.Date),
		},
		{
			Name:       "wrongBody",
			WantStatus: 400,
			WantRes:    `{"status":"Error","error":"Can't decode json body"}`,
		},
	}

	handler := NewHandler(log, mock)

	router := http.NewServeMux()

	router.HandleFunc("POST /songs", http.HandlerFunc(handler.Create))

	for _, testCase := range testCases {
		testCase.Url = "/songs"
		testCase.Method = "POST"

		test.TestEndpoint(t, router, testCase)
	}
}

func TestGetAll(t *testing.T) {
	mock := mocks.NewServiceIface(t)

	song := models.Song{
		Id:    1,
		Song:  "TestSong",
		Group: "TestGroup",
		Text:  "TestText TestText",
		Link:  "TestLink",
		Date:  "TestDate",
	}

	filters := models.AllFilters{
		Limit:  1,
		Offset: 1,
		Song:   song.Song,
		Group:  song.Group,
		Date:   song.Date,
	}

	log := logger.NewTextLogger("")

	mock.On("GetAll", logger.NewCtxWithLog(context.Background(), log), filters).
		Return([]models.Song{song}, nil)

	testCases := []test.TestCase{
		{
			Name:       "success",
			Url:        "/songs?song=TestSong&group=TestGroup&date=TestDate&limit=1&offset=1",
			WantStatus: 200,
			WantRes:    `{"status":"Ok","result":[{"id":1,"song":"TestSong","group":"TestGroup","text":"TestText TestText","link":"TestLink","releaseDate":"TestDate"}]}`,
		},
		{
			Name:       "invalid filters",
			Url:        "/songs?limit=one&offset=two",
			WantStatus: 400,
			WantRes:    `{"status":"Error","error":"limit and offset must be int"}`,
		},
	}

	handler := NewHandler(log, mock)

	router := http.NewServeMux()

	router.HandleFunc("GET /songs", http.HandlerFunc(handler.GetAll))

	for _, testCase := range testCases {
		testCase.Method = "GET"

		test.TestEndpoint(t, router, testCase)
	}
}

func TestGetById(t *testing.T) {
	mock := mocks.NewServiceIface(t)

	song := models.Song{
		Id:    1,
		Song:  "TestSong",
		Group: "TestGroup",
		Text:  "TestText TestText",
		Link:  "TestLink",
		Date:  "TestDate",
	}

	filters := models.SongFilters{
		Id:    1,
		Verse: 1,
	}

	log := logger.NewTextLogger("")

	mock.On("GetById", logger.NewCtxWithLog(context.Background(), log), filters).
		Return(song, nil)

	testCases := []test.TestCase{
		{
			Name:       "success",
			Url:        "/songs/1?verse=1",
			WantStatus: 200,
			WantRes:    `{"status":"Ok","result":[{"id":1,"song":"TestSong","group":"TestGroup","text":"TestText TestText","link":"TestLink","releaseDate":"TestDate"}]}`,
		},
		{
			Name:       "invalid id",
			Url:        "/songs/kasjdf?verse=1",
			WantStatus: 400,
			WantRes:    `{"status":"Error","error":"id must be int"}`,
		},
		{
			Name:       "invalid verse",
			Url:        "/songs/1?verse=asdf",
			WantStatus: 400,
			WantRes:    `{"status":"Error","error":"verse must be int"}`,
		},
	}

	handler := NewHandler(log, mock)

	router := http.NewServeMux()

	router.HandleFunc("GET /songs/{id}", http.HandlerFunc(handler.GetById))

	for _, testCase := range testCases {
		testCase.Method = "GET"

		test.TestEndpoint(t, router, testCase)
	}
}

func TestUpdate(t *testing.T) {
	mock := mocks.NewServiceIface(t)

	song := models.Song{
		Id:    1,
		Song:  "TestSong",
		Group: "TestGroup",
		Text:  "TestText TestText",
		Link:  "TestLink",
		Date:  "TestDate",
	}

	log := logger.NewTextLogger("")

	mock.On("Update", logger.NewCtxWithLog(context.Background(), log), song).
		Return(true, nil)

	testCases := []test.TestCase{
		{
			Name:       "success",
			Url:        "/songs",
			Body:       `{"id":1,"song":"TestSong", "group":"TestGroup","text":"TestText TestText","link":"TestLink","releaseDate":"TestDate"}`,
			WantStatus: 200,
			WantRes:    `{"status":"Ok"}`,
		},
		{
			Name:       "invalid id",
			Url:        "/songs",
			Body:       `{"id":"zero"}`,
			WantStatus: 400,
			WantRes:    `{"status":"Error","error":"Can't decode json body"}`,
		},
	}

	handler := NewHandler(log, mock)

	router := http.NewServeMux()

	router.HandleFunc("PUT /songs", http.HandlerFunc(handler.Update))

	for _, testCase := range testCases {
		testCase.Method = "PUT"

		test.TestEndpoint(t, router, testCase)
	}
}

func TestFailUpdate(t *testing.T) {
	mock := mocks.NewServiceIface(t)

	song := models.Song{
		Id: 1,
	}

	log := logger.NewTextLogger("")

	mock.On("Update", logger.NewCtxWithLog(context.Background(), log), song).
		Return(false, nil)

	testCase := test.TestCase{
		Name:       "fail",
		Url:        "/songs",
		Method:     "PUT",
		Body:       `{"id":1}`,
		WantStatus: 400,
		WantRes:    `{"status":"Error","error":"Song not exists"}`,
	}

	handler := NewHandler(log, mock)

	router := http.NewServeMux()

	router.HandleFunc("PUT /songs", http.HandlerFunc(handler.Update))

	test.TestEndpoint(t, router, testCase)
}

func TestDelete(t *testing.T) {
	mock := mocks.NewServiceIface(t)

	id := 1

	log := logger.NewTextLogger("")

	mock.On("Delete", logger.NewCtxWithLog(context.Background(), log), id).
		Return(true, nil)

	testCases := []test.TestCase{
		{
			Name:       "success",
			Url:        "/songs/1",
			Body:       `{}`,
			WantStatus: 204,
			WantRes:    `{"status":"Ok"}`,
		},
		{
			Name:       "invalid id",
			Url:        "/songs/ieunf",
			Body:       `{}`,
			WantStatus: 400,
			WantRes:    `{"status":"Error","error":"id must be int"}`,
		},
	}

	handler := NewHandler(log, mock)

	router := http.NewServeMux()

	router.HandleFunc("DELETE /songs/{id}", http.HandlerFunc(handler.Delete))

	for _, testCase := range testCases {
		testCase.Method = "DELETE"

		test.TestEndpoint(t, router, testCase)
	}
}

func TestFailDelete(t *testing.T) {
	mock := mocks.NewServiceIface(t)

	id := 1

	log := logger.NewTextLogger("")

	mock.On("Delete", logger.NewCtxWithLog(context.Background(), log), id).
		Return(false, nil)

	testCase := test.TestCase{
		Name:       "fail",
		Url:        "/songs/1",
		Method:     "DELETE",
		Body:       `{}`,
		WantStatus: 400,
		WantRes:    `{"status":"Error","error":"Song not exists"}`,
	}

	handler := NewHandler(log, mock)

	router := http.NewServeMux()

	router.HandleFunc("DELETE /songs/{id}", http.HandlerFunc(handler.Delete))

	test.TestEndpoint(t, router, testCase)
}
