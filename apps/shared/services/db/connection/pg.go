package connection

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"

	"github.com/ilya-mezentsev/micro-dep/shared/types/configs"
)

func MustGetConnection(c configs.DB) *sqlx.DB {
	var (
		db  *sqlx.DB
		err error
	)
	tryNumber := 1
	for {
		db, err = sqlx.Open("postgres", BuildPostgresString(c))
		if err != nil {
			log.Errorf("Unable to open DB connection: %v. try number #%d", err, tryNumber)
			time.Sleep(time.Second * time.Duration(c.Connection.RetryTimeout))
		} else if err = db.Ping(); err != nil {
			log.Errorf("Unable to ping DB: %v. try number #%d", err, tryNumber)
			time.Sleep(time.Second * time.Duration(c.Connection.RetryTimeout))
		} else {
			break
		}

		tryNumber++
		if tryNumber > c.Connection.RetryCount {
			break
		}
	}

	if err != nil {
		log.Fatalf("Unable to create DB connection: %v", err)
	}

	return db
}

func BuildPostgresString(c configs.DB) string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host,
		c.Port,
		c.User,
		c.Password,
		c.DBName,
	)
}
