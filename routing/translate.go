package routing

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"preventis.io/translationApi/model"
)

type stringKeyValidation struct {
	ProjectId int    `form:"projectId" json:"projectId" xml:"projectId"  binding:"required"`
	Key       string `form:"key" json:"key" xml:"key"  binding:"required"`
}

func createKey(c *gin.Context) {
	var json stringKeyValidation
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var key model.StringIdentifier
	key.Identifier = json.Key
	var project model.Project
	if err := db.Where("id = ?", json.ProjectId).First(&project).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	}
	key.Project = project

	if dbc := db.Create(&key); dbc.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": dbc.Error.Error()})
		return
	}
	c.Status(201)
}

type updateKeyValidation struct {
	Key string `form:"key" json:"key" xml:"key"  binding:"required"`
}

func updateKey(c *gin.Context) {
	id := c.Param("id")
	var json updateKeyValidation
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var key model.StringIdentifier
	if err := db.Where("id = ?", id).First(&key).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
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
	}
	translation.Identifier = key

	var lang model.Language
	if err := db.Where("IsoCode = ?", json.Language).First(&lang).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
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
	}
	translation.Translation = json.Translation
	db.Save(&translation)
	c.Status(200)
}
func setApproved(c *gin.Context) {
	id := c.Param("id")
	var translation model.Translation
	if err := db.Where("id = ?", id).First(&translation).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
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
	}

	var key model.StringIdentifier
	if err := db.Where("id = ?", id).First(&key).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
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
	}
	db.Delete(&key)
	c.Status(200)
}
