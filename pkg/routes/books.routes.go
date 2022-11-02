package routes

import (
	"net/http"

	"github.com/SairamVemula/booksland-backend-go/pkg/controllers"
	"github.com/SairamVemula/booksland-backend-go/pkg/middlewares"
	"github.com/SairamVemula/booksland-backend-go/pkg/models"
	"github.com/SairamVemula/booksland-backend-go/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
)

var RegisterBooksRoutes = func(router *mux.Router, logger hclog.Logger, configs *utils.Configurations, validator *models.Validation) {
	c := controllers.NewBookController(logger, configs, validator)
	m := middlewares.NewMiddleware(logger, configs, validator)

	sr := router.PathPrefix("/books").Subrouter()

	sr.Handle("", m.AuthWithRoles([]string{"admin"}, middlewares.AuthenticatedHandler(c.Create))).Methods(http.MethodPost)
	sr.Handle("", m.AuthWithRoles([]string{}, middlewares.AuthenticatedHandler(c.Get))).Methods(http.MethodGet)
	sr.Handle("/{book_id}", m.AuthWithRoles([]string{}, middlewares.AuthenticatedHandler(c.GetById))).Methods(http.MethodGet)
	sr.Handle("/{book_id}", m.AuthWithRoles([]string{"admin"}, middlewares.AuthenticatedHandler(c.Update))).Methods(http.MethodPatch)
	sr.Handle("/{book_id}", m.AuthWithRoles([]string{"admin"}, middlewares.AuthenticatedHandler(c.Delete))).Methods(http.MethodDelete)
}
