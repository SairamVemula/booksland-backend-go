package controllers

import (
	"net/http"

	"github.com/SairamVemula/booksland-backend-go/pkg/models"
	"github.com/SairamVemula/booksland-backend-go/pkg/services"
	"github.com/SairamVemula/booksland-backend-go/pkg/utils"
	"github.com/hashicorp/go-hclog"
)

type AuthController struct {
	authService *services.AuthService
	logger      hclog.Logger
	configs     *utils.Configurations
	validator   *models.Validation
}

func NewAuthController(logger hclog.Logger, configs *utils.Configurations, validator *models.Validation) *AuthController {
	return &AuthController{services.NewAuthService(logger, configs, validator), logger, configs, validator}
}

func (ac *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	perr := utils.ParseBody(r, user)
	if perr != nil {
		utils.ResponseStringError(&w, perr.Error())
		return
	}
	err := ac.validator.Struct(user)
	if err != nil {
		utils.ResponseValidationError(&w, &err)
		return
	}
	res, e := ac.authService.Register(r.Context(), user)
	if e != nil {
		utils.ResponseError(&w, e)
		return
	}
	utils.ResponseSuccess(&w, res)

}
func (ac *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	user := &models.LoginUser{}
	perr := utils.ParseBody(r, user)
	if perr != nil {
		utils.ResponseStringError(&w, perr.Error())
		return
	}
	err := ac.validator.Struct(user)
	if err != nil {
		utils.ResponseValidationError(&w, &err)
		return
	}
	res, e := ac.authService.Login(r.Context(), user)
	if e != nil {
		utils.ResponseError(&w, e)
		return
	}
	utils.ResponseSuccess(&w, res)

}
func (ac *AuthController) ResentVerification(w http.ResponseWriter, r *http.Request) {

}
func (ac *AuthController) Verify(w http.ResponseWriter, r *http.Request) {

}
func (ac *AuthController) ForgotPassword(w http.ResponseWriter, r *http.Request) {

}
func (ac *AuthController) ResetPassword(w http.ResponseWriter, r *http.Request) {

}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

func (ac *AuthController) Refresh(w http.ResponseWriter, r *http.Request) {
	refresh_token := r.URL.Query().Get("refresh_token")
	if refresh_token == "" {
		utils.ResponseStringError(&w, "refresh_token is required")
	}
	user_id, err := ac.authService.ValidateRefreshToken(refresh_token)
	if err != nil {
		utils.ResponseStringError(&w, err.Error())
		return
	}
	token, err := ac.authService.GenerateAccessToken(r.Context(), user_id, "admin")
	if err != nil {
		utils.ResponseStringError(&w, err.Error())
		return
	}
	utils.ResponseSuccess(&w, &TokenResponse{token})
}

func (ac *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	ac.authService.Logout(r.Context(), &models.User{})
	utils.ResponseSuccess(&w, &utils.Response{Code: 200, Message: "success", Data: ""})
}
