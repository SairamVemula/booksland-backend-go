package controllers

import (
	"net/http"

	"github.com/SairamVemula/booksland-backend-go/pkg/models"
	"github.com/SairamVemula/booksland-backend-go/pkg/services"
	"github.com/SairamVemula/booksland-backend-go/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
)

type FeedController struct {
	feedService *services.FeedService
	logger      hclog.Logger
	configs     *utils.Configurations
	validator   *models.Validation
}

func NewFeedController(logger hclog.Logger, configs *utils.Configurations, validator *models.Validation) *FeedController {
	return &FeedController{services.NewFeedService(logger, configs, validator), logger, configs, validator}
}

func (fc *FeedController) Create(w http.ResponseWriter, r *http.Request, authUser *models.User) {
	feed := &models.Feed{}
	err := utils.ParseBody(r, feed)
	if err != nil {
		utils.ResponseStringError(&w, err.Error())
		return
	}
	models.NewFeed(feed)
	e := fc.validator.Struct(feed)
	if e != nil {
		utils.ResponseValidationError(&w, &e)
		return
	}
	res, error := fc.feedService.Create(r.Context(), feed)
	if error != nil {
		utils.ResponseError(&w, error)
		return
	}
	utils.ResponseSuccess(&w, res)
}

func (fc *FeedController) Get(w http.ResponseWriter, r *http.Request, authUser *models.User) {
	var query services.GetQuery
	err := decoder.Decode(&query, r.URL.Query())
	if err != nil {
		utils.ResponseStringError(&w, err.Error())
		return
	}
	services.NewGetQuery(&query)

	res, e := fc.feedService.Find(r.Context(), &query)
	if e != nil {
		utils.ResponseError(&w, e)
		return
	}
	utils.ResponseSuccess(&w, res)
}
func (fc *FeedController) GetById(w http.ResponseWriter, r *http.Request, authUser *models.User) {
	user_id := r.URL.Query().Get("id")
	if user_id == "" {
		utils.ResponseStringError(&w, "user_id is required")
	}
	res, e := fc.feedService.FindById(r.Context(), user_id)
	if e != nil {
		utils.ResponseError(&w, e)
		return
	}
	utils.ResponseSuccess(&w, res)

}
func (fc *FeedController) Update(w http.ResponseWriter, r *http.Request, authUser *models.User) {
	params := mux.Vars(r)
	if params["feed_id"] == "" {
		utils.ResponseStringError(&w, "feed_id is required")
		return
	}
	feed := &models.Feed{}
	err := utils.ParseBody(r, feed)
	if err != nil {
		utils.ResponseStringError(&w, err.Error())
		return
	}
	feed.CreatedBy = authUser.ID
	res, e := fc.feedService.UpdateById(r.Context(), params["feed_id"], feed)
	if e != nil {
		utils.ResponseError(&w, e)
		return
	}
	utils.ResponseSuccess(&w, res)
}

func (fc *FeedController) Delete(w http.ResponseWriter, r *http.Request, authUser *models.User) {
	params := mux.Vars(r)
	if params["feed_id"] == "" {
		utils.ResponseStringError(&w, "feed_id is required")
		return
	}
	e := fc.feedService.DeleteById(r.Context(), params["feed_id"])
	if e != nil {
		utils.ResponseError(&w, e)
		return
	}
	utils.ResponseSuccess(&w, utils.Response{Code: 200, Message: "success"})

}
