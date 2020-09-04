package config

import "os"

const (
	dbName     = "DB_NAME"
	dbPassword = "DB_PASS"
	dbUser     = "DB_USER"
	dbHost     = "DB_HOST"
	dbPort     = "DB_PORT"
)

var (
	databaseName     = os.Getenv(dbName)
	databaseUser     = os.Getenv(dbUser)
	databasePassword = os.Getenv(dbPassword)
	databaseHost     = os.Getenv(dbHost)
	databasePort     = os.Getenv(dbPort)
)

func GetDatabaseHost() string {
	return databaseHost
}

func GetDatabaseName() string {
	return databaseName
}

func GetDatabaseUser() string {
	return databaseUser
}

func GetDatabasePassword() string {
	return databasePassword
}

func GetDatabasePort() string {
	return databasePort
}
