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

var cs *CourseService

type CourseService struct {
	cc        *mongo.Collection
	logger    hclog.Logger
	configs   *utils.Configurations
	validator *models.Validation
}

func NewCourseService(logger hclog.Logger, configs *utils.Configurations, validator *models.Validation) *CourseService {
	return &CourseService{models.CoursesCollection, logger, configs, validator}
}

func (cs *CourseService) Create(ctx context.Context, course *models.Course) (*models.Course, *utils.RestError) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*20)
	defer cancel()
	course = models.NewCourse(course)
	result, err := cs.cc.InsertOne(ctx, course)
	if err != nil {
		RestError := utils.InternalErr("can't insert user to the database.")
		return nil, RestError
	}
	course.ID = result.InsertedID.(primitive.ObjectID)
	return course, nil
}

type GetQuery struct {
	ID       string `schema:"id"`
	Page     int64  `schema:"page"`
	Limit    int64  `schema:"limit"`
	Search   string `schema:"search"`
	Sort     string `schema:"sort"`
	CourseID string `schema:"course_id"`
	BookID   string `schema:"book_id"`
	Paralink string `schema:"paralink"`
}

func NewGetQuery(q *GetQuery) {
	if q.Limit == 0 {
		q.Limit = 20
	}
	if q.Page == 0 {
		q.Page = 1
	}
}
func (cs *CourseService) Find(ctx context.Context, params *GetQuery) (bson.M, *utils.RestError) {
	skip := (params.Page - 1) * params.Limit
	query := bson.M{}
	if params.ID != "" {
		_id, err := primitive.ObjectIDFromHex(params.ID)
		if err != nil {
			return nil, utils.InternalErr(err.Error())
		}
		query["_id"] = _id
	}
	if params.CourseID != "" {
		_id, err := primitive.ObjectIDFromHex(params.CourseID)
		if err != nil {
			return nil, utils.InternalErr(err.Error())
		}
		query["course_id"] = _id
	}
	// opts := options.Find().SetSkip(skip).SetLimit(params.Limit)

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
	setStage := bson.D{{Key: "$addFields", Value: bson.M{"image.url": bson.D{{"$concat", bson.A{cs.configs.AssetsUrl, "$image.path"}}}}}}
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
	unwindStage := bson.D{{"$unwind", "$total"}}

	pipeline := mongo.Pipeline{matchStage, imagePipelineStage, imageUnwindStage, setStage, facetStage, unwindStage}

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	cursor, err := cs.cc.Aggregate(ctx, pipeline)

	if err != nil {
		// return nil, utils.InternalErr("Internal Server Error")
		return nil, utils.InternalErr(err.Error())
	}
	defer cursor.Close(ctx)

	var courses []bson.M
	if err = cursor.All(context.TODO(), &courses); err != nil {
		RestError := utils.InternalErr("Internal Server Error")
		return nil, RestError
	}
	if len(courses) == 0 {
		return bson.M{"docs": []bson.M{}, "total": bson.M{"count": 0}}, nil
	}
	return courses[0], nil
}
func (cs *CourseService) FindById(ctx context.Context, course_id string) (*models.Course, *utils.RestError) {
	var course models.Course
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	id, e := primitive.ObjectIDFromHex(course_id)
	if e != nil {
		RestError := utils.NotFound("Invalid user_id")
		return nil, RestError
	}
	err := cs.cc.FindOne(ctx, bson.M{"_id": id}).Decode(&course)
	if err != nil {
		RestError := utils.NotFound("user not found.")
		return nil, RestError
	}
	return &course, nil
}

func (cs *CourseService) DeleteById(ctx context.Context, course_id string) *utils.RestError {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	id, e := primitive.ObjectIDFromHex(course_id)
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
		RestError := utils.NotFound("course not found.")
		return RestError
	}
	return nil
}

func (cs *CourseService) UpdateById(ctx context.Context, course_id string, updateCourse *models.Course) (*models.Course, *utils.RestError) {
	id, e := primitive.ObjectIDFromHex(course_id)
	if e != nil {
		RestError := utils.NotFound("Invalid user_id")
		return nil, RestError
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	course, RestError := cs.FindById(ctx, course_id)
	if RestError != nil {
		return nil, RestError
	}
	after := options.After
	opts := options.FindOneAndUpdateOptions{ReturnDocument: &after}
	result := cs.cc.FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.M{"$set": updateCourse}, &opts)
	if result.Err() != nil {
		return nil, utils.InternalErr(result.Err().Error())
	}
	// Decode the result
	decodeErr := result.Decode(course)
	if decodeErr != nil {
		return nil, utils.InternalErr(decodeErr.Error())
	}
	return course, nil
}
