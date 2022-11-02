package controllers

import (
	"net/http"
	"strings"

	"github.com/SairamVemula/booksland-backend-go/pkg/models"
	"github.com/SairamVemula/booksland-backend-go/pkg/services"
	"github.com/SairamVemula/booksland-backend-go/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
)

type BookController struct {
	bookService *services.BookService
	logger      hclog.Logger
	configs     *utils.Configurations
	validator   *models.Validation
}

func NewBookController(logger hclog.Logger, configs *utils.Configurations, validator *models.Validation) *BookController {
	return &BookController{services.NewBookService(logger, configs, validator), logger, configs, validator}
}

func (bc *BookController) Create(w http.ResponseWriter, r *http.Request, authUser *models.User) {
	book := &models.Book{}
	err := utils.ParseBody(r, book)
	if err != nil {
		utils.ResponseStringError(&w, err.Error())
		return
	}
	if len(book.Tags) > 0 {
		book.Tags = strings.Split(strings.Trim(book.Tags[0], " "), ",")
	}
	models.NewBook(book)
	e := bc.validator.Struct(book)
	if e != nil {
		utils.ResponseValidationError(&w, &e)
		return
	}
	res, error := bc.bookService.Create(r.Context(), book)
	if error != nil {
		utils.ResponseError(&w, error)
		return
	}
	utils.ResponseSuccess(&w, res)
}

func (bc *BookController) Get(w http.ResponseWriter, r *http.Request, authUser *models.User) {
	var query services.GetQuery
	err := decoder.Decode(&query, r.URL.Query())
	if err != nil {
		utils.ResponseStringError(&w, err.Error())
		return
	}
	services.NewGetQuery(&query)

	res, e := bc.bookService.Find(r.Context(), &query)
	if e != nil {
		utils.ResponseError(&w, e)
		return
	}
	utils.ResponseSuccess(&w, res)
}
func (bc *BookController) GetById(w http.ResponseWriter, r *http.Request, authUser *models.User) {
	user_id := r.URL.Query().Get("id")
	if user_id == "" {
		utils.ResponseStringError(&w, "user_id is required")
	}
	res, e := bc.bookService.FindById(r.Context(), user_id)
	if e != nil {
		utils.ResponseError(&w, e)
		return
	}
	utils.ResponseSuccess(&w, res)

}
func (bc *BookController) Update(w http.ResponseWriter, r *http.Request, authUser *models.User) {
	params := mux.Vars(r)
	if params["book_id"] == "" {
		utils.ResponseStringError(&w, "book_id is required")
		return
	}
	book := &models.Book{}
	err := utils.ParseBody(r, book)
	if err != nil {
		utils.ResponseStringError(&w, err.Error())
		return
	}
	if len(book.Tags) > 0 {
		book.Tags = strings.Split(strings.Trim(book.Tags[0], " "), ",")
	}
	book.CreatedBy = authUser.ID
	res, e := bc.bookService.UpdateById(r.Context(), params["book_id"], book)
	if e != nil {
		utils.ResponseError(&w, e)
		return
	}
	utils.ResponseSuccess(&w, res)
}

func (bc *BookController) Delete(w http.ResponseWriter, r *http.Request, authUser *models.User) {
	params := mux.Vars(r)
	if params["book_id"] == "" {
		utils.ResponseStringError(&w, "book_id is required")
		return
	}
	e := bc.bookService.DeleteById(r.Context(), params["book_id"])
	if e != nil {
		utils.ResponseError(&w, e)
		return
	}
	utils.ResponseSuccess(&w, utils.Response{Code: 200, Message: "success"})

}
