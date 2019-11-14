package routing

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"preventis.io/translationApi/model"

	"golang.org/x/crypto/bcrypt"
)

func login(c *gin.Context) {
	var json loginValidation
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		println(err.Error())
		return
	}
	var users []model.User
	db.Where("Username = ?", json.LoginName, json.LoginName).Find(&users)
	if len(users) != 1 {
		c.Status(http.StatusNotFound)
		fmt.Printf("found %d users instead of one", len(users))
		return
	}
	user := users[0]

	pw := []byte(json.Password)
	hash := []byte(user.Password)
	err := bcrypt.CompareHashAndPassword(hash, pw)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		println(err.Error())
		return
	}

	session := sessions.Default(c)
	session.Set("username", user.Username)
	if user.Admin {
		session.Set("admin", user.Admin)
	}
	session.Set("userID", user.ID)
	err = session.Save()
	if err != nil {
		println(err.Error())
	}
	c.Status(http.StatusOK)
}

func logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	c.Status(http.StatusOK)
}

type userValidation struct {
	Name     string `form:"name" json:"name" xml:"name"  binding:"required"`
	Password string `form:"password" json:"password" xml:"password"  binding:"required"`
	Admin    bool   `form:"admin" json:"admin" xml:"admin"  binding:"required"`
	Mail     string `form:"mail" json:"mail" xml:"mail"  binding:"required"`
}

type loginValidation struct {
	LoginName string `form:"loginName" json:"loginName" xml:"loginName"  binding:"required"`
	Password  string `form:"password" json:"password" xml:"password"  binding:"required"`
}

func createUser(c *gin.Context) {
	session := sessions.Default(c)
	admin := session.Get("admin")
	if admin == nil {
		println("no admin")
		c.Status(http.StatusForbidden)
		return
	}
	var json userValidation
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		println(err.Error())
		return
	}

	saltedBytes := []byte(json.Password)
	hashedBytes, _ := bcrypt.GenerateFromPassword(saltedBytes, bcrypt.DefaultCost)

	hash := string(hashedBytes[:])

	var user model.User
	user.Username = json.Name
	user.Admin = json.Admin
	user.Mail = json.Mail
	user.Password = hash

	if dbc := db.Create(&user); dbc.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": dbc.Error.Error()})
		return
	}
	c.Status(http.StatusCreated)
}
