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

var RegisterCartRoutes = func(router *mux.Router, logger hclog.Logger, configs *utils.Configurations, validator *models.Validation) {
	c := controllers.NewCartItemController(logger, configs, validator)
	m := middlewares.NewMiddleware(logger, configs, validator)

	sr := router.PathPrefix("/cart").Subrouter()

	sr.Handle("", m.AuthWithRoles([]string{"admin", "user"}, middlewares.AuthenticatedHandler(c.Create))).Methods(http.MethodPost)
	sr.Handle("", m.AuthWithRoles([]string{"admin", "user"}, middlewares.AuthenticatedHandler(c.Get))).Methods(http.MethodGet)
	sr.Handle("/{cartItem_id}", m.AuthWithRoles([]string{"admin", "user"}, middlewares.AuthenticatedHandler(c.GetById))).Methods(http.MethodGet)
	sr.Handle("/{cartItem_id}", m.AuthWithRoles([]string{"admin", "user"}, middlewares.AuthenticatedHandler(c.Update))).Methods(http.MethodPatch)
	sr.Handle("/{cartItem_id}", m.AuthWithRoles([]string{"admin", "user"}, middlewares.AuthenticatedHandler(c.Delete))).Methods(http.MethodDelete)
}
