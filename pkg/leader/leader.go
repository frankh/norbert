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
	Table         string
	LeaseDuration time.Duration
}

var DefaultConfig Config = Config{
	Table:         "leases",
	LeaseDuration: 10 * time.Second,
}

func NewElector(db *sql.DB, config Config) (LeaderElector, error) {
	_, err := db.Exec(`
    CREATE TABLE IF NOT EXISTS "` + config.Table + `" (
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
		config: config,
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
				wasLeader := e.IsLeader()
				e.elect()
				if !wasLeader && e.IsLeader() {
					log.Println("Elected leader")
				}
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
	e.db.Ping()
	tx, err := e.db.Begin()
	if err != nil {
		log.Println("Error starting transaction: ", err)
		return err
	}
	defer tx.Rollback()

	var l lease

	row := tx.QueryRow(`SELECT id, lease_end, CURRENT_TIMESTAMP FROM "` + e.config.Table + `" WHERE singleton='leader'`)
	err = row.Scan(&l.id, &l.leaseEnd, &l.dbTime)
	if err != nil && err != sql.ErrNoRows {
		log.Println("Error selecting current lease: ", err)
		return err
	}
	l.localTime = time.Now()

	if err == sql.ErrNoRows || l.id == e.id || l.dbTime.After(l.leaseEnd) {
		_, err := tx.Exec(`
      INSERT INTO "`+e.config.Table+`" (singleton, id, lease_end)
      VALUES ('leader', $1, $2)
      ON CONFLICT(singleton) DO UPDATE SET id=EXCLUDED.id, lease_end=EXCLUDED.lease_end
    `, e.id, l.dbTime.Add(e.config.LeaseDuration))

		if err != nil {
			log.Println("Error updating lease: ", err)
			return err
		}

		// Re-fetch updated lease
		row := tx.QueryRow(`SELECT id, lease_end, CURRENT_TIMESTAMP FROM "` + e.config.Table + `" WHERE singleton='leader'`)
		row.Scan(&l.id, &l.leaseEnd, &l.dbTime)
	}

	e.currentLease = &l
	tx.Commit()
	return nil
}
