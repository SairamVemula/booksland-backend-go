package controllers

import (
	"net/http"
	"strconv"

	"github.com/SairamVemula/booksland-backend-go/pkg/models"
	"github.com/SairamVemula/booksland-backend-go/pkg/services"
	"github.com/SairamVemula/booksland-backend-go/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
)

type UserController struct {
	userService *services.UserService
	logger      hclog.Logger
	configs     *utils.Configurations
	validator   *models.Validation
}

func NewUserController(logger hclog.Logger, configs *utils.Configurations, validator *models.Validation) *UserController {
	return &UserController{services.NewUserService(logger, configs, validator), logger, configs, validator}
}

func (uc *UserController) CreateUser(w http.ResponseWriter, r *http.Request, authUser *models.User) {
	user := models.User{}
	perr := utils.ParseBody(r, &user)
	if perr != nil {
		utils.ResponseStringError(&w, perr.Error())
		return
	}
	err := uc.validator.Struct(user)
	if err != nil {
		utils.ResponseValidationError(&w, &err)
		return
	}
	res, e := uc.userService.Create(r.Context(), &user)
	if e != nil {
		utils.ResponseError(&w, e)
		return
	}
	utils.ResponseSuccess(&w, res)
}
func (uc *UserController) GetUserDetails(w http.ResponseWriter, r *http.Request, authUser *models.User) {
	res, e := uc.userService.FindById(r.Context(), authUser.ID.Hex())
	if e != nil {
		utils.ResponseError(&w, e)
		return
	}
	utils.ResponseSuccess(&w, res)
	return
}
func (uc *UserController) GetUsers(w http.ResponseWriter, r *http.Request, authUser *models.User) {
	query := r.URL.Query()
	var (
		limit int64
		page  int64
	)
	if query.Get("limit") == "" {
		limit = 20
	} else {
		value, err := strconv.Atoi(query.Get("limit"))
		if err != nil {
			utils.ResponseStringError(&w, err.Error())
			return
		}
		limit = int64(value)
	}
	if query.Get("page") == "" {
		page = 1
	} else {
		value, err := strconv.Atoi(query.Get("page"))
		if err != nil {
			utils.ResponseStringError(&w, err.Error())
			return
		}
		page = int64(value)
	}

	res, e := uc.userService.Find(r.Context(), page, limit)
	if e != nil {
		utils.ResponseError(&w, e)
		return
	}
	utils.ResponseSuccess(&w, res)
}
func (uc *UserController) GetUserById(w http.ResponseWriter, r *http.Request, authUser *models.User) {
	user_id := r.URL.Query().Get("id")
	if user_id == "" {
		utils.ResponseStringError(&w, "user_id is required")
	}
	res, e := uc.userService.FindById(r.Context(), user_id)
	if e != nil {
		utils.ResponseError(&w, e)
		return
	}
	utils.ResponseSuccess(&w, res)

}
func (uc *UserController) UpdateUser(w http.ResponseWriter, r *http.Request, authUser *models.User) {
	params := mux.Vars(r)
	if params["user_id"] == "" {
		utils.ResponseStringError(&w, "user_id is required")
		return
	}
	user := &models.UpdateUser{}
	perr := utils.ParseBody(r, user)
	if perr != nil {
		utils.ResponseStringError(&w, perr.Error())
		return
	}
	res, e := uc.userService.UpdateById(r.Context(), params["user_id"], user)
	if e != nil {
		utils.ResponseError(&w, e)
		return
	}
	utils.ResponseSuccess(&w, res)
}

func (uc *UserController) ChangePassword(w http.ResponseWriter, r *http.Request, authUser *models.User) {
	params := mux.Vars(r)
	if params["user_id"] == "" {
		utils.ResponseStringError(&w, "user_id is required")
		return
	}
	cp := &models.ChangePassword{}
	perr := utils.ParseBody(r, cp)
	if perr != nil {
		utils.ResponseStringError(&w, perr.Error())
		return
	}
	res, e := uc.userService.ChangePassword(r.Context(), params["user_id"], cp)
	if e != nil {
		utils.ResponseError(&w, e)
		return
	}
	utils.ResponseSuccess(&w, res)
}

func (uc *UserController) DeleteUser(w http.ResponseWriter, r *http.Request, authUser *models.User) {
	params := mux.Vars(r)
	if params["user_id"] == "" {
		utils.ResponseStringError(&w, "user_id is required")
		return
	}
	e := uc.userService.DeleteById(r.Context(), params["user_id"])
	if e != nil {
		utils.ResponseError(&w, e)
		return
	}
	utils.ResponseSuccess(&w, utils.Response{Code: 200, Message: "success"})

}
