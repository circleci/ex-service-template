package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (API) ping(c *gin.Context) {
	type response struct {
		Message string `json:"message"`
	}

	c.JSON(http.StatusOK, response{Message: "pong"})
}
