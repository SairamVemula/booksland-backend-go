package utils

import (
	"encoding/json"
	"net/http"

	"github.com/SairamVemula/booksland-backend-go/pkg/models"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func setupCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "*")
}
func ResponseSuccess(w *http.ResponseWriter, response interface{}) {
	setupCORS(w)
	res, _ := json.Marshal(Response{http.StatusOK, "success", response})
	(*w).WriteHeader(http.StatusOK)
	(*w).Write(res)
}

func ResponseValidationError(w *http.ResponseWriter, err *models.ValidationErrors) {
	ResponseStringError(w, err.Errors()[0])
}

func ResponseStringError(w *http.ResponseWriter, s string) {
	setupCORS(w)
	res, _ := json.Marshal(BadRequest(s))
	(*w).WriteHeader(http.StatusBadRequest)
	(*w).Write(res)
}
func ResponseError(w *http.ResponseWriter, response *RestError) {
	setupCORS(w)
	res, _ := json.Marshal(response)
	(*w).WriteHeader(response.Code)
	(*w).Write(res)
}
