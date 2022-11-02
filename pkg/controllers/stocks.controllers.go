package controllers

import (
	"net/http"

	"github.com/SairamVemula/booksland-backend-go/pkg/models"
	"github.com/SairamVemula/booksland-backend-go/pkg/services"
	"github.com/SairamVemula/booksland-backend-go/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
)

type StockController struct {
	stockService *services.StockService
	logger       hclog.Logger
	configs      *utils.Configurations
	validator    *models.Validation
}

func NewStockController(logger hclog.Logger, configs *utils.Configurations, validator *models.Validation) *StockController {
	return &StockController{services.NewStockService(logger, configs, validator), logger, configs, validator}
}

func (sc *StockController) Create(w http.ResponseWriter, r *http.Request, authUser *models.User) {
	stock := models.Stock{}
	perr := utils.ParseBody(r, &stock)
	if perr != nil {
		utils.ResponseStringError(&w, perr.Error())
		return
	}
	stock.CreatedBy = authUser.ID
	models.NewStock(&stock)
	err := sc.validator.Struct(stock)
	if err != nil {
		utils.ResponseValidationError(&w, &err)
		return
	}
	res, e := sc.stockService.Create(r.Context(), &stock)
	if e != nil {
		utils.ResponseError(&w, e)
		return
	}
	utils.ResponseSuccess(&w, res)
}

func (sc *StockController) Get(w http.ResponseWriter, r *http.Request, authUser *models.User) {
	var query services.GetQuery
	err := decoder.Decode(&query, r.URL.Query())
	if err != nil {
		utils.ResponseStringError(&w, err.Error())
		return
	}
	services.NewGetQuery(&query)

	res, e := sc.stockService.Find(r.Context(), &query)
	if e != nil {
		utils.ResponseError(&w, e)
		return
	}
	utils.ResponseSuccess(&w, res)
}
func (sc *StockController) GetById(w http.ResponseWriter, r *http.Request, authUser *models.User) {
	user_id := r.URL.Query().Get("id")
	if user_id == "" {
		utils.ResponseStringError(&w, "user_id is required")
	}
	res, e := sc.stockService.FindById(r.Context(), user_id)
	if e != nil {
		utils.ResponseError(&w, e)
		return
	}
	utils.ResponseSuccess(&w, res)

}
func (sc *StockController) Update(w http.ResponseWriter, r *http.Request, authUser *models.User) {
	params := mux.Vars(r)
	if params["stock_id"] == "" {
		utils.ResponseStringError(&w, "stock_id is required")
		return
	}
	stock := &models.Stock{}
	perr := utils.ParseBody(r, &stock)
	if perr != nil {
		utils.ResponseStringError(&w, perr.Error())
		return
	}

	res, e := sc.stockService.UpdateById(r.Context(), params["stock_id"], stock)
	if e != nil {
		utils.ResponseError(&w, e)
		return
	}
	utils.ResponseSuccess(&w, res)
}

func (sc *StockController) Delete(w http.ResponseWriter, r *http.Request, authUser *models.User) {
	params := mux.Vars(r)
	if params["stock_id"] == "" {
		utils.ResponseStringError(&w, "stock_id is required")
		return
	}
	e := sc.stockService.DeleteById(r.Context(), params["stock_id"])
	if e != nil {
		utils.ResponseError(&w, e)
		return
	}
	utils.ResponseSuccess(&w, utils.Response{Code: 200, Message: "success"})

}
