package postgres

import (
	"context"
	"testing"
	"time"

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

func TestGetAll(t *testing.T) {
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
		Date:  time.Now().Format("02.01.2006"),
	}

	filters := models.GetFilters{
		Limit:  1,
		Offset: 1,
		Id:     1,
		Song:   song.Song,
		Group:  song.Group,
		Date:   song.Date,
	}

	mock.ExpectQuery("^SELECT (.+) FROM songs WHERE (.+)$").
		WithArgs(filters.Id, filters.Song, filters.Group, filters.Date, filters.Limit, filters.Offset).
		WillReturnRows(pgxmock.NewRows([]string{"id", "song", "group", "text", "link", "date"}).
			AddRow(song.Id, song.Song, song.Group, song.Text, song.Link, song.Date).
			AddRow(song.Id, song.Song, song.Group, song.Text, song.Link, song.Date))

	db := NewStorage(mock)

	songs, err := db.GetAll(context.Background(), filters)
	if err != nil {
		t.Fatalf("error not expected while get all songs: %s", err)
	}

	if len(songs) != 2 {
		t.Fatal("error: must get same songs as in storage")
	}

	for _, storedSong := range songs {
		if storedSong.Id != song.Id {
			t.Fatalf("error: returned id must be the same")
		}

		if storedSong.Song != song.Song {
			t.Fatalf("error: returned song must be the same as in storage")
		}

		if storedSong.Group != song.Group {
			t.Fatalf("error: returned group must be the same as in storage")
		}

		if storedSong.Text != song.Text {
			t.Fatalf("error: returned text must be the same as in storage")
		}

		if storedSong.Link != song.Link {
			t.Fatalf("error: returned link must be the same as in storage")
		}

		if storedSong.Date != song.Date {
			t.Fatalf("error: returned date must be the same as in storage")
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %s", err)
	}
}
