package postgres

import (
	"context"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/s3nn1k/ef-mob-task/internal/models"
)

func TestCreate(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}

	song := models.Song{
		Song:  "TestSong",
		Group: "TestGroup",
		Text:  "TestText",
		Link:  "TestLink",
		Date:  "TestDate",
	}

	id := 1

	mock.ExpectQuery("^INSERT INTO songs (.+) RETURNING id$").
		WithArgs(song.Song, song.Group, song.Text, song.Link, song.Date).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(id))

	db := NewStorage(mock)

	storedId, err := db.Create(context.Background(), song)
	if err != nil {
		t.Fatalf("error not expected while creating: %s", err)
	}

	if storedId != id {
		t.Fatal("error: id must be returned after creating")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %s", err)
	}
}
