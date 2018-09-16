package leader

import (
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
)

type leaderElector struct {
	id     string
	table  string
	config Config
	db     *sql.DB

	currentLease *lease
}

type lease struct {
	id        string
	leaseEnd  time.Time
	dbTime    time.Time
	localTime time.Time
}

type LeaderElector interface {
	// IsRunning() bool
	IsLeader() bool
	Start()
}

type Config struct {
	LeaseDuration time.Duration
}

var defaultConfig Config = Config{
	LeaseDuration: 10 * time.Second,
}

func NewElector(table string, db *sql.DB) (LeaderElector, error) {
	_, err := db.Exec(`
    CREATE TABLE IF NOT EXISTS "` + table + `" (
      singleton TEXT UNIQUE,
      id TEXT,
      lease_end TIMESTAMP
    )
  `)

	if err != nil {
		log.Println("Failed to create lease table: ", err)
		return nil, err
	}

	return &leaderElector{
		id:     uuid.New().String(),
		table:  table,
		config: defaultConfig,
		db:     db,
	}, nil
}

func (e *leaderElector) Start() {
	go func() {
		e.elect()
		// Make the ticker interval 50% of the lease duration
		interval := e.config.LeaseDuration / time.Duration(2)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				e.elect()
			}
		}
	}()
}

func (e *leaderElector) IsLeader() bool {
	l := e.currentLease
	if l == nil || l.id != e.id {
		return false
	}

	elapsed := time.Since(l.localTime)
	// We are the leader if the current lease matches our id, and
	// the lease hasn't expired.
	// We always use the database time as canonical, but we trust
	// our local time to measure the time since we got the lease.
	return l.dbTime.Add(elapsed).Before(l.leaseEnd)
}

func (e *leaderElector) elect() error {
	tx, err := e.db.Begin()
	if err != nil {
		log.Println("Error starting transaction: ", err)
		return err
	}
	defer tx.Rollback()

	var l lease

	row := tx.QueryRow(`SELECT id, lease_end, CURRENT_TIMESTAMP FROM "` + e.table + `" WHERE singleton='leader'`)
	err = row.Scan(&l.id, &l.leaseEnd, &l.dbTime)
	if err != nil && err != sql.ErrNoRows {
		log.Println("Error selecting current lease: ", err)
		return err
	}
	l.localTime = time.Now()

	if err == sql.ErrNoRows || l.id == e.id || l.dbTime.After(l.leaseEnd) {
		_, err := tx.Exec(`
      INSERT INTO "`+e.table+`" (singleton, id, lease_end)
      VALUES ('leader', $1, $2)
      ON CONFLICT(singleton) DO UPDATE SET id=EXCLUDED.id, lease_end=EXCLUDED.lease_end
    `, e.id, l.dbTime.Add(e.config.LeaseDuration))

		if err != nil {
			log.Println("Error updating lease: ", err)
			return err
		}
	}

	e.currentLease = &l
	tx.Commit()
	return nil
}
