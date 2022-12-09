package bookmarks

import (
	"charum/business/bookmarks"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type bookmarkRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) bookmarks.Repository {
	return &bookmarkRepository{
		collection: db.Collection("bookmarks"),
	}
}

/*
Create
*/

func (br *bookmarkRepository) Create(domain *bookmarks.Domain) (bookmarks.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := br.collection.InsertOne(ctx, FromDomain(domain))
	if err != nil {
		return bookmarks.Domain{}, err
	}

	result, err := br.GetByUserIDAndThreadID(domain.UserID, domain.ThreadID)
	if err != nil {
		return bookmarks.Domain{}, err
	}

	return result, err
}

/*
Read
*/

func (br *bookmarkRepository) GetByUserIDAndThreadID(UserID primitive.ObjectID, ThreadID primitive.ObjectID) (bookmarks.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result Model
	err := br.collection.FindOne(ctx, bson.M{
		"userID":   UserID,
		"threadID": ThreadID,
	}).Decode(&result)
	if err != nil {
		return bookmarks.Domain{}, err
	}

	return result.ToDomain(), nil
}

func (br *bookmarkRepository) GetAllByUserID(UserID primitive.ObjectID) ([]bookmarks.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	cursor, err := br.collection.Find(ctx, bson.M{
		"userID": UserID,
	})
	if err != nil {
		return []bookmarks.Domain{}, err
	}

	var results []Model
	if err = cursor.All(ctx, &results); err != nil {
		return []bookmarks.Domain{}, err
	}

	domains := ToDomainArray(results)

	return domains, nil
}

func (br *bookmarkRepository) CountByThreadID(threadID primitive.ObjectID) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	count, err := br.collection.CountDocuments(ctx, bson.M{
		"threadID": threadID,
	})
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

/*
Delete
*/

func (br *bookmarkRepository) Delete(domain *bookmarks.Domain) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := br.collection.DeleteOne(ctx, bson.M{
		"userID":   domain.UserID,
		"threadID": domain.ThreadID,
	})
	if err != nil {
		return err
	}

	return nil
}

func (br *bookmarkRepository) DeleteAllByUserID(userID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := br.collection.DeleteMany(ctx, bson.M{
		"userID": userID,
	})
	if err != nil {
		return err
	}

	return nil
}

func (br *bookmarkRepository) DeleteAllByThreadID(threadID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := br.collection.DeleteMany(ctx, bson.M{
		"threadID": threadID,
	})
	if err != nil {
		return err
	}

	return nil
}
