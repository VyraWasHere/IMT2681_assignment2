package main

import (
	"fmt"
	"os"
)

// Fills out with values from env
var cfg = newConfig(os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_ID"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

type config struct {
	DBName,
	DBUser,
	DBPass,
	DBuri string
}

func newConfig(dbuser, dbpass, mLabID, port, dbname string) (cfg config) {
	cfg.DBName = dbname
	cfg.DBUser = dbuser
	cfg.DBPass = dbpass
	cfg.DBuri = fmt.Sprintf("mongodb://%s:%s@%s.mlab.com:%s/%s", dbuser, dbpass, mLabID, port, dbname)
	return
}
