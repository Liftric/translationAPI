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
	db.AutoMigrate(&StringIdentifier{})
	db.AutoMigrate(&Translation{})
}

type Project struct {
	gorm.Model
	Name              string             `gorm:"unique;not null"`
	BaseLanguage      Language           `gorm:"foreignkey:BaseLanguageRefer"`
	BaseLanguageRefer string             `gorm:"not null"`
	Languages         []Language         `gorm:"many2many:project_languages;"`
	Archived          bool               `gorm:"default:'false'"`
	Identifiers       []StringIdentifier `gorm:"foreignkey:ProjectID"`
}

type Language struct {
	IsoCode string `gorm:"primary_key"`
	Name    string `gorm:"unique;not null"`
}

type StringIdentifier struct {
	gorm.Model
	Identifier   string        `gorm:"unique_index:idx_identifier_project"`
	Project      Project       `gorm:"foreignkey:ProjectID"`
	ProjectID    uint          `gorm:"not null;unique_index:idx_identifier_project"`
	Translations []Translation `gorm:"foreignkey:StringIdentifierID"`
}

type Translation struct {
	gorm.Model
	Translation        string
	Identifier         StringIdentifier `gorm:"foreignkey:StringIdentifierID"`
	StringIdentifierID uint             `gorm:"not null"`
	Language           Language         `gorm:"foreignkey:LanguageRefer"`
	LanguageRefer      string           `gorm:"not null"`
	Approved           bool             `gorm:"default:'false'"`
}
