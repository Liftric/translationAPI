package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func StartDB(dbDialect string, dbArgs string) *gorm.DB {
	db, err := gorm.Open(dbDialect, dbArgs)
	if err != nil {
		panic(fmt.Sprintf("failed to connect database: %s", err.Error()))
	}
	InitDB(db)
	return db
}

func InitDB(db *gorm.DB) {
	db.AutoMigrate(&Language{})
	db.AutoMigrate(&Project{})
	db.AutoMigrate(&StringKey{})
	db.AutoMigrate(&Translation{})
}

type Project struct {
	gorm.Model
	Name              string   `gorm:"unique;not null"`
	BaseLanguage      Language `gorm:"foreignkey:BaseLanguageRefer"`
	BaseLanguageRefer string
	Languages         []Language `gorm:"many2many:project_languages;"`
	Archived          bool
}

type Language struct {
	IsoCode string `gorm:"primary_key"`
	Name    string `gorm:"unique;not null"`
}

type StringKey struct {
	gorm.Model
	Key          string
	Project      Project `gorm:"foreignkey:ProjectRefer"`
	ProjectRefer uint
}

type Translation struct {
	gorm.Model
	Translation   string
	Key           StringKey `gorm:"foreignkey:KeyRefer"`
	KeyRefer      uint
	Language      Language `gorm:"foreignkey:LanguageRefer"`
	LanguageRefer string
	Approved      bool `gorm:"default:'false'"`
}
