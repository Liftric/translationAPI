package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"os"
	"preventis.io/translationApi/model"
	"preventis.io/translationApi/routing"
)

func main() {
	dbDialect := getEnv("DATABASE_TYPE", "mysql")
	dbArgs := getDbArgs(dbDialect)

	db, err := gorm.Open(dbDialect, dbArgs)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	model.InitDB(db)
	routing.StartRouter(db)
}

func getDbArgs(dbDialect string) string {
	dbHost := getEnv("DATABASE_HOST", "localhost")
	dbPort := getEnv("DATABASE_PORT", "3306")
	dbUser := getEnv("DATABASE_USER", "translation")
	dbPassword := getEnv("DATABASE_PASSWORD", "translation")
	dbDatabase := getEnv("DATABASE_NAME", "translation")

	if dbDialect == "mysql" {
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbPort, dbDatabase)
	} else if dbDialect == "postgres" {
		return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s", dbHost, dbPort, dbUser, dbDatabase, dbPassword)
	}
	return ""
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
