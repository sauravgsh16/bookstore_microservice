package usersdb

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	// postgres driver
	_ "github.com/lib/pq"
)

const (
	dBHost = "localhost"
	dBPort = 5432
	dBUser = "postgres"
	dBPwd  = "postgres"
	dBName = "bookstore"
)

// DB connection
var DB dbConn

type dbConn struct {
	conn *sql.DB
}

func init() {
	connInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dBHost,
		dBPort,
		dBUser,
		dBPwd,
		dBName,
	)
	var err error
	DB = dbConn{}
	DB.conn, err = sql.Open("postgres", connInfo)
	if err != nil {
		panic(err)
	}

	if err = DB.conn.Ping(); err != nil {
		panic(err)
	}
	log.Println("Successfully configured database")
}

func (db *dbConn) GetConn(ctx context.Context) *sql.Conn {
	conn, err := db.conn.Conn(ctx)
	if err != nil {
		panic(err.Error())
	}
	return conn
}
