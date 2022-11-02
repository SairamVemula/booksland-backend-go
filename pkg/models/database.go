package models

import (
	"context"
	"log"

	"github.com/hashicorp/go-hclog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	DB                *mongo.Database
	Client            *mongo.Client
	UsersCollection   *mongo.Collection
	MediaCollection   *mongo.Collection
	BooksCollection   *mongo.Collection
	CoursesCollection *mongo.Collection
	StocksCollection  *mongo.Collection
	OrdersCollection  *mongo.Collection
	FeedsCollection   *mongo.Collection
)

func Connect(uri string, dbname string, logger hclog.Logger) error {

	clientOptions := options.Client().ApplyURI(uri) // Connect to //MongoDB
	Client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		logger.Error(err.Error())
	}
	// Check the connection
	err = Client.Ping(context.TODO(), nil)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	DB = Client.Database(dbname)

	UsersCollection = DB.Collection("users")
	BooksCollection = DB.Collection("books")
	MediaCollection = DB.Collection("media")
	FeedsCollection = DB.Collection("feeds")
	CoursesCollection = DB.Collection("courses")
	StocksCollection = DB.Collection("stocks")
	OrdersCollection = DB.Collection("orders")

	log.Println("Connected to MongoDB!")
	return nil
}
