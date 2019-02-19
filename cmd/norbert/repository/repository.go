package repository

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/frankh/norbert/cmd/norbert/config"
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
	GetService(string) (*models.Service, error)
}

var schema = `
CREATE TABLE IF NOT EXISTS check_results (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  checkid TEXT,
  starttime TIMESTAMP,
  endtime TIMESTAMP,
  resultcode TEXT,
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
		log.Printf("Unable to connect to database (%s), retrying (%d attempts left)\n", err, tries)
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
	rows, err := db.NamedQuery(`INSERT INTO check_results
      (checkid, starttime, endtime, resultcode, errormsg)
    VALUES (:checkid, :starttime, :endtime, :resultcode, :errormsg)
    RETURNING id
  `, result)
	if err != nil {
		return err
	}
	rows.Next()
	cols, err := rows.SliceScan()
	if err != nil || len(cols) == 0 {
		return err
	}
	insertedId, ok := cols[0].([]byte)
	if !ok {
		return fmt.Errorf("Couldn't get inserted it")
	}
	result.Id = string(insertedId)
	return nil
}

func (db *sqlRepository) CheckResults(checkId string) ([]*models.CheckResult, error) {
	results := []*models.CheckResult{}
	err := db.Select(&results, "SELECT * FROM check_results WHERE checkId=$1 ORDER BY starttime DESC LIMIT 100", checkId)
	return results, err
}

func (db *sqlRepository) GetService(serviceName string) (*models.Service, error) {
	return config.Services[serviceName], nil
}
