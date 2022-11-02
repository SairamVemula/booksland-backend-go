package services

import (
	"github.com/SairamVemula/booksland-backend-go/pkg/models"
	"github.com/SairamVemula/booksland-backend-go/pkg/utils"
	"github.com/hashicorp/go-hclog"
	"go.mongodb.org/mongo-driver/mongo"
)

type MailService interface {
	sendMail(user *models.User, code string) error
	sendSMS(user *models.User, code string) error
	verify(user *models.User, code string) error
	sendForgotPasswordMail(user *models.User, code string) error
	sendForgotPasswordSMS(user *models.User, code string) error
	verifyForgetPassword(user *models.User, code string) error
}

var ms *Mail

type Mail struct {
	vc        *mongo.Collection
	fpc       *mongo.Collection
	logger    hclog.Logger
	configs   *utils.Configurations
	validator *models.Validation
}

func NewMailService(logger hclog.Logger, configs *utils.Configurations, validator *models.Validation) *Mail {
	if ms != nil {
		return ms
	}
	return &Mail{models.DB.Collection("verification"), models.DB.Collection("forgotpassword"), logger, configs, validator}
}
