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

var RegisterCourseRoutes = func(router *mux.Router, logger hclog.Logger, configs *utils.Configurations, validator *models.Validation) {
	c := controllers.NewCourseController(logger, configs, validator)
	m := middlewares.NewMiddleware(logger, configs, validator)

	sr := router.PathPrefix("/courses").Subrouter()

	//Need to add AuthMiddleware With Access Control

	sr.Handle("", m.AuthWithRoles([]string{"admin"}, middlewares.AuthenticatedHandler(c.CreateCourse))).Methods(http.MethodPost)
	sr.Handle("", m.AuthWithRoles([]string{}, middlewares.AuthenticatedHandler(c.GetCourse))).Methods(http.MethodGet)
	sr.Handle("/{course_id}", m.AuthWithRoles([]string{}, middlewares.AuthenticatedHandler(c.GetCourseById))).Methods(http.MethodGet)
	sr.Handle("/{course_id}", m.AuthWithRoles([]string{"admin"}, middlewares.AuthenticatedHandler(c.UpdateCourse))).Methods(http.MethodPatch)
	sr.Handle("/{course_id}", m.AuthWithRoles([]string{"admin"}, middlewares.AuthenticatedHandler(c.DeleteCourse))).Methods(http.MethodDelete)
}
