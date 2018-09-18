package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"preventis.io/translationApi/model"
	"preventis.io/translationApi/routing"
)

func main() {
	db, err := gorm.Open("mysql", "translation:translation@/translation?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	model.InitDB(db)
	routing.StartRouter()
}
