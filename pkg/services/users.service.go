package services

import (
	"context"
	"log"
	"time"

	"github.com/SairamVemula/booksland-backend-go/pkg/models"
	"github.com/SairamVemula/booksland-backend-go/pkg/utils"
	"github.com/hashicorp/go-hclog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var us *UserService

type UserService struct {
	uc        *mongo.Collection
	logger    hclog.Logger
	configs   *utils.Configurations
	validator *models.Validation
}

func NewUserService(logger hclog.Logger, configs *utils.Configurations, validator *models.Validation) *UserService {
	return &UserService{models.UsersCollection, logger, configs, validator}
}

func (us *UserService) Create(ctx context.Context, user *models.User) (*models.User, *utils.RestError) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*20)
	defer cancel()
	models.NewUser(user)
	newPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		RestError := utils.InternalErr("can't insert user to the database.")
		return nil, RestError
	}
	user.Password = string(newPassword)
	if user.Type == "" {
		user.Type = "user"
	}
	result, err := us.uc.InsertOne(ctx, user)
	if err != nil {
		RestError := utils.InternalErr("can't insert user to the database.")
		return nil, RestError
	}
	user.ID = result.InsertedID.(primitive.ObjectID)
	user.Password = ""
	return user, nil
}

func (us *UserService) Find(ctx context.Context, page int64, limit int64) (*[]*models.User, *utils.RestError) {
	skip := (page - 1) * limit
	query := bson.M{}
	opts := options.Find().SetSkip(skip).SetLimit(limit)

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	cursor, err := us.uc.Find(ctx, query, opts)
	if err != nil {
		// Handle error
	}
	defer cursor.Close(ctx)

	var users []*models.User
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			RestError := utils.InternalErr("Internal Server Error")
			return nil, RestError
		}
		user.Password = ""
		users = append(users, &user)
	}
	return &users, nil
}
func (us *UserService) FindById(ctx context.Context, user_id string) (*models.User, *utils.RestError) {
	var user models.User
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	id, e := primitive.ObjectIDFromHex(user_id)
	if e != nil {
		RestError := utils.NotFound("Invalid user_id")
		return nil, RestError
	}
	err := us.uc.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		RestError := utils.NotFound("user not found.")
		return nil, RestError
	}

	return &user, nil
}

func (us *UserService) DeleteById(ctx context.Context, user_id string) *utils.RestError {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	id, e := primitive.ObjectIDFromHex(user_id)
	if e != nil {
		RestError := utils.NotFound("Invalid user_id")
		return RestError
	}
	result, err := us.uc.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		RestError := utils.NotFound("faild to delete.")
		return RestError
	}
	if result.DeletedCount == 0 {
		RestError := utils.NotFound("user not found.")
		return RestError
	}
	return nil
}

func (us *UserService) UpdateById(ctx context.Context, user_id string, updateUser *models.UpdateUser) (*models.User, *utils.RestError) {
	id, e := primitive.ObjectIDFromHex(user_id)
	if e != nil {
		RestError := utils.NotFound("Invalid user_id")
		return nil, RestError
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	after := options.After
	opts := options.FindOneAndUpdateOptions{ReturnDocument: &after}
	result := us.uc.FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.M{"$set": updateUser}, &opts)
	if result.Err() != nil {
		return nil, utils.InternalErr(result.Err().Error())
	}
	// Decode the result
	user := models.User{}
	decodeErr := result.Decode(&user)
	if decodeErr != nil {
		return nil, utils.InternalErr(decodeErr.Error())
	}
	user.Password = ""
	return &user, nil
}

func (us *UserService) ChangePassword(ctx context.Context, user_id string, changePwd *models.ChangePassword) (*models.User, *utils.RestError) {
	user, RestError := us.FindById(ctx, user_id)
	if RestError != nil {
		return nil, RestError
	}
	pe := bcrypt.CompareHashAndPassword([]byte(changePwd.Password), []byte(user.Password))
	if pe != nil {
		RestError := utils.InternalErr("Incorrect password")
		return nil, RestError
	}
	var newPassword []byte
	if changePwd.Password != "" {
		var err error
		newPassword, err = bcrypt.GenerateFromPassword([]byte(changePwd.Password), 10)
		if err != nil {
			RestError := utils.InternalErr("Internal Server Error")
			return nil, RestError
		}
	}

	id, e := primitive.ObjectIDFromHex(user_id)
	if e != nil {
		RestError := utils.NotFound("Invalid user_id")
		return nil, RestError
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	result, err := us.uc.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"password": string(newPassword)}})
	if err != nil {
		RestError := utils.InternalErr("can not update.")
		return nil, RestError
	}
	if result.MatchedCount == 0 {
		RestError := utils.NotFound("user not found.")
		return nil, RestError
	}
	if result.ModifiedCount == 0 {
		RestError := utils.BadRequest("no such field")
		return nil, RestError
	}
	user.Password = ""
	return user, nil
}

func (us *UserService) FindUsernameAndPassword(ctx context.Context, username string, password string) (user *models.User, RestError *utils.RestError) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	err := us.uc.FindOne(ctx, bson.M{"$or": bson.A{bson.M{"phone": username}, bson.M{"email": username}}}).Decode(&user)
	if err != nil {
		log.Println(err.Error())
		RestError := utils.NotFound("User not Found.")
		return nil, RestError
	}

	pe := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if pe != nil {
		RestError := utils.InternalErr("Incorrect password")
		return nil, RestError
	}

	return user, nil
}
