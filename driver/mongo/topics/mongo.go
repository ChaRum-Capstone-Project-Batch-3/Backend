package topics

import (
	"charum/business/topics"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type topicRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) topics.Repository {
	return &topicRepository{
		collection: db.Collection("topics"),
	}
}

/*
Create
*/

func (tr *topicRepository) CreateTopic(domain *topics.Domain) (topics.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	res, err := tr.collection.InsertOne(ctx, FromDomain(domain))
	if err != nil {
		return topics.Domain{}, err
	}
	result, err := tr.GetByID(res.InsertedID.(primitive.ObjectID))
	if err != nil {
		return topics.Domain{}, err
	}

	return result, err
}

/*
Read
*/

func (tr *topicRepository) GetByID(id primitive.ObjectID) (topics.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result Model
	err := tr.collection.FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&result)
	if err != nil {
		return topics.Domain{}, err
	}

	return result.ToDomain(), nil
}

/*
Update
*/

func (tr *topicRepository) UpdateTopic(domain *topics.Domain) (topics.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := tr.collection.UpdateOne(ctx, bson.M{
		"_id": domain.Id,
	}, bson.M{
		"$set": FromDomain(domain),
	})
	if err != nil {
		return topics.Domain{}, err
	}

	result, err := tr.GetByID(domain.Id)
	if err != nil {
		return topics.Domain{}, err
	}

	return result, nil
}

/*
Delete
*/

func (tr *topicRepository) DeleteTopic(id primitive.ObjectID) error {
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
