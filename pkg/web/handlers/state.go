package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/stv0g/vand/pkg/store"
)

func HandleStateWith(store *store.Store) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Header("Content-type", "application/json")
		c.JSON(200, &store.State)
	}
}
