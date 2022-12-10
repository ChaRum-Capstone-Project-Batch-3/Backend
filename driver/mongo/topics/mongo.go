package topics

import (
	"charum/business/topics"
	dtoQuery "charum/dto/query"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (tr *topicRepository) Create(domain *topics.Domain) (topics.Domain, error) {
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

func (tr *topicRepository) GetManyWithPagination(query dtoQuery.Request, domain *topics.Domain) ([]topics.Domain, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	skip64 := int64(query.Skip)
	limit64 := int64(query.Limit)

	var result []Model
	filter := bson.M{}

	if domain.Topic != "" {
		filter["topic"] = bson.M{"$regex": domain.Topic}
	}

	cursor, err := tr.collection.Find(ctx, filter, &options.FindOptions{
		Skip:  &skip64,
		Limit: &limit64,
		Sort:  bson.M{query.Sort: query.Order},
	})
	if err != nil {
		return []topics.Domain{}, 0, err
	}

	totalData, err := tr.collection.CountDocuments(ctx, filter)
	if err != nil {
		return []topics.Domain{}, 0, err
	}

	if err = cursor.All(ctx, &result); err != nil {
		return []topics.Domain{}, 0, err
	}

	return ToArrayDomain(result), int(totalData), nil
}

func (tr *topicRepository) GetByTopic(topic string) (topics.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result Model
	err := tr.collection.FindOne(ctx, bson.M{
		"topic": topic,
	}).Decode(&result)
	if err != nil {
		return topics.Domain{}, err
	}

	return result.ToDomain(), nil
}

/*
Update
*/

func (tr *topicRepository) Update(domain *topics.Domain) (topics.Domain, error) {
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

func (tr *topicRepository) Delete(id primitive.ObjectID) error {
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
