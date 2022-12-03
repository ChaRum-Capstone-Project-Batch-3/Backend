package follow_threads

import (
	followthreads "charum/business/follow_threads"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type followThreadRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) followthreads.Repository {
	return &followThreadRepository{
		collection: db.Collection("followThreads"),
	}
}

/*
Create
*/

func (ftr *followThreadRepository) Create(domain *followthreads.Domain) (followthreads.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	res, err := ftr.collection.InsertOne(ctx, FromDomain(domain))
	if err != nil {
		return followthreads.Domain{}, err
	}

	result, err := ftr.GetByID(res.InsertedID.(primitive.ObjectID))
	if err != nil {
		return followthreads.Domain{}, err
	}

	return result, err
}

/*
Read
*/

func (ftr *followThreadRepository) GetAllByUserID(userID primitive.ObjectID) ([]followthreads.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result []Model
	cursor, err := ftr.collection.Find(ctx, bson.M{
		"userID": userID,
	})
	if err != nil {
		return []followthreads.Domain{}, err
	}

	err = cursor.All(ctx, &result)
	if err != nil {
		return []followthreads.Domain{}, err
	}

	return ToDomainArray(result), nil
}

func (ftr *followThreadRepository) GetByID(id primitive.ObjectID) (followthreads.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result Model
	err := ftr.collection.FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&result)
	if err != nil {
		return followthreads.Domain{}, err
	}

	return result.ToDomain(), nil
}

func (ftr *followThreadRepository) GetByUserIDAndThreadID(userID primitive.ObjectID, threadID primitive.ObjectID) (followthreads.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result Model
	err := ftr.collection.FindOne(ctx, bson.M{
		"userID":   userID,
		"threadID": threadID,
	}).Decode(&result)
	if err != nil {
		return followthreads.Domain{}, err
	}

	return result.ToDomain(), nil
}

func (ftr *followThreadRepository) CountByThreadID(threadID primitive.ObjectID) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	count, err := ftr.collection.CountDocuments(ctx, bson.M{
		"threadID": threadID,
	})
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

/*
Update
*/

/*
Delete
*/

func (ftr *followThreadRepository) Delete(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := ftr.collection.DeleteOne(ctx, bson.M{
		"_id": id,
	})
	if err != nil {
		return err
	}

	return nil
}
