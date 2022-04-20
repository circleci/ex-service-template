package api

import (
	"context"
	"net/http"

	"github.com/circleci/ex/httpserver/ginrouter"
	"github.com/gin-gonic/gin"

	"github.com/circleci/ex-service-template/books"
)

type API struct {
	router *gin.Engine
	store  *books.Store
}

type Options struct {
	Store *books.Store
}

func New(ctx context.Context, opts Options) *API {
	r := ginrouter.Default(ctx, "api")
	a := &API{
		router: r,
		store:  opts.Store,
	}

	r.GET("/api/ping", a.ping)
	r.GET("/api/books/:id", a.getBook)
	r.GET("/api/books", a.getAllBooks)
	r.POST("/api/books", a.postBook)

	return a
}

func (a *API) Handler() http.Handler {
	return a.router
}
