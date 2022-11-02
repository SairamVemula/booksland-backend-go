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

var RegisterUsersRoutes = func(router *mux.Router, logger hclog.Logger, configs *utils.Configurations, validator *models.Validation) {
	c := controllers.NewUserController(logger, configs, validator)
	m := middlewares.NewMiddleware(logger, configs, validator)

	sr := router.PathPrefix("/users").Subrouter()

	sr.Handle("", m.AuthWithRoles([]string{"admin"}, middlewares.AuthenticatedHandler(c.CreateUser))).Methods(http.MethodPost)
	sr.Handle("", m.AuthWithRoles([]string{"admin"}, middlewares.AuthenticatedHandler(c.GetUsers))).Methods(http.MethodGet)
	sr.Handle("/details", m.AuthWithRoles([]string{"admin", "user"}, middlewares.AuthenticatedHandler(c.GetUserDetails))).Methods(http.MethodGet)
	sr.Handle("/{user_id}", m.AuthWithRoles([]string{"admin"}, middlewares.AuthenticatedHandler(c.GetUserById))).Methods(http.MethodGet)
	sr.Handle("/{user_id}", m.AuthWithRoles([]string{"admin"}, middlewares.AuthenticatedHandler(c.UpdateUser))).Methods(http.MethodPatch)
	sr.Handle("/{user_id}", m.AuthWithRoles([]string{"admin"}, middlewares.AuthenticatedHandler(c.DeleteUser))).Methods(http.MethodDelete)
}
