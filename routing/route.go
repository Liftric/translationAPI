package routing

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

func StartRouter(database *gorm.DB) {
	db = database

	r := gin.Default()
	// display all projects (with statistics)
	r.GET("/projects", getAllActiveProjects)
	// display all archived projects
	r.GET("/projects/archived", getAllArchivedProjects)
	// display all strings, statistics and config of a project
	r.GET("/project/:id", getProject)
	// create project
	r.PUT("/project", createProject)
	// create language
	r.PUT("/language", createLanguage)
	// create key in project
	r.PUT("/project/:id/key", createKey)
	// change key
	r.POST("/project/:id/key", updateKey)
	// change translation
	r.POST("/project/:id/translation", updateTranslation)
	// set revised for translation in a language
	r.POST("/project/:id/revised/:key", setRevised)
	// change project name
	r.POST("/projectName", renameProject)
	// move key to another project
	r.POST("/project/:id/moveKey", moveKey)
	// add language to project
	r.PUT("/project/:id/languages", addLanguage)
	// set base language of project
	r.POST("/project/:id/baseLanguage", setBaseLanguage)
	// archive project
	r.DELETE("/project/:id", archiveProject)
	// delete key
	r.DELETE("/project/:id/:key", deleteKey)
	// diff iOS strings file and db
	r.POST("/project/:id/ios", diffIOS)
	// diff android strings file and db
	r.POST("/project/:id/android", diffAndroid)
	// diff excel file and db
	r.POST("/project/:id/excel", diffExcel)
	// export ios strings
	r.GET("/project/:id/ios", exportIOS)
	// export android strings
	r.GET("/project/:id/android", exportAndroid)
	// export to excel
	r.GET("/project/:id/excel", exportExcel)

	r.Run() // listen and serve on 0.0.0.0:8080
}

func createLanguage(c *gin.Context) {
	// TODO
}
func createKey(c *gin.Context) {
	// TODO
}
func updateKey(c *gin.Context) {
	// TODO
}
func updateTranslation(c *gin.Context) {
	// TODO
}
func setRevised(c *gin.Context) {
	// TODO
}

func moveKey(c *gin.Context) {
	// TODO
}
func addLanguage(c *gin.Context) {
	// TODO
}
func setBaseLanguage(c *gin.Context) {
	// TODO
}

func deleteKey(c *gin.Context) {
	// TODO
}
func diffIOS(c *gin.Context) {
	// TODO
}
func diffAndroid(c *gin.Context) {
	// TODO
}
func diffExcel(c *gin.Context) {
	// TODO
}
func exportIOS(c *gin.Context) {
	// TODO
}
func exportAndroid(c *gin.Context) {
	// TODO
}
func exportExcel(c *gin.Context) {
	// TODO
}
