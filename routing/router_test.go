package routing

import (
	"github.com/gin-gonic/gin"
	"preventis.io/translationApi/model"
)

func setupTestRouter() *gin.Engine {
	db = model.StartDB("sqlite3", ":memory:")
	router := setupRouter()
	return router
}
