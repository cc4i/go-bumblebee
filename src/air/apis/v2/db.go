package v2

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

func sample() {
	db, err := sql.Open("mysql", "user:password@/dbname")
	if err != nil {
		log.Errorf("%s\n", err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Errorf("%s\n", err)
	}
}

func LogCall() {

}
