package assignment2

import "fmt"

// Cfg : Configuration file for MongoDB credentials
var Cfg = newConfig("imt2681_assignment2", "guest", "imt2681", "ds227035", "27035")

type Configuration struct {
	DBName         string
	DBUser, DBPass string
	DBurl          string
}

func newConfig(dbName, dbUser, dbPass, id, port string) (conf Configuration) {
	conf.DBName = dbName
	conf.DBUser = dbUser
	conf.DBPass = dbPass
	conf.DBurl = fmt.Sprintf("mongodb://%s:%s@%s.mlab.com:%s/%s", dbUser, dbPass, id, port, dbName)
	return
}
