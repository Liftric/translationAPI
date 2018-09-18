package model

import "github.com/jinzhu/gorm"

func InitDB(db *gorm.DB) {
	db.AutoMigrate(&Language{})
	db.AutoMigrate(&Project{})
}

type Project struct {
	gorm.Model
	Name              string
	BaseLanguage      Language `gorm:"foreignkey:BaseLanguageRefer"`
	BaseLanguageRefer uint
	Languages         []Language `gorm:"many2many:project_languages;"`
}

type Language struct {
	gorm.Model
	Name string
}

type StringKey struct {
	gorm.Model
	Key           string
	Language      Language
	LanguageRefer uint
	Project       Project
	ProjectRefer  uint
}

type Translation struct {
	gorm.Model
	Translation string
	Key         StringKey
	KeyRefer    uint
}
