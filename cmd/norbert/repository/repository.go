package repository

import (
	"log"
	"net/url"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type sqlRepository struct {
	DB *sqlx.DB
}

func NewRepository(dbURI string) (*sqlRepository, error) {
	url, err := url.Parse(dbURI)
	if err != nil {
		return nil, err
	}

	var db *sqlx.DB
	for tries := 10; tries > 0; tries -= 1 {
		db, err = sqlx.Connect(url.Scheme, url.String())
		if err == nil {
			break
		}
		log.Printf("Unable to connect to database, retrying (%d attempts left)\n", tries)
		time.Sleep(500 * time.Millisecond)
	}
	if err != nil {
		return nil, err
	}

	return &sqlRepository{db}, nil
}
