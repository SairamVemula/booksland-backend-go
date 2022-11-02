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

var RegisterMediaRoutes = func(router *mux.Router, logger hclog.Logger, configs *utils.Configurations, validator *models.Validation) {
	c := controllers.NewMediaController(logger, configs, validator)
	m := middlewares.NewMiddleware(logger, configs, validator)

	sr := router.PathPrefix("/media").Subrouter()

	// sr.HandleFunc("", c.Create).Methods(http.MethodPost)
	sr.HandleFunc("", c.Get).Methods(http.MethodGet)
	sr.HandleFunc("/{media_id}", c.GetById).Methods(http.MethodGet)
	sr.HandleFunc("/{media_id}", c.Delete).Methods(http.MethodDelete)

	sr.Handle("", m.AuthWithRoles([]string{"admin"}, middlewares.AuthenticatedHandler(c.Create))).Methods(http.MethodPost)
	// sr.Handle("", m.AuthWithRoles([]string{}, middlewares.AuthenticatedHandler(c.Delete))).Methods(http.MethodGet)
	// sr.Handle("/{book_id}", m.AuthWithRoles([]string{}, middlewares.AuthenticatedHandler(c.GetById))).Methods(http.MethodGet)
	// sr.Handle("/{book_id}", m.AuthWithRoles([]string{"admin"}, middlewares.AuthenticatedHandler(c.Update))).Methods(http.MethodPatch)
	// sr.Handle("/{book_id}", m.AuthWithRoles([]string{"admin"}, middlewares.AuthenticatedHandler(c.Delete))).Methods(http.MethodDelete)
}
