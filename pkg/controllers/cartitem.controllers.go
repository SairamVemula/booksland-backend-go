package controllers

import (
	"net/http"

	"github.com/SairamVemula/booksland-backend-go/pkg/models"
	"github.com/SairamVemula/booksland-backend-go/pkg/services"
	"github.com/SairamVemula/booksland-backend-go/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
)

type CartItemController struct {
	cartItemService *services.CartItemService
	logger          hclog.Logger
	configs         *utils.Configurations
	validator       *models.Validation
}

func NewCartItemController(logger hclog.Logger, configs *utils.Configurations, validator *models.Validation) *CartItemController {
	return &CartItemController{services.NewCartItemService(logger, configs, validator), logger, configs, validator}
}

func (cic *CartItemController) Create(w http.ResponseWriter, r *http.Request, authUser *models.User) {
	cartItem := models.CartItem{}
	perr := utils.ParseBody(r, &cartItem)
	if perr != nil {
		utils.ResponseStringError(&w, perr.Error())
		return
	}
	models.NewCartItem(&cartItem)
	err := cic.validator.Struct(cartItem)
	if err != nil {
		utils.ResponseValidationError(&w, &err)
		return
	}
	cartItem.UserID = authUser.ID
	res, e := cic.cartItemService.Create(r.Context(), &cartItem)
	if e != nil {
		utils.ResponseError(&w, e)
		return
	}
	utils.ResponseSuccess(&w, res)
}

func (cic *CartItemController) Get(w http.ResponseWriter, r *http.Request, authUser *models.User) {
	var query services.GetQuery
	err := decoder.Decode(&query, r.URL.Query())
	if err != nil {
		utils.ResponseStringError(&w, err.Error())
		return
	}
	services.NewGetQuery(&query)

	res, e := cic.cartItemService.Find(r.Context(), &query)
	if e != nil {
		utils.ResponseError(&w, e)
		return
	}
	utils.ResponseSuccess(&w, res)
}
func (cic *CartItemController) GetById(w http.ResponseWriter, r *http.Request, authUser *models.User) {
	user_id := r.URL.Query().Get("id")
	if user_id == "" {
		utils.ResponseStringError(&w, "user_id is required")
	}
	res, e := cic.cartItemService.FindById(r.Context(), user_id)
	if e != nil {
		utils.ResponseError(&w, e)
		return
	}
	utils.ResponseSuccess(&w, res)

}
func (cic *CartItemController) Update(w http.ResponseWriter, r *http.Request, authUser *models.User) {
	params := mux.Vars(r)
	if params["cartItem_id"] == "" {
		utils.ResponseStringError(&w, "cartItem_id is required")
		return
	}
	cartItem := &models.CartItem{}
	perr := utils.ParseBody(r, &cartItem)
	if perr != nil {
		utils.ResponseStringError(&w, perr.Error())
		return
	}

	res, e := cic.cartItemService.UpdateById(r.Context(), params["cartItem_id"], cartItem)
	if e != nil {
		utils.ResponseError(&w, e)
		return
	}
	utils.ResponseSuccess(&w, res)
}

func (cic *CartItemController) Delete(w http.ResponseWriter, r *http.Request, authUser *models.User) {
	params := mux.Vars(r)
	if params["cartItem_id"] == "" {
		utils.ResponseStringError(&w, "cartItem_id is required")
		return
	}
	e := cic.cartItemService.DeleteById(r.Context(), params["cartItem_id"])
	if e != nil {
		utils.ResponseError(&w, e)
		return
	}
	utils.ResponseSuccess(&w, utils.Response{Code: 200, Message: "success"})

}
