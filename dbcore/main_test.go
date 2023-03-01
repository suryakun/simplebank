package dbcore

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5555/simple_bank?sslmode=disable"
)

var testQUeries *Queries
var dbConn *sql.DB

func TestMain(m *testing.M) {
	var err error
	dbConn, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect db:", err)
	}
	testQUeries = New(dbConn)

	os.Exit(m.Run())
}
