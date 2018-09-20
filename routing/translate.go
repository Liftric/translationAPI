package routing

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"preventis.io/translationApi/model"
)

type stringIdentifierValidation struct {
	ProjectId  int    `form:"projectId" json:"projectId" xml:"projectId"  binding:"required"`
	Identifier string `form:"identifier" json:"identifier" xml:"identifier"  binding:"required"`
}

func createIdentifier(c *gin.Context) {
	var json stringIdentifierValidation
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var key model.StringIdentifier
	key.Identifier = json.Identifier
	var project model.Project
	if err := db.Where("id = ?", json.ProjectId).First(&project).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
		return
	}
	key.Project = project

	var existingKeys []model.StringIdentifier
	db.Where("Identifier = ? AND project_id = ?", "newIdentifier", 1).Find(&existingKeys)
	if len(existingKeys) > 0 {
		c.Status(409)
		return
	}

	if dbc := db.Create(&key); dbc.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": dbc.Error.Error()})
		return
	}
	c.Status(201)
}

type updateIdentifierValidation struct {
	Key string `form:"identifier" json:"identifier" xml:"identifier"  binding:"required"`
}

func updateIdentifier(c *gin.Context) {
	id := c.Param("id")
	var json updateIdentifierValidation
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var key model.StringIdentifier
	if err := db.Where("id = ?", id).First(&key).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
		return
	}
	key.Identifier = json.Key
	db.Save(&key)
	c.Status(200)
}

type translationValidation struct {
	KeyId       int    `form:"keyId" json:"keyId" xml:"keyId"  binding:"required"`
	Translation string `form:"translation" json:"translation" xml:"translation"  binding:"required"`
	Language    string `form:"languageCode" json:"languageCode" xml:"languageCode"  binding:"required"`
}

func createTranslation(c *gin.Context) {
	var json translationValidation
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var translation model.Translation
	translation.Translation = json.Translation
	var key model.StringIdentifier
	if err := db.Where("id = ?", json.KeyId).First(&key).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
		return
	}
	translation.Identifier = key

	var lang model.Language
	if err := db.Where("iso_code = ?", json.Language).First(&lang).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
		return
	}
	translation.Language = lang

	if dbc := db.Create(&translation); dbc.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": dbc.Error.Error()})
		return
	}
	c.Status(201)
}

type updateTranslationValidation struct {
	Translation string `form:"translation" json:"translation" xml:"translation"  binding:"required"`
}

func updateTranslation(c *gin.Context) {
	id := c.Param("id")
	var json updateTranslationValidation
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var translation model.Translation
	if err := db.Where("id = ?", id).First(&translation).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
		return
	}
	translation.Translation = json.Translation
	translation.Approved = false
	db.Save(&translation)
	c.Status(200)
}
func setApproved(c *gin.Context) {
	id := c.Param("id")
	var translation model.Translation
	if err := db.Where("id = ?", id).First(&translation).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
		return
	} else {
		translation.Approved = true
		db.Save(&translation)
		c.Status(200)
	}
}

func moveKey(c *gin.Context) {
	id := c.Param("id")
	projectId := c.Param("projectId")

	var project model.Project
	if err := db.Where("id = ?", projectId).First(&project).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
		return
	}

	var key model.StringIdentifier
	if err := db.Where("id = ?", id).First(&key).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
		return
	}
	key.Project = project
	db.Save(&key)
	c.Status(200)
}

func deleteKey(c *gin.Context) {
	id := c.Param("id")
	var key model.StringIdentifier
	if err := db.Where("id = ?", id).First(&key).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
		return
	}
	db.Delete(&key)
	c.Status(200)
}
