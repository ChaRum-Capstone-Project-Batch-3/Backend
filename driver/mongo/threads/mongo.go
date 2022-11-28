package threads

import (
	"charum/business/threads"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type threadRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) threads.Repository {
	return &threadRepository{
		collection: db.Collection("threads"),
	}
}

/*
Create
*/

func (tr *threadRepository) Create(domain *threads.Domain) (threads.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	res, err := tr.collection.InsertOne(ctx, FromDomain(domain))
	if err != nil {
		return threads.Domain{}, err
	}

	result, err := tr.GetByID(res.InsertedID.(primitive.ObjectID))
	if err != nil {
		return threads.Domain{}, err
	}

	return result, err
}

/*
Read
*/

func (tr *threadRepository) GetWithSortAndOrder(skip int, limit int, sort string, order int) ([]threads.Domain, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	skip64 := int64(skip)
	limit64 := int64(limit)

	var result []Model

	cursor, err := tr.collection.Find(ctx, bson.M{}, &options.FindOptions{
		Skip:  &skip64,
		Limit: &limit64,
		Sort:  bson.M{sort: order},
	})
	if err != nil {
		return []threads.Domain{}, 0, err
	}

	// count total data in collection
	totalData, err := tr.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return []threads.Domain{}, 0, err
	}

	if err = cursor.All(ctx, &result); err != nil {
		return []threads.Domain{}, 0, err
	}

	return ToArrayDomain(result), int(totalData), nil
}

func (ur *threadRepository) GetByID(id primitive.ObjectID) (threads.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result Model
	err := ur.collection.FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&result)
	if err != nil {
		return threads.Domain{}, err
	}

	return result.ToDomain(), nil
}

/*
Update
*/

func (tr *threadRepository) Update(domain *threads.Domain) (threads.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := tr.collection.UpdateOne(ctx, bson.M{
		"_id": domain.Id,
	}, bson.M{
		"$set": FromDomain(domain),
	})
	if err != nil {
		return threads.Domain{}, err
	}

	result, err := tr.GetByID(domain.Id)
	if err != nil {
		return threads.Domain{}, err
	}

	return result, nil
}

/*
Delete
*/

func (tr *threadRepository) Delete(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := tr.collection.DeleteOne(ctx, bson.M{
		"_id": id,
	})
	if err != nil {
		return err
	}

	return nil
}
