package main

import (
	"fmt"

	"os"
	"preventis.io/translationApi/model"
	"preventis.io/translationApi/routing"
)

func main() {
	dbDialect := getEnv("DATABASE_TYPE", "sqlite3")
	dbArgs := getDbArgs(dbDialect)

	db := model.StartDB(dbDialect, dbArgs)
	defer db.Close()
	routing.StartRouter(db)
}

func getDbArgs(dbDialect string) string {
	dbHost := getEnv("DATABASE_HOST", "")
	dbPort := getEnv("DATABASE_PORT", "")
	dbUser := getEnv("DATABASE_USER", "")
	dbPassword := getEnv("DATABASE_PASSWORD", "")
	dbDatabase := getEnv("DATABASE_NAME", "/tmp/gorm.db")
	dbSSL := getEnv("DATABASE_SSL", "")

	if dbDialect == "sqlite3" {
		return dbDatabase
	} else if dbDialect == "mysql" {
		if dbSSL == "" {
			dbSSL = "false"
		}
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local&tls=%s", dbUser, dbPassword, dbHost, dbPort, dbDatabase, dbSSL)
	} else if dbDialect == "postgres" {
		if dbSSL == "" {
			dbSSL = "disabled"
		}
		return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", dbHost, dbPort, dbUser, dbDatabase, dbPassword, dbSSL)
	}
	return ""
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
