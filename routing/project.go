package routing

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"preventis.io/translationApi/model"
)

type simpleProjectDTO struct {
	Id           uint
	Name         string
	BaseLanguage model.Language
	Languages    []model.Language
}

func getAllActiveProjects(c *gin.Context) {
	var projects []model.Project
	db.Where("archived = ?", false).Preload("BaseLanguage").Preload("Languages").Find(&projects)
	result := convertProjectListToDTO(projects)
	c.JSON(200, result)
}

func getAllArchivedProjects(c *gin.Context) {
	var projects []model.Project
	db.Where("archived = ?", true).Find(&projects)
	result := convertProjectListToDTO(projects)
	c.JSON(200, result)
}

func convertProjectListToDTO(projects []model.Project) []simpleProjectDTO {
	var result []simpleProjectDTO
	result = []simpleProjectDTO{}
	for _, e := range projects {
		var languages = e.Languages
		if languages == nil {
			languages = []model.Language{}
		}
		p := simpleProjectDTO{e.ID, e.Name, e.BaseLanguage, languages}
		result = append(result, p)
	}
	return result
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

type projectValidation struct {
	Name    string `form:"name" json:"name" xml:"name"  binding:"required"`
	IsoCode string `form:"baseLanguageCode" json:"baseLanguageCode" xml:"baseLanguageCode"  binding:"required"`
}

func createProject(c *gin.Context) {
	var json projectValidation
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var project model.Project
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
	c.JSON(201, project)
}

type projectRenameValidation struct {
	Name string `form:"name" json:"name" xml:"name"  binding:"required"`
}

func renameProject(c *gin.Context) {
	id := c.Param("id")
	var json projectRenameValidation
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

type projectArchiveValidation struct {
	Archive bool `form:"archive" json:"archive" xml:"archive"  binding:"required"`
}

func archiveProject(c *gin.Context) {
	id := c.Param("id")
	var json projectArchiveValidation
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

type projectLanguageValidation struct {
	IsoCode string `form:"languageCode" json:"languageCode" xml:"languageCode"  binding:"required"`
}

func addLanguageToProject(c *gin.Context) {
	id := c.Param("id")
	var json projectLanguageValidation
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var project model.Project
	if err := db.Where("id = ?", id).First(&project).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		var lang model.Language
		if err := db.Where("iso_code = ?", json.IsoCode).First(&lang).Error; err != nil {
			c.AbortWithStatus(404)
			fmt.Println(err)
		}
		db.Model(&project).Association("Languages").Append(lang)
		db.Save(&project)
		c.JSON(200, project)
	}
}

func setBaseLanguage(c *gin.Context) {
	id := c.Param("id")
	var json projectLanguageValidation
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var project model.Project
	if err := db.Where("id = ?", id).First(&project).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		var baseLang model.Language
		if err := db.Where("iso_code = ?", json.IsoCode).First(&baseLang).Error; err != nil {
			c.AbortWithStatus(404)
			fmt.Println(err)
		}
		project.BaseLanguage = baseLang
		db.Save(&project)
		c.JSON(200, project)
	}
}
