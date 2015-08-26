package session

import (
	"database/sql"
	log "github.com/Sirupsen/logrus"
	"github.com/coopernurse/gorp"
	_ "github.com/lib/pq"
	stdlog "log"
	"os"
        "fmt"
)

var dbmap *gorp.DbMap
var mydb *sql.DB

func init() {
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Panic("Missing required environment variable 'DATABASE_URL'.")
	}

	db, err := sql.Open("postgres", dbUrl)
fmt.Println("dbUrl", dbUrl)
        checkErr(err)
	if nil != err {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Database connection error")
                log.Panic("failed to open db")
	}

	dbmap = &gorp.DbMap{
		Db:      db,
		Dialect: gorp.PostgresDialect{},
	}

	if os.Getenv("DEBUG") == "true" {
		dbmap.TraceOn("[gorp]", stdlog.New(os.Stdout, "authentication:", stdlog.Lmicroseconds))
	}

	t := dbmap.AddTableWithName(User{}, "users").SetKeys(true, "Id")
        log.Info(t)
	t = dbmap.AddTableWithName(Event{}, "events").SetKeys(true, "Id")
        log.Info(t)
	t = dbmap.AddTableWithName(Attendance{}, "attendance")
        log.Info(t)
   
        mydb = db

}

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}
