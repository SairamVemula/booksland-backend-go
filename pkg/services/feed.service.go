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

var fs *FeedService

type FeedService struct {
	bc        *mongo.Collection
	logger    hclog.Logger
	configs   *utils.Configurations
	validator *models.Validation
}

func NewFeedService(logger hclog.Logger, configs *utils.Configurations, validator *models.Validation) *FeedService {
	return &FeedService{models.FeedsCollection, logger, configs, validator}
}

func (fs *FeedService) Create(ctx context.Context, feed *models.Feed) (*models.Feed, *utils.RestError) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*20)
	defer cancel()
	feed = models.NewFeed(feed)
	result, err := fs.bc.InsertOne(ctx, feed)
	if err != nil {
		RestError := utils.InternalErr("can't insert user to the database.")
		return nil, RestError
	}
	feed.ID = result.InsertedID.(primitive.ObjectID)
	return feed, nil
}

func (fs *FeedService) Find(ctx context.Context, params *GetQuery) (bson.M, *utils.RestError) {
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
	if params.Paralink != "" {
		query["paralink"] = params.Paralink
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
	sectionsUnWindStage := bson.D{{"$unwind", bson.D{{"path", "$sections"}, {"preserveNullAndEmptyArrays", true}}}}

	coursePipelineStage := bson.D{
		{"$lookup", bson.D{
			{"from", "courses"},
			{"let", bson.M{"course_id": "$sections.course"}},
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
				bson.D{{
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
				}},
				bson.D{{"$unwind", bson.D{{"path", "$image"}, {"preserveNullAndEmptyArrays", true}}}},
				bson.D{{Key: "$addFields", Value: bson.M{"image.url": bson.D{{"$concat", bson.A{fs.configs.AssetsUrl, "$image.path"}}}}}},
			},
			},
			{"as", "sections.course"},
		},
		},
	}
	courseUnwindStage := bson.D{{"$unwind", bson.D{{"path", "$sections.course"}, {"preserveNullAndEmptyArrays", true}}}}

	bookPipelineStage := bson.D{
		{"$lookup", bson.D{
			{"from", "books"},
			{"let", bson.M{"book_id": "$sections.book"}},
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
				bson.D{{
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
				}},
				bson.D{{"$unwind", bson.D{{"path", "$image"}, {"preserveNullAndEmptyArrays", true}}}},
				bson.D{{Key: "$addFields", Value: bson.M{"image.url": bson.D{{"$concat", bson.A{fs.configs.AssetsUrl, "$image.path"}}}}}},
				bson.D{{
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
				}},
			},
			},
			{"as", "sections.book"},
		},
		},
	}
	bookUnwindStage := bson.D{{"$unwind", bson.D{{"path", "$sections.book"}, {"preserveNullAndEmptyArrays", true}}}}

	imagePipelineStage := bson.D{
		{"$lookup", bson.D{
			{"from", "media"},
			{"let", bson.M{"image_id": "$sections.image"}},
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
			{"as", "sections.image"},
		},
		},
	}
	imageUnwindStage := bson.D{{"$unwind", bson.D{{"path", "$sections.image"}, {"preserveNullAndEmptyArrays", true}}}}
	setStage := bson.D{{Key: "$addFields", Value: bson.M{"sections.image.url": bson.D{{"$concat", bson.A{fs.configs.AssetsUrl, "$sections.image.path"}}}}}}

	groupStage1 := bson.D{
		{"$group", bson.D{
			{"_id", "$_id"},
			{"sections", bson.M{"$push": "$sections"}},
			{"created_on", bson.M{"$first": "$created_on"}},
			{"updated_on", bson.M{"$first": "$updated_on"}},
			{"name", bson.M{"$first": "$name"}},
			{"paralink", bson.M{"$first": "$paralink"}},
			{"title", bson.M{"$first": "$title"}},
			{"type", bson.M{"$first": "$type"}},
			{"view_type", bson.M{"$first": "$view_type"}},
			{"created_by", bson.M{"$first": "$created_by"}},
			{"order", bson.M{"$first": "$order"}},
		}},
	}

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
	pipeline := mongo.Pipeline{
		matchStage,
		sectionsUnWindStage,
		imagePipelineStage,
		imageUnwindStage,
		coursePipelineStage,
		courseUnwindStage,
		bookPipelineStage,
		bookUnwindStage,
		setStage,
		groupStage1,
		facetStage,
		unwindStage2,
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	cursor, err := fs.bc.Aggregate(ctx, pipeline)

	if err != nil {
		// return nil, utils.InternalErr("Internal Server Error")
		return nil, utils.InternalErr(err.Error())
	}
	defer cursor.Close(ctx)

	var feeds []bson.M
	if err = cursor.All(context.TODO(), &feeds); err != nil {
		// RestError := utils.InternalErr("Internal Server Error")
		RestError := utils.InternalErr(err.Error())
		return nil, RestError
	}
	if len(feeds) == 0 {
		return bson.M{"docs": []bson.M{}, "total": bson.M{"count": 0}}, nil
	}
	return feeds[0], nil
}
func (fs *FeedService) FindById(ctx context.Context, feed_id string) (*models.Feed, *utils.RestError) {
	var feed models.Feed
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	id, e := primitive.ObjectIDFromHex(feed_id)
	if e != nil {
		RestError := utils.NotFound("Invalid user_id")
		return nil, RestError
	}
	err := fs.bc.FindOne(ctx, bson.M{"_id": id}).Decode(&feed)
	if err != nil {
		RestError := utils.NotFound("user not found.")
		return nil, RestError
	}
	return &feed, nil
}

func (fs *FeedService) DeleteById(ctx context.Context, feed_id string) *utils.RestError {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	id, e := primitive.ObjectIDFromHex(feed_id)
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
		RestError := utils.NotFound("feed not found.")
		return RestError
	}
	return nil
}

func (fs *FeedService) UpdateById(ctx context.Context, feed_id string, updateFeed *models.Feed) (*models.Feed, *utils.RestError) {
	id, e := primitive.ObjectIDFromHex(feed_id)
	if e != nil {
		RestError := utils.NotFound("Invalid user_id")
		return nil, RestError
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	feed, RestError := fs.FindById(ctx, feed_id)
	if RestError != nil {
		return nil, RestError
	}
	after := options.After
	opts := options.FindOneAndUpdateOptions{ReturnDocument: &after}
	result := fs.bc.FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.M{"$set": updateFeed}, &opts)
	if result.Err() != nil {
		return nil, utils.InternalErr(result.Err().Error())
	}
	// Decode the result
	decodeErr := result.Decode(feed)
	if decodeErr != nil {
		return nil, utils.InternalErr(decodeErr.Error())
	}
	return feed, nil
}
