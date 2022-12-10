package controllers

import (
	"net/http"

	"github.com/SairamVemula/booksland-backend-go/pkg/models"
	"github.com/SairamVemula/booksland-backend-go/pkg/services"
	"github.com/SairamVemula/booksland-backend-go/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
)

type MediaController struct {
	mediaService *services.MediaService
	logger       hclog.Logger
	configs      *utils.Configurations
	validator    *models.Validation
}

func NewMediaController(logger hclog.Logger, configs *utils.Configurations, validator *models.Validation) *MediaController {
	return &MediaController{services.NewMediaService(logger, configs, validator), logger, configs, validator}
}

func (mdc *MediaController) Create(w http.ResponseWriter, r *http.Request, authUser *models.User) {
	path, err := utils.UploadFile(r, "file")
	if err != nil || path == "empty" {
		utils.ResponseStringError(&w, err.Error())
		return
	}
	if path == "empty" {
		utils.ResponseStringError(&w, "File is required")
		return
	}
	if err := r.ParseForm(); err != nil {
		utils.ResponseStringError(&w, err.Error())
		return
	}
	media := &models.Media{Path: path, CreatedBy: authUser.ID}
	media = models.NewMedia(media)
	e := mdc.validator.Struct(media)
	if e != nil {
		utils.ResponseValidationError(&w, &e)
		return
	}
	res, error := mdc.mediaService.Create(r.Context(), media)
	if error != nil {
		utils.ResponseError(&w, error)
		return
	}
	utils.ResponseSuccess(&w, res)
}

func (mdc *MediaController) Get(w http.ResponseWriter, r *http.Request /**, authUser *models.User) **/) {
	var query services.GetQuery
	err := decoder.Decode(&query, r.URL.Query())
	if err != nil {
		utils.ResponseStringError(&w, err.Error())
		return
	}
	services.NewGetQuery(&query)

	res, e := mdc.mediaService.Find(r.Context(), &query)
	if e != nil {
		utils.ResponseError(&w, e)
		return
	}
	utils.ResponseSuccess(&w, res)
}
func (mdc *MediaController) GetById(w http.ResponseWriter, r *http.Request /**, authUser *models.User) **/) {
	media_id := mux.Vars(r)["media_id"]
	if media_id == "" {
		utils.ResponseStringError(&w, "media_id is required")
		return
	}
	res, e := mdc.mediaService.FindById(r.Context(), media_id)
	if e != nil {
		utils.ResponseError(&w, e)
		return
	}
	// data, err := os.ReadFile("."+res.Path);
	// if err != nil {
	// 	utils.ResponseStringError(&w, err.Error())
	// 	return
	// }
	// w.Write(data)
	http.ServeFile(w, r, "."+res.Path)
	// utils.ResponseSuccess(&w, res)
}

func (mdc *MediaController) Delete(w http.ResponseWriter, r *http.Request /**, authUser *models.User) **/) {
	params := mux.Vars(r)
	if params["media_id"] == "" {
		utils.ResponseStringError(&w, "media_id is required")
		return
	}
	e := mdc.mediaService.DeleteById(r.Context(), params["media_id"])
	if e != nil {
		utils.ResponseError(&w, e)
		return
	}
	utils.ResponseSuccess(&w, utils.Response{Code: 200, Message: "success"})

}
