package routes

import (
	"net/http"

	"github.com/SairamVemula/booksland-backend-go/pkg/controllers"
	"github.com/SairamVemula/booksland-backend-go/pkg/models"
	"github.com/SairamVemula/booksland-backend-go/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
)

var RegisterAuthToutes = func(router *mux.Router, logger hclog.Logger, configs *utils.Configurations, validator *models.Validation) {
	c := controllers.NewAuthController(logger, configs, validator)

	sr := router.PathPrefix("/auth").Subrouter()

	sr.HandleFunc("/login", c.Login).Methods(http.MethodPost)
	sr.HandleFunc("/register", c.Register).Methods(http.MethodPost)
	sr.HandleFunc("/forgot-password/{email}", c.ForgotPassword).Methods(http.MethodGet)
	sr.HandleFunc("/reset-verification", c.ResentVerification).Methods(http.MethodPost)
	sr.HandleFunc("/reset-password", c.ResetPassword).Methods(http.MethodPost)
	sr.HandleFunc("/verify", c.Verify).Methods(http.MethodPost)
	sr.HandleFunc("/refresh", c.Refresh).Methods(http.MethodGet)
	sr.HandleFunc("/logout", c.Logout).Methods(http.MethodGet)
}
