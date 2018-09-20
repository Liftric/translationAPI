// database package that contains models for the gorm library. Database connection can be acquired here and the database
// gets initialized here.
package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Connects to database, initializes the models in the database and returns a gorm database object
func StartDB(dbDialect string, dbArgs string) *gorm.DB {
	db, err := gorm.Open(dbDialect, dbArgs)
	if err != nil {
		panic(fmt.Sprintf("failed to connect database: %s", err.Error()))
	}
	InitDB(db)
	return db
}

// Registers the models in the database
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

// Languages the strings can be translated in
type Language struct {
	IsoCode string `gorm:"primary_key"`
	Name    string `gorm:"unique;not null"`
}

// Identifies the string, contains the different translations and corresponds to a project
type StringIdentifier struct {
	gorm.Model
	Identifier   string        `gorm:"unique_index:idx_identifier_project"`
	Project      Project       `gorm:"foreignkey:ProjectID"`
	ProjectID    uint          `gorm:"not null;unique_index:idx_identifier_project"`
	Translations []Translation `gorm:"foreignkey:StringIdentifierID"`
}

// Translation of a string in one language
type Translation struct {
	gorm.Model
	Translation        string
	Identifier         StringIdentifier `gorm:"foreignkey:StringIdentifierID"`
	StringIdentifierID uint             `gorm:"not null;unique_index:idx_identifier_language"`
	Language           Language         `gorm:"foreignkey:LanguageRefer"`
	LanguageRefer      string           `gorm:"not null;unique_index:idx_identifier_language"`
	Approved           bool             `gorm:"default:'false'"`
}
