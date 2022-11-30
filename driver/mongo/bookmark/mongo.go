package bookmark

import (
	"charum/business/bookmarks"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
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

func (br *bookmarkRepository) AddBookmark(domain *bookmarks.Domain) (bookmarks.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	res, err := br.collection.InsertOne(ctx, FromDomain(domain))
	if err != nil {
		return bookmarks.Domain{}, err
	}

	// return data
	result, err := br.GetByID(res.InsertedID.(primitive.ObjectID))
	if err != nil {
		return bookmarks.Domain{}, err
	}

	return result, err
}

/*
Read
*/

func (br *bookmarkRepository) GetByID(id primitive.ObjectID) (bookmarks.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result Bookmark
	err := br.collection.FindOne(ctx, bson.M{
		"userId": id,
	}).Decode(&result)
	if err != nil {
		return bookmarks.Domain{}, err
	}
	return result.ToDomain(), nil
}

// update
func (br *bookmarkRepository) UpdateBookmark(domain *bookmarks.Domain) (bookmarks.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// only update threads & updatedAt
	_, err := br.collection.UpdateOne(ctx, bson.M{
		"_id": domain.Id,
	}, bson.M{
		"$set": FromDomain(domain),
	})

	if err != nil {
		return bookmarks.Domain{}, err
	}

	// return data
	result, err := br.GetByID(domain.UserID)
	if err != nil {
		return bookmarks.Domain{}, err
	}

	return result, err
}
