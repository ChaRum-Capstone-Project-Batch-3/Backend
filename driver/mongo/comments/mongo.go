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

func (cr *commentRepository) CountByThreadID(threadID primitive.ObjectID) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	count, err := cr.collection.CountDocuments(ctx, bson.M{
		"threadID": threadID,
	})
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (cr *commentRepository) GetByIDAndThreadID(id primitive.ObjectID, threadID primitive.ObjectID) (comments.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result Model
	err := cr.collection.FindOne(ctx, bson.M{
		"_id":      id,
		"threadID": threadID,
	}).Decode(&result)
	if err != nil {
		return comments.Domain{}, err
	}

	return result.ToDomain(), nil
}

func (cr *commentRepository) GetAllByUserID(userID primitive.ObjectID) ([]comments.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result []Model
	cursor, err := cr.collection.Find(ctx, bson.M{
		"userID": userID,
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

func (cr *commentRepository) DeleteAllByUserID(userID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := cr.collection.DeleteMany(ctx, bson.M{
		"userID": userID,
	})
	if err != nil {
		return err
	}

	return nil
}

func (cr *commentRepository) DeleteAllByThreadID(threadID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := cr.collection.DeleteMany(ctx, bson.M{
		"threadID": threadID,
	})
	if err != nil {
		return err
	}

	return nil
}
