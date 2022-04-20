package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/circleci/ex-service-template/books"
)

type book struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Price string    `json:"price"`
}

func (a *API) getBook(c *gin.Context) {
	ctx := c.Request.Context()
	idString := c.Param("id")

	id, err := uuid.Parse(idString)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
		return
	}

	b, err := a.store.ByID(ctx, id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{})
		return
	}

	c.JSON(http.StatusOK, book(*b))
}

func (a *API) getAllBooks(c *gin.Context) {
	ctx := c.Request.Context()

	res, err := a.store.All(ctx)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{})
		return
	}

	all := make([]book, len(res))
	for i := range res {
		all[i] = book(res[i])
	}

	c.JSON(http.StatusOK, all)
}

func (a *API) postBook(c *gin.Context) {
	type request struct {
		Name  string `json:"name" binding:"required"`
		Price string `db:"price" binding:"required"`
	}
	type response struct {
		ID uuid.UUID `json:"id"`
	}

	ctx := c.Request.Context()

	var req request
	err := c.BindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
		return
	}

	id, err := a.store.Add(ctx, books.ToAdd(req))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusOK, response{
		ID: id,
	})
}
