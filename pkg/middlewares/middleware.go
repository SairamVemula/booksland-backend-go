package middlewares

import (
	"github.com/SairamVemula/booksland-backend-go/pkg/models"
	"github.com/SairamVemula/booksland-backend-go/pkg/utils"
	"github.com/hashicorp/go-hclog"
)

type Middleware struct {
	logger    hclog.Logger
	configs   *utils.Configurations
	validator *models.Validation
}

func NewMiddleware(logger hclog.Logger, configs *utils.Configurations, validator *models.Validation) *Middleware {
	return &Middleware{logger, configs, validator}
}
