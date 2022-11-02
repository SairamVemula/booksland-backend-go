package controllers

import (
	"net/http"
	"strings"

	"github.com/SairamVemula/booksland-backend-go/pkg/models"
	"github.com/SairamVemula/booksland-backend-go/pkg/services"
	"github.com/SairamVemula/booksland-backend-go/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/hashicorp/go-hclog"
)

var decoder = schema.NewDecoder()

type CourseController struct {
	courseService *services.CourseService
	logger        hclog.Logger
	configs       *utils.Configurations
	validator     *models.Validation
}

func NewCourseController(logger hclog.Logger, configs *utils.Configurations, validator *models.Validation) *CourseController {
	return &CourseController{services.NewCourseService(logger, configs, validator), logger, configs, validator}
}

func (cc *CourseController) CreateCourse(w http.ResponseWriter, r *http.Request, authUser *models.User) {
	course := &models.Course{}
	err := utils.ParseBody(r, course)
	if err != nil {
		utils.ResponseStringError(&w, err.Error())
		return
	}
	// if len(course.Tags) > 0 {
	// 	course.Tags = strings.Split(strings.Trim(course.Tags[0], " "), ",")
	// }
	course.CreatedBy = authUser.ID
	models.NewCourse(course)
	e := cc.validator.Struct(course)
	if e != nil {
		utils.ResponseValidationError(&w, &e)
		return
	}
	res, error := cc.courseService.Create(r.Context(), course)
	if error != nil {
		utils.ResponseError(&w, error)
		return
	}
	utils.ResponseSuccess(&w, res)
}

func (cc *CourseController) GetCourse(w http.ResponseWriter, r *http.Request, authUser *models.User) {
	var query services.GetQuery
	err := decoder.Decode(&query, r.URL.Query())
	if err != nil {
		utils.ResponseStringError(&w, err.Error())
		return
	}
	services.NewGetQuery(&query)

	res, e := cc.courseService.Find(r.Context(), &query)
	if e != nil {
		utils.ResponseError(&w, e)
		return
	}
	utils.ResponseSuccess(&w, res)
}
func (cc *CourseController) GetCourseById(w http.ResponseWriter, r *http.Request, authUser *models.User) {
	course_id := r.URL.Query().Get("course_id")
	if course_id == "" {
		utils.ResponseStringError(&w, "course_id is required")
	}
	res, e := cc.courseService.FindById(r.Context(), course_id)
	if e != nil {
		utils.ResponseError(&w, e)
		return
	}
	utils.ResponseSuccess(&w, res)

}
func (cc *CourseController) UpdateCourse(w http.ResponseWriter, r *http.Request, authUser *models.User) {
	params := mux.Vars(r)
	if params["course_id"] == "" {
		utils.ResponseStringError(&w, "course_id is required")
		return
	}
	course := &models.Course{}
	err := utils.ParseBody(r, course)
	if err != nil {
		utils.ResponseStringError(&w, err.Error())
		return
	}
	if len(course.Tags) > 0 {
		course.Tags = strings.Split(strings.Trim(course.Tags[0], " "), ",")
	}
	res, e := cc.courseService.UpdateById(r.Context(), params["course_id"], course)
	if e != nil {
		utils.ResponseError(&w, e)
		return
	}
	utils.ResponseSuccess(&w, res)
}

func (cc *CourseController) DeleteCourse(w http.ResponseWriter, r *http.Request, authUser *models.User) {
	params := mux.Vars(r)
	if params["course_id"] == "" {
		utils.ResponseStringError(&w, "course_id is required")
		return
	}
	e := cc.courseService.DeleteById(r.Context(), params["course_id"])
	if e != nil {
		utils.ResponseError(&w, e)
		return
	}
	utils.ResponseSuccess(&w, utils.Response{Code: 200, Message: "success"})

}
