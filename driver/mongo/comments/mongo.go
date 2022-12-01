package comments

import (
	"charum/business/comments"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type commentRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) comments.Repository {
	return &commentRepository{
		collection: db.Collection("comments"),
	}
}

/*
Create
*/

func (cr *commentRepository) Create(domain *comments.Domain) (comments.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	res, err := cr.collection.InsertOne(ctx, FromDomain(domain))
	if err != nil {
		return comments.Domain{}, err
	}

	result, err := cr.GetByID(res.InsertedID.(primitive.ObjectID))
	if err != nil {
		return comments.Domain{}, err
	}

	return result, err
}

/*
Read
*/

func (cr *commentRepository) GetByID(id primitive.ObjectID) (comments.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result Model
	err := cr.collection.FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&result)
	if err != nil {
		return comments.Domain{}, err
	}

	return result.ToDomain(), nil
}

func (cr *commentRepository) GetByThreadID(threadID primitive.ObjectID) ([]comments.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result []Model
	// get commment by thread id sorted by createdAt descending
	cursor, err := cr.collection.Find(ctx, bson.M{
		"threadID": threadID,
	}, &options.FindOptions{
		Sort: bson.M{
			"createdAt": -1,
		},
	})
	if err != nil {
		return []comments.Domain{}, err
	}

	if err = cursor.All(ctx, &result); err != nil {
		return []comments.Domain{}, err
	}

	return ToDomainArray(result), nil
}

/*
Update
*/

func (cr *commentRepository) Update(domain *comments.Domain) (comments.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := cr.collection.UpdateOne(ctx, bson.M{
		"_id": domain.Id,
	}, bson.M{
		"$set": FromDomain(domain),
	})
	if err != nil {
		return comments.Domain{}, err
	}

	result, err := cr.GetByID(domain.Id)
	if err != nil {
		return comments.Domain{}, err
	}

	return result, nil
}

/*
Delete
*/

func (cr *commentRepository) Delete(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := cr.collection.DeleteOne(ctx, bson.M{
		"_id": id,
	})
	if err != nil {
		return err
	}

	return nil
}
