package routing

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"preventis.io/translationApi/model"
)

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
	id := c.Param("id")
	lang := c.Param("lang")
	var project model.Project
	if err := db.Where("id = ?", id).
		Preload("Languages").
		Preload("Identifiers").
		Preload("Identifiers.Translations").
		First(&project).
		Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
		return
	} else {
		if !containsLanguage(lang, project) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Language not present in project"})
			println("language not in project")
			return
		}
		resList := getAndroidExportStrings(project, lang)
		res := ""
		for _, e := range resList {
			res += fmt.Sprintf(`"%s" = "%s";\n`, e.Identifier, e.Translation)
		}

		// not using c.XML because of formatting
		c.String(http.StatusOK, res)
	}
}

type androidResource struct {
	XMLName xml.Name        `xml:"resources"`
	Strings []androidString `xml:"string"`
}

type androidString struct {
	Identifier  string `xml:"name,attr"`
	Translation string `xml:",chardata"`
}

func exportAndroid(c *gin.Context) {
	id := c.Param("id")
	lang := c.Param("lang")
	var project model.Project
	if err := db.Where("id = ?", id).
		Preload("Languages").
		Preload("Identifiers").
		Preload("Identifiers.Translations").
		First(&project).
		Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
		return
	} else {
		if !containsLanguage(lang, project) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Language not present in project"})
			println("language not in project")
			return
		}
		resList := getAndroidExportStrings(project, lang)
		resource := androidResource{Strings: resList}
		w := &bytes.Buffer{}
		enc := xml.NewEncoder(w)
		enc.Indent("", "  ")
		if err := enc.Encode(resource); err != nil {
			panic(err)
		}

		// not using c.XML because of formatting
		c.String(http.StatusOK, w.String())
	}
}

func getAndroidExportStrings(project model.Project, lang string) []androidString {
	var resList []androidString
	for _, e := range project.Identifiers {
		for _, t := range e.Translations {
			if t.LanguageRefer == lang {
				translation := androidString{Identifier: e.Identifier, Translation: t.Translation}
				resList = append(resList, translation)
			}
		}
	}
	return resList
}

func exportExcel(c *gin.Context) {
	// TODO
}
