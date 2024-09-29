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
			Body:       ``,
			WantStatus: 400,
			WantRes:    `{"status":"Error","error":"Can't decode json body"}`,
		},
	}

	handler := NewHandler(log, mock)

	for _, testCase := range testCases {
		testCase.Url = "/songs"
		testCase.Method = "POST"

		test.TestEndpoint(t, http.HandlerFunc(handler.Create), testCase)
	}
}

func TestGet(t *testing.T) {
	mock := mocks.NewServiceIface(t)

	filter := models.Song{
		Id:    1,
		Song:  "TestSong",
		Group: "TestGroup",
		Text:  "TestText TestText",
		Link:  "TestLink",
		Date:  "TestDate",
	}

	filters := models.Filters{
		Verse:  1,
		Limit:  1,
		Offset: 1,
	}

	log := logger.NewTextLogger("")

	mock.On("Get", logger.NewCtxWithLog(context.Background(), log), filter, filters).
		Return([]models.Song{filter}, nil)

	testCases := []test.TestCase{
		{
			Name:       "success",
			Url:        "/songs?id=1&song=TestSong&group=TestGroup&text=TestText TestText&link=TestLink&date=TestDate&verse=1&limit=1&offset=1",
			Body:       ``,
			WantStatus: 200,
			WantRes:    `{"status":"Ok","result":[{"id":1,"song":"TestSong","group":"TestGroup","text":"TestText TestText","link":"TestLink","releaseDate":"TestDate"}]}`,
		},
		{
			Name:       "invalid filters",
			Url:        "/songs?verse=zero&limit=one&offset=two",
			Body:       ``,
			WantStatus: 400,
			WantRes:    `{"status":"Error","error":"verse, limit and offset must be int"}`,
		},
		{
			Name:       "invalid Id",
			Url:        "/songs?id=zero",
			Body:       ``,
			WantStatus: 400,
			WantRes:    `{"status":"Error","error":"id must be int"}`,
		},
	}

	handler := NewHandler(log, mock)

	for _, testCase := range testCases {
		testCase.Method = "GET"

		test.TestEndpoint(t, http.HandlerFunc(handler.Get), testCase)
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
			Name:       "json success",
			Url:        "/songs?id=1",
			Body:       `{"song":"TestSong", "group":"TestGroup","text":"TestText TestText","link":"TestLink","releaseDate":"TestDate"}`,
			WantStatus: 200,
			WantRes:    `{"status":"Ok"}`,
		},
		{
			Name:       "invalid id",
			Url:        "/songs?id=zero",
			Body:       `{}`,
			WantStatus: 400,
			WantRes:    `{"status":"Error","error":"id must be int"}`,
		},
		{
			Name:       "invalid json id",
			Url:        "/songs",
			Body:       `{"id":"zero"}`,
			WantStatus: 400,
			WantRes:    `{"status":"Error","error":"Can't decode json body"}`,
		},
	}

	handler := NewHandler(log, mock)

	for _, testCase := range testCases {
		testCase.Method = "PUT"

		test.TestEndpoint(t, http.HandlerFunc(handler.Update), testCase)
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
		Url:        "/songs?id=1",
		Method:     "PUT",
		Body:       `{}`,
		WantStatus: 400,
		WantRes:    `{"status":"Error","error":"Song not exists"}`,
	}

	handler := NewHandler(log, mock)

	test.TestEndpoint(t, http.HandlerFunc(handler.Update), testCase)
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
			Url:        "/songs?id=1",
			Body:       `{}`,
			WantStatus: 204,
			WantRes:    `{"status":"Ok"}`,
		},
		{
			Name:       "invalid id",
			Url:        "/songs?id=one",
			Body:       `{}`,
			WantStatus: 400,
			WantRes:    `{"status":"Error","error":"id must be int"}`,
		},
	}

	handler := NewHandler(log, mock)

	for _, testCase := range testCases {
		testCase.Method = "DELETE"

		test.TestEndpoint(t, http.HandlerFunc(handler.Delete), testCase)
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
		Url:        "/songs?id=1",
		Method:     "DELETE",
		Body:       `{}`,
		WantStatus: 400,
		WantRes:    `{"status":"Error","error":"Song not exists"}`,
	}

	handler := NewHandler(log, mock)

	test.TestEndpoint(t, http.HandlerFunc(handler.Delete), testCase)
}
