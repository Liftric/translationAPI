package main

import (
	"github.com/jinzhu/gorm"
)
import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
)


func main() {
	db, err := gorm.Open("mysql", "translation:translation@/translation?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	initDB(db)
	startRouter()
}



