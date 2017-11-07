package main

import (
	"fmt"
)

// Fill out with your own values
var cfg = newConfig("guest", "imt2681", "ds227035", "27035", "imt2681_assignment2")

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
