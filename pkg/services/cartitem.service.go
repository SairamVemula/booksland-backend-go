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

var cis *CartItemService

type CartItemService struct {
	cic       *mongo.Collection
	logger    hclog.Logger
	configs   *utils.Configurations
	validator *models.Validation
}

func NewCartItemService(logger hclog.Logger, configs *utils.Configurations, validator *models.Validation) *CartItemService {
	return &CartItemService{models.CartItemsCollection, logger, configs, validator}
}

func (cis *CartItemService) Create(ctx context.Context, cartItem *models.CartItem) (*models.CartItem, *utils.RestError) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*20)
	defer cancel()
	cartItem = models.NewCartItem(cartItem)
	result, err := cis.cic.InsertOne(ctx, cartItem)
	if err != nil {
		RestError := utils.InternalErr("can't insert user to the database.")
		return nil, RestError
	}
	cartItem.ID = result.InsertedID.(primitive.ObjectID)
	return cartItem, nil
}

func (cis *CartItemService) Find(ctx context.Context, params *GetQuery) (bson.M, *utils.RestError) {
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
	stockLookup := bson.D{{
		"$lookup", bson.D{
			{"from", "stocks"},
			{"let", bson.M{"book_id": "$_id"}},
			{"pipeline", bson.A{
				bson.D{{
					"$match", bson.D{{
						"$expr",
						bson.D{{
							"$and",
							bson.A{
								bson.D{{"$eq", bson.A{"$book_id", "$$book_id"}}},
								bson.D{{"$eq", bson.A{"$status", "available"}}},
							},
						}},
					}},
				}},
				bson.D{{
					"$group", bson.D{
						{"_id", bson.D{{"publisher", "$publisher"}, {"year", "$year"}}},
						{"prices", bson.D{{"$push", "$price"}}},
						{"discount_percents", bson.D{{"$push", "$discount_percent"}}},
						{"count", bson.D{{"$sum", 1}}},
					},
				}},
			},
			},
			{"as", "stocks"},
		},
	}}
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
				stockLookup,
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

	sortStage := bson.D{{Key: "$sort", Value: bson.D{{"_id", 1}}}}
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
	cursor, err := cis.cic.Aggregate(ctx, pipeline)

	if err != nil {
		// return nil, utils.InternalErr("Internal Server Error")
		return nil, utils.InternalErr(err.Error())
	}
	defer cursor.Close(ctx)

	var cartItems []bson.M
	if err = cursor.All(context.TODO(), &cartItems); err != nil {
		// RestError := utils.InternalErr("Internal Server Error")
		RestError := utils.InternalErr(err.Error())
		return nil, RestError
	}
	if len(cartItems) == 0 {
		return bson.M{"docs": []bson.M{}, "total": bson.M{"count": 0}}, nil
	}
	return cartItems[0], nil
}
func (cis *CartItemService) FindById(ctx context.Context, cartItem_id string) (*models.CartItem, *utils.RestError) {
	var cartItem models.CartItem
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	id, e := primitive.ObjectIDFromHex(cartItem_id)
	if e != nil {
		RestError := utils.NotFound("Invalid user_id")
		return nil, RestError
	}
	err := cis.cic.FindOne(ctx, bson.M{"_id": id}).Decode(&cartItem)
	if err != nil {
		RestError := utils.NotFound("user not found.")
		return nil, RestError
	}
	return &cartItem, nil
}

func (cis *CartItemService) DeleteById(ctx context.Context, cartItem_id string) *utils.RestError {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	id, e := primitive.ObjectIDFromHex(cartItem_id)
	if e != nil {
		RestError := utils.NotFound("Invalid cartItem_id")
		return RestError
	}
	result, err := cis.cic.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		RestError := utils.NotFound("faild to delete.")
		return RestError
	}
	if result.DeletedCount == 0 {
		RestError := utils.NotFound("cartItem not found.")
		return RestError
	}
	return nil
}

func (cis *CartItemService) UpdateById(ctx context.Context, cartItem_id string, updateCartItem *models.CartItem) (*models.CartItem, *utils.RestError) {
	id, e := primitive.ObjectIDFromHex(cartItem_id)
	if e != nil {
		RestError := utils.NotFound("Invalid cartitem_id")
		return nil, RestError
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	cartItem, RestError := cis.FindById(ctx, cartItem_id)
	if RestError != nil {
		return nil, RestError
	}
	after := options.After
	opts := options.FindOneAndUpdateOptions{ReturnDocument: &after}
	result := cis.cic.FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.M{"$set": updateCartItem}, &opts)
	if result.Err() != nil {
		return nil, utils.InternalErr(result.Err().Error())
	}
	// Decode the result
	decodeErr := result.Decode(cartItem)
	if decodeErr != nil {
		return nil, utils.InternalErr(decodeErr.Error())
	}
	return cartItem, nil
}
