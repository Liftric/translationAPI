package routing

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"preventis.io/translationApi/model"
)

func getAllActiveProjects(c *gin.Context) {
	var projects []model.Project
	db.Find(&projects)
	db.Where("archived = ?", false).Find(&projects)
	c.JSON(200, projects)
}

func getAllArchivedProjects(c *gin.Context) {
	var projects []model.Project
	db.Find(&projects)
	db.Where("archived = ?", true).Find(&projects)
	c.JSON(200, projects)
}

func getProject(c *gin.Context) {
	id := c.Param("id")
	var project model.Project
	if err := db.Where("id = ?", id).First(&project).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, project)
	}
}

type ProjectValidation struct {
	Name    string `form:"name" json:"name" xml:"name"  binding:"required"`
	IsoCode string `form:"baseLanguageCode" json:"baseLanguageCode" xml:"baseLanguageCode"  binding:"required"`
}

func createProject(c *gin.Context) {
	var project model.Project
	var json ProjectValidation
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	project.Name = json.Name
	var baseLang model.Language
	if err := db.Where("iso_code = ?", json.IsoCode).First(&baseLang).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	}

	if dbc := db.Create(&project); dbc.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": dbc.Error.Error()})
		return
	}
	c.JSON(200, project)
}

type ProjectRenameValidation struct {
	Id   uint   `form:"id" json:"id" xml:"id"  binding:"required"`
	Name string `form:"name" json:"name" xml:"name"  binding:"required"`
}

func renameProject(c *gin.Context) {
	var json ProjectRenameValidation
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var project model.Project
	if err := db.Where("id = ?", id).First(&project).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		project.Name = json.Name
		db.Save(&project)
		c.JSON(200, project)
	}
}

type ProjectArchiveValidation struct {
	Id      uint `form:"id" json:"id" xml:"id"  binding:"required"`
	Archive bool `form:"archive" json:"archive" xml:"archive"  binding:"required"`
}

func archiveProject(c *gin.Context) {
	var json ProjectArchiveValidation
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var project model.Project
	if err := db.Where("id = ?", id).First(&project).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		project.Archived = json.Archive
		db.Save(&project)
		c.JSON(200, project)
	}
}
