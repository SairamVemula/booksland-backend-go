package controllers

import (
	"net/http"
	"time"

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
	models.NewUser(user)
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
	// cookie := http.Cookie{
	// 	Name:     "access_token",
	// 	Value:    res.Token,
	// 	Path:     "/",
	// 	Expires:  time.Now().Add(time.Minute * time.Duration(ac.configs.JwtExpiration)),
	// 	HttpOnly: true,
	// 	Domain:   "http://localhost:3000/",
	// }
	// http.SetCookie(w, &cookie)
	// cookie = http.Cookie{
	// 	Name:     "refresh_token",
	// 	Value:    res.RefreshToken,
	// 	Path:     "/",
	// 	Expires:  time.Now().Add(time.Minute * time.Duration(ac.configs.RefreshJwtExpiration)),
	// 	HttpOnly: true,
	// 	Domain:   "http://localhost:3000/",
	// }
	// http.SetCookie(w, &cookie)
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
	TokenExpiry int64  `json:"token_expiry"`
}

func (ac *AuthController) Refresh(w http.ResponseWriter, r *http.Request) {
	refresh_token := r.URL.Query().Get("refresh_token")
	if refresh_token == "" {
		utils.ResponseStringError(&w, "refresh_token is required")
	}
	user_id, user_type, err := ac.authService.ValidateRefreshToken(r.Context(), refresh_token)
	if err != nil {
		utils.ResponseStringError(&w, err.Error())
		return
	}
	token, err := ac.authService.GenerateAccessToken(r.Context(), user_id, user_type)
	if err != nil {
		utils.ResponseStringError(&w, err.Error())
		return
	}
	// cookie := http.Cookie{
	// 	Name:     "access_token",
	// 	Value:    token,
	// 	Path:     "/",
	// 	Expires:  time.Now().Add(time.Minute * time.Duration(ac.configs.JwtExpiration)),
	// 	HttpOnly: true,
	// 	Domain:   "http://localhost:3000/",
	// }
	// http.SetCookie(w, &cookie)
	utils.ResponseSuccess(&w, &TokenResponse{token, (time.Now().UnixMilli() + int64(ac.configs.JwtExpiration*60000)) - (1 * 60 * 1000)})
}

func (ac *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	ac.authService.Logout(r.Context(), &models.User{})
	utils.ResponseSuccess(&w, &utils.Response{Code: 200, Message: "success", Data: ""})
}
