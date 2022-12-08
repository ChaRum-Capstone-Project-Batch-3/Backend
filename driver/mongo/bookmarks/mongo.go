package bookmark

import (
	"charum/business/bookmarks"
	"context"
	"fmt"
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

	_, err := br.collection.InsertOne(ctx, FromDomain(domain))
	if err != nil {
		return bookmarks.Domain{}, err
	}

	// get threadid from res
	result, err := br.GetByID(domain.UserID, domain.ThreadID)
	fmt.Println(result)
	fmt.Println(domain.UserID, domain.ThreadID)
	if err != nil {
		return bookmarks.Domain{}, err
	}

	return result, err
}

/*
Read
*/

func (br *bookmarkRepository) GetByID(UserID primitive.ObjectID, ThreadID primitive.ObjectID) (bookmarks.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result Bookmark
	err := br.collection.FindOne(ctx, bson.M{
		"userId":   UserID,
		"threadId": ThreadID,
	}).Decode(&result)
	if err != nil {
		return bookmarks.Domain{}, err
	}

	return result.ToDomain(), nil
}

func (br *bookmarkRepository) GetAllBookmark(UserID primitive.ObjectID) ([]bookmarks.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// get all bookmark by userID
	cursor, err := br.collection.Find(ctx, bson.M{
		"userId": UserID,
	})
	if err != nil {
		return []bookmarks.Domain{}, err
	}

	// convert to array
	var results []Bookmark
	if err = cursor.All(ctx, &results); err != nil {
		return []bookmarks.Domain{}, err
	}

	// convert to domain
	var domains []bookmarks.Domain
	for _, result := range results {
		domains = append(domains, result.ToDomain())
	}

	return domains, nil
}

/*
Update
*/

func (br *bookmarkRepository) UpdateBookmark(userID primitive.ObjectID, threadID primitive.ObjectID, domain *bookmarks.Domain) (bookmarks.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// update
	// only update threads & updatedAt
	_, err := br.collection.UpdateOne(ctx, bson.M{
		"_id": userID,
	}, bson.M{
		"$set": FromDomain(domain),
	})

	if err != nil {
		return bookmarks.Domain{}, err
	}

	// return data
	result, err := br.GetByID(domain.UserID, domain.ThreadID)
	if err != nil {
		return bookmarks.Domain{}, err
	}

	return result, err
}

/*
Delete
*/

func (br *bookmarkRepository) DeleteBookmark(userID primitive.ObjectID, threadID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// delete
	_, err := br.collection.DeleteOne(ctx, bson.M{
		"userId":   userID,
		"threadId": threadID,
	})

	if err != nil {
		return err
	}

	return nil
}
