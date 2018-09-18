package repository

import (
	"log"
	"net/url"
	"time"

	"github.com/frankh/norbert/cmd/norbert/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type sqlRepository struct {
	*sqlx.DB
}

type Repository interface {
	SaveCheckResult(*models.CheckResult) error
	CheckResults(string) ([]*models.CheckResult, error)
}

var schema = `
CREATE TABLE IF NOT EXISTS check_results (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  checkid TEXT,
  starttime TIMESTAMP,
  endtime TIMESTAMP,
  resultcode INT,
  errormsg TEXT
);

CREATE INDEX IF NOT EXISTS check_results_starttime ON check_results(starttime);
`

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
	_, err = db.Exec(schema)
	if err != nil {
		return nil, err
	}

	return &sqlRepository{db}, nil
}

func (db *sqlRepository) SaveCheckResult(result *models.CheckResult) error {
	_, err := db.NamedExec(`INSERT INTO check_results
      (checkid, starttime, endtime, resultcode, errormsg)
    VALUES (:checkid, :starttime, :endtime, :resultcode, :errormsg)
  `, result)
	return err
}

func (db *sqlRepository) CheckResults(checkId string) ([]*models.CheckResult, error) {
	results := []*models.CheckResult{}
	err := db.Select(&results, "SELECT * FROM check_results WHERE checkId=$1 ORDER BY starttime DESC LIMIT 100", checkId)
	return results, err
}
