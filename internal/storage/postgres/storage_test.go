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

func TestUpdate(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}

	song := models.Song{
		Id:    1,
		Song:  "TestSong",
		Group: "TestGroup",
		Text:  "TestText",
		Link:  "TestLink",
		Date:  "TestDate",
	}

	mock.ExpectExec("^UPDATE songs SET (.+) WHERE (.+)$").
		WithArgs(song.Song, song.Group, song.Text, song.Link, song.Date, song.Id).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	db := NewStorage(mock)

	ok, err := db.Update(context.Background(), song)
	if err != nil {
		t.Fatalf("error not expected while updating: %s", err)
	}

	if !ok {
		t.Fatal("error: result of updating must be true")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %s", err)
	}
}

func TestDelete(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}

	id := 1

	mock.ExpectExec("^DELETE FROM songs WHERE (.+)$").
		WithArgs(id).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	db := NewStorage(mock)

	ok, err := db.Delete(context.Background(), id)
	if err != nil {
		t.Fatalf("error not expected while deleting: %s", err)
	}

	if !ok {
		t.Fatal("error: result of deleting must be true")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %s", err)
	}
}

func TestGet(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}

	song := models.Song{
		Id:    1,
		Song:  "TestSong",
		Group: "TestGroup",
		Text:  "TestText",
		Link:  "TestLink",
		Date:  "TestDate",
	}

	limit := 1
	offset := 1

	mock.ExpectQuery("^SELECT (.+) FROM songs WHERE (.+)$").
		WithArgs(song.Id, song.Song, song.Group, song.Text, song.Link, song.Date, limit, offset).
		WillReturnRows(pgxmock.NewRows([]string{"id", "song", "group", "text", "link", "date"}).
			AddRow(song.Id, song.Song, song.Group, song.Text, song.Link, song.Date))

	db := NewStorage(mock)

	songs, err := db.Get(context.Background(), song, limit, offset)
	if err != nil {
		t.Fatalf("error not expected while updating: %s", err)
	}

	if len(songs) != 1 {
		t.Fatal("error: result of updating must be true")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %s", err)
	}
}
