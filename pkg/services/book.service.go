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

var bs *BookService

type BookService struct {
	bc        *mongo.Collection
	logger    hclog.Logger
	configs   *utils.Configurations
	validator *models.Validation
}

func NewBookService(logger hclog.Logger, configs *utils.Configurations, validator *models.Validation) *BookService {
	return &BookService{models.BooksCollection, logger, configs, validator}
}

func (bs *BookService) Create(ctx context.Context, book *models.Book) (*models.Book, *utils.RestError) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*20)
	defer cancel()
	book = models.NewBook(book)
	result, err := bs.bc.InsertOne(ctx, book)
	if err != nil {
		RestError := utils.InternalErr("can't insert user to the database.")
		return nil, RestError
	}
	book.ID = result.InsertedID.(primitive.ObjectID)
	return book, nil
}

func (bs *BookService) Find(ctx context.Context, params *GetQuery) (bson.M, *utils.RestError) {
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
	imagePipelineStage := bson.D{
		{"$lookup", bson.D{
			{"from", "media"},
			{"let", bson.M{"image_id": "$image"}},
			{"pipeline", bson.A{
				bson.D{{
					"$match", bson.D{{
						"$expr",
						bson.D{{
							"$and",
							bson.A{
								bson.D{{"$eq", bson.A{"$_id", "$$image_id"}}},
							},
						}},
					}},
				}},
			},
			},
			{"as", "image"},
		},
		},
	}
	imageUnwindStage := bson.D{{"$unwind", bson.D{{"path", "$image"}, {"preserveNullAndEmptyArrays", true}}}}
	setStage := bson.D{{Key: "$addFields", Value: bson.M{"image.url": bson.D{{"$concat", bson.A{bs.configs.AssetsUrl, "$image.path"}}}}}}
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
	pipeline := mongo.Pipeline{matchStage, imagePipelineStage, imageUnwindStage, setStage, stockLookup, coursePipelineStage, courseUnwindStage, facetStage, unwindStage2}

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	cursor, err := bs.bc.Aggregate(ctx, pipeline)

	if err != nil {
		// return nil, utils.InternalErr("Internal Server Error")
		return nil, utils.InternalErr(err.Error())
	}
	defer cursor.Close(ctx)

	var books []bson.M
	if err = cursor.All(context.TODO(), &books); err != nil {
		// RestError := utils.InternalErr("Internal Server Error")
		RestError := utils.InternalErr(err.Error())
		return nil, RestError
	}
	if len(books) == 0 {
		return bson.M{"docs": []bson.M{}, "total": bson.M{"count": 0}}, nil
	}
	return books[0], nil
}
func (bs *BookService) FindById(ctx context.Context, book_id string) (*primitive.M, *utils.RestError) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	id, e := primitive.ObjectIDFromHex(book_id)
	if e != nil {
		RestError := utils.NotFound("Invalid user_id")
		return nil, RestError
	}
	matchStage := bson.D{{"$match", bson.M{"_id": id}}}
	mediaLookup := bson.D{{
		"$lookup", bson.D{
			{"from", "media"},
			{"let", bson.M{"image_id": "$image"}},
			{"pipeline", bson.A{
				bson.D{{
					"$match", bson.D{{
						"$expr",
						bson.D{{
							"$and",
							bson.A{
								bson.D{{"$eq", bson.A{"$_id", "$$image_id"}}},
							},
						}},
					}},
				}},
			},
			},
			{"as", "image"},
		},
	}}
	mediaUnWind := bson.D{{"$unwind", bson.D{{"path", "$image"}, {"preserveNullAndEmptyArrays", true}}}}
	mediaConcat := bson.D{{Key: "$addFields", Value: bson.M{"image.url": bson.D{{"$concat", bson.A{bs.configs.AssetsUrl, "$image.path"}}}}}}
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
	pipeline := mongo.Pipeline{matchStage, mediaLookup, mediaUnWind, mediaConcat, stockLookup}
	cursor, err := bs.bc.Aggregate(ctx, pipeline)

	if err != nil {
		// return nil, utils.InternalErr("Internal Server Error")
		return nil, utils.InternalErr(err.Error())
	}
	defer cursor.Close(ctx)

	var books []bson.M
	if err = cursor.All(context.TODO(), &books); err != nil {
		// RestError := utils.InternalErr("Internal Server Error")
		RestError := utils.InternalErr(err.Error())
		return nil, RestError
	}
	if len(books) == 0 {
		return nil, utils.BadRequest("No Book found")
	}
	return &books[0], nil
}

func (bs *BookService) DeleteById(ctx context.Context, book_id string) *utils.RestError {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	id, e := primitive.ObjectIDFromHex(book_id)
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
		RestError := utils.NotFound("book not found.")
		return RestError
	}
	return nil
}

func (bs *BookService) UpdateById(ctx context.Context, book_id string, updateBook *models.Book) (*primitive.M, *utils.RestError) {
	id, e := primitive.ObjectIDFromHex(book_id)
	if e != nil {
		RestError := utils.NotFound("Invalid user_id")
		return nil, RestError
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	book, RestError := bs.FindById(ctx, book_id)
	if RestError != nil {
		return nil, RestError
	}
	after := options.After
	opts := options.FindOneAndUpdateOptions{ReturnDocument: &after}
	result := bs.bc.FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.M{"$set": updateBook}, &opts)
	if result.Err() != nil {
		return nil, utils.InternalErr(result.Err().Error())
	}
	// Decode the result
	decodeErr := result.Decode(book)
	if decodeErr != nil {
		return nil, utils.InternalErr(decodeErr.Error())
	}
	return book, nil
}
