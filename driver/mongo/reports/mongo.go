package reports

import (
	"charum/business/reports"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type reportRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) reports.Repository {
	return &reportRepository{
		collection: db.Collection("reports"),
	}
}

/*
Create
*/

func (rr *reportRepository) Create(domain *reports.Domain) (reports.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	res, err := rr.collection.InsertOne(ctx, FromDomain(domain))
	if err != nil {
		return reports.Domain{}, err
	}

	result, err := rr.GetByID(res.InsertedID.(primitive.ObjectID))
	if err != nil {
		return reports.Domain{}, err
	}

	return result, nil
}

/*
Read
*/
func (rr *reportRepository) GetByID(id primitive.ObjectID) (reports.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result Model
	err := rr.collection.FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&result)
	if err != nil {
		return reports.Domain{}, err
	}

	return result.ToDomain(), nil
}

func (rr *reportRepository) GetByReportedID(id primitive.ObjectID) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	cursor, err := rr.collection.Find(ctx, bson.M{
		"reportedId": id,
	})
	if err != nil {
		return 0, err
	}

	var result []Model
	if err = cursor.All(ctx, &result); err != nil {
		return 0, err
	}

	domains := len(result)
	return domains, nil
}

func (rr *reportRepository) CheckByUserID(UserID primitive.ObjectID, ReportedID primitive.ObjectID) (reports.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result Model
	err := rr.collection.FindOne(ctx, bson.M{
		"reportedId": ReportedID,
		"userId":     UserID,
	}).Decode(&result)
	if err != nil {
		return reports.Domain{}, err
	}

	return result.ToDomain(), nil
}

func (rr *reportRepository) GetAll() ([]reports.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	cursor, err := rr.collection.Find(ctx, bson.M{})
	if err != nil {
		return []reports.Domain{}, err
	}

	var result []Model
	if err = cursor.All(ctx, &result); err != nil {
		return []reports.Domain{}, err
	}

	domains := ToDomainArray(result)
	return domains, nil
}

func (rr *reportRepository) GetAllReportedUsers() ([]reports.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	cursor, err := rr.collection.Find(ctx, bson.M{
		"reportedType": "user",
	})
	if err != nil {
		return []reports.Domain{}, err
	}

	var result []Model
	if err = cursor.All(ctx, &result); err != nil {
		return []reports.Domain{}, err
	}

	domains := ToDomainArray(result)
	return domains, nil
}

func (rr *reportRepository) GetAllReportedThreads() ([]reports.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	cursor, err := rr.collection.Find(ctx, bson.M{
		"reportedType": "thread",
	})
	if err != nil {
		return []reports.Domain{}, err
	}

	var result []Model
	if err = cursor.All(ctx, &result); err != nil {
		return []reports.Domain{}, err
	}

	domains := ToDomainArray(result)
	return domains, nil
}
