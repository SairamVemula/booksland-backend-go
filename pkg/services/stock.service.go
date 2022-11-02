package services

import (
	"context"
	"time"

	"github.com/SairamVemula/booksland-backend-go/pkg/models"
	"github.com/SairamVemula/booksland-backend-go/pkg/utils"
	"github.com/hashicorp/go-hclog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ss *StockService

type StockService struct {
	sc        *mongo.Collection
	logger    hclog.Logger
	configs   *utils.Configurations
	validator *models.Validation
}

func NewStockService(logger hclog.Logger, configs *utils.Configurations, validator *models.Validation) *StockService {
	return &StockService{models.StocksCollection, logger, configs, validator}
}

func (ss *StockService) Create(ctx context.Context, stock *models.Stock) (*models.Stock, *utils.RestError) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*20)
	defer cancel()
	stock = models.NewStock(stock)
	result, err := ss.sc.InsertOne(ctx, stock)
	if err != nil {
		RestError := utils.InternalErr("can't insert user to the database.")
		return nil, RestError
	}
	stock.ID = result.InsertedID.(primitive.ObjectID)
	return stock, nil
}

func (ss *StockService) Find(ctx context.Context, params *GetQuery) (bson.M, *utils.RestError) {
	skip := (params.Page - 1) * params.Limit
	query := bson.M{}
	// opts := options.Find().SetSkip(skip).SetLimit(params.Limit)
	if params.ID != "" {
		_id, err := primitive.ObjectIDFromHex(params.ID)
		if err != nil {
			return nil, utils.InternalErr(err.Error())
		}
		query["_id"] = _id
	}

	if params.Search != "" {
		query["$or"] = bson.A{
			bson.M{
				"name": bson.M{"$regex": params.Search, "$options": "i"},
			},
			bson.M{
				"tags": bson.M{"$elemMatch": bson.M{"$regex": params.Search, "$options": "i"}},
			},
		}
	}

	matchStage := bson.D{{"$match", query}}
	bookPipelineStage := bson.D{
		{"$lookup", bson.D{
			{"from", "books"},
			{"let", bson.M{"book_id": "$book_id"}},
			{"pipeline", bson.A{
				bson.D{{
					"$match", bson.D{{
						"$expr",
						bson.D{{
							"$and",
							bson.A{
								bson.D{{"$eq", bson.A{"$_id", "$$book_id"}}},
							},
						}},
					}},
				}},
			},
			},
			{"as", "book"},
		},
		},
	}
	bookUnwindStage := bson.D{{"$unwind", bson.D{{"path", "$book"}, {"preserveNullAndEmptyArrays", true}}}}

	coursePipelineStage := bson.D{
		{"$lookup", bson.D{
			{"from", "courses"},
			{"let", bson.M{"course_id": "$course_id"}},
			{"pipeline", bson.A{
				bson.D{{
					"$match", bson.D{{
						"$expr",
						bson.D{{
							"$and",
							bson.A{
								bson.D{{"$eq", bson.A{"$_id", "$$course_id"}}},
							},
						}},
					}},
				}},
			},
			},
			{"as", "course"},
		},
		},
	}
	courseUnwindStage := bson.D{{"$unwind", bson.D{{"path", "$course"}, {"preserveNullAndEmptyArrays", true}}}}

	sortStage := bson.D{{Key: "$sort", Value: bson.D{{"order", 1}}}}
	skipStage := bson.D{{Key: "$skip", Value: skip}}
	limitStage := bson.D{{Key: "$limit", Value: params.Limit}}
	countStage := bson.D{{Key: "$count", Value: "count"}}

	facetStage := bson.D{{
		"$facet", bson.D{
			{"docs", bson.A{sortStage, skipStage, limitStage}},
			{"total", bson.A{countStage}},
		},
	}}
	unwindStage2 := bson.D{{"$unwind", "$total"}}
	pipeline := mongo.Pipeline{matchStage, bookPipelineStage, bookUnwindStage, coursePipelineStage, courseUnwindStage, facetStage, unwindStage2}

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	cursor, err := ss.sc.Aggregate(ctx, pipeline)

	if err != nil {
		// return nil, utils.InternalErr("Internal Server Error")
		return nil, utils.InternalErr(err.Error())
	}
	defer cursor.Close(ctx)

	var stocks []bson.M
	if err = cursor.All(context.TODO(), &stocks); err != nil {
		// RestError := utils.InternalErr("Internal Server Error")
		RestError := utils.InternalErr(err.Error())
		return nil, RestError
	}
	if len(stocks) == 0 {
		return bson.M{"docs": []bson.M{}, "total": bson.M{"count": 0}}, nil
	}
	return stocks[0], nil
}
func (ss *StockService) FindById(ctx context.Context, stock_id string) (*models.Stock, *utils.RestError) {
	var stock models.Stock
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	id, e := primitive.ObjectIDFromHex(stock_id)
	if e != nil {
		RestError := utils.NotFound("Invalid user_id")
		return nil, RestError
	}
	err := ss.sc.FindOne(ctx, bson.M{"_id": id}).Decode(&stock)
	if err != nil {
		RestError := utils.NotFound("user not found.")
		return nil, RestError
	}
	return &stock, nil
}

func (ss *StockService) DeleteById(ctx context.Context, stock_id string) *utils.RestError {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	id, e := primitive.ObjectIDFromHex(stock_id)
	if e != nil {
		RestError := utils.NotFound("Invalid user_id")
		return RestError
	}
	result, err := cs.cc.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		RestError := utils.NotFound("faild to delete.")
		return RestError
	}
	if result.DeletedCount == 0 {
		RestError := utils.NotFound("stock not found.")
		return RestError
	}
	return nil
}

func (ss *StockService) UpdateById(ctx context.Context, stock_id string, updateBook *models.Stock) (*models.Stock, *utils.RestError) {
	id, e := primitive.ObjectIDFromHex(stock_id)
	if e != nil {
		RestError := utils.NotFound("Invalid user_id")
		return nil, RestError
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	stock, RestError := ss.FindById(ctx, stock_id)
	if RestError != nil {
		return nil, RestError
	}
	after := options.After
	opts := options.FindOneAndUpdateOptions{ReturnDocument: &after}
	result := ss.sc.FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.M{"$set": updateBook}, &opts)
	if result.Err() != nil {
		return nil, utils.InternalErr(result.Err().Error())
	}
	// Decode the result
	decodeErr := result.Decode(stock)
	if decodeErr != nil {
		return nil, utils.InternalErr(decodeErr.Error())
	}
	return stock, nil
}
