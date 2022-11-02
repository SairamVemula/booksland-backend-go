package services

import (
	"context"
	"os"
	"time"

	"github.com/SairamVemula/booksland-backend-go/pkg/models"
	"github.com/SairamVemula/booksland-backend-go/pkg/utils"
	"github.com/hashicorp/go-hclog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var mds *MediaService

type MediaService struct {
	mdc       *mongo.Collection
	logger    hclog.Logger
	configs   *utils.Configurations
	validator *models.Validation
}

func NewMediaService(logger hclog.Logger, configs *utils.Configurations, validator *models.Validation) *MediaService {
	return &MediaService{models.MediaCollection, logger, configs, validator}
}

func (mds *MediaService) Create(ctx context.Context, media *models.Media) (*models.Media, *utils.RestError) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*20)
	defer cancel()

	result, err := mds.mdc.InsertOne(ctx, media)
	if err != nil {
		RestError := utils.InternalErr("can't insert user to the database.")
		return nil, RestError
	}
	media.ID = result.InsertedID.(primitive.ObjectID)
	media.Path = mds.configs.AppendUrl(media.Path)
	return media, nil
}

func (mds *MediaService) Find(ctx context.Context, params *GetQuery) (bson.M, *utils.RestError) {
	skip := (params.Page - 1) * params.Limit
	query := bson.M{}
	// opts := options.Find().SetSkip(skip).SetLimit(params.Limit)

	// log.Printf("%+v\n", params)

	if params.Search != "" {
		query["$or"] = bson.A{
			bson.M{
				"path": bson.M{"$regex": params.Search, "$options": "i"},
			},
		}
	}

	matchStage := bson.D{{"$match", query}}
	// projectStage := bson.D{
	// 	{"$project", bson.D{
	// 		{"_id", 1},
	// 		// {"path", bson.D{{"$concat", bson.A{mds.configs.AssetsUrl, "$path"}}}},
	// 		{"created_by", 1},
	// 		{"created_on", 1},
	// 		{"updated_on", 1},
	// 	},
	// 	},
	// }
	setStage := bson.D{{Key: "$addFields", Value: bson.M{"url": bson.D{{"$concat", bson.A{mds.configs.AssetsUrl, "$path"}}}}}}
	sortStage := bson.D{{Key: "$sort", Value: bson.D{{"_id", -1}}}}
	skipStage := bson.D{{Key: "$skip", Value: skip}}
	limitStage := bson.D{{Key: "$limit", Value: params.Limit}}
	countStage := bson.D{{Key: "$count", Value: "count"}}

	facetStage := bson.D{{
		"$facet", bson.D{
			{"docs", bson.A{sortStage, skipStage, limitStage}},
			{"total", bson.A{countStage}},
		},
	}}
	unwindStage := bson.D{{"$unwind", "$total"}}

	pipeline := mongo.Pipeline{matchStage, setStage, facetStage, unwindStage}

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	cursor, err := mds.mdc.Aggregate(ctx, pipeline)
	if err != nil {
		// return nil, utils.InternalErr("Internal Server Error")
		return nil, utils.InternalErr(err.Error())
	}
	defer cursor.Close(ctx)

	var medias []bson.M
	if err = cursor.All(context.TODO(), &medias); err != nil {
		RestError := utils.InternalErr("Internal Server Error")
		return nil, RestError
	}
	if len(medias) == 0 {
		return bson.M{"docs": []bson.M{}, "total": bson.M{"count": 0}}, nil
	}
	return medias[0], nil
}
func (mds *MediaService) FindById(ctx context.Context, media_id string) (*models.Media, *utils.RestError) {
	var media models.Media
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	id, e := primitive.ObjectIDFromHex(media_id)
	if e != nil {
		RestError := utils.NotFound("Invalid media_id")
		return nil, RestError
	}
	err := mds.mdc.FindOne(ctx, bson.M{"_id": id}).Decode(&media)
	if err != nil {
		RestError := utils.NotFound("media not found.")
		return nil, RestError
	}
	return &media, nil
}

func (mds *MediaService) DeleteById(ctx context.Context, media_id string) *utils.RestError {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	id, e := primitive.ObjectIDFromHex(media_id)
	if e != nil {
		RestError := utils.NotFound("Invalid media_id")
		return RestError
	}
	var media models.Media
	err := cs.cc.FindOne(ctx, bson.M{"_id": id}).Decode(&media)
	if err != nil {
		RestError := utils.NotFound("media not found")
		return RestError
	}
	err = os.Remove("." + media.Path)
	if err != nil {
		RestError := utils.NotFound(err.Error())
		return RestError
	}

	result, err := cs.cc.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		RestError := utils.NotFound("media to delete.")
		return RestError
	}
	if result.DeletedCount == 0 {
		RestError := utils.NotFound("media not found.")
		return RestError
	}
	return nil
}
