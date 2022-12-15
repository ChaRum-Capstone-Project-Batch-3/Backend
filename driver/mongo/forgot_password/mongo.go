package forgot_password

import (
	"charum/business/forgot_password"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type forgotPasswordRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) forgot_password.Repository {
	return &forgotPasswordRepository{
		collection: db.Collection("forgotPassword"),
	}
}

/*
Create
*/

func (fr *forgotPasswordRepository) Generate(domain *forgot_password.Domain) (forgot_password.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	res, err := fr.collection.InsertOne(ctx, FromDomain(domain))
	if err != nil {
		return forgot_password.Domain{}, err
	}

	result, err := fr.GetByID(res.InsertedID.(primitive.ObjectID))
	if err != nil {
		return forgot_password.Domain{}, err
	}

	return result, nil
}

/*
Read
*/

func (fr *forgotPasswordRepository) GetByID(id primitive.ObjectID) (forgot_password.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result Model
	err := fr.collection.FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&result)
	if err != nil {
		return forgot_password.Domain{}, err
	}

	return result.ToDomain(), nil
}

func (fr *forgotPasswordRepository) GetByToken(token string) (forgot_password.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result Model
	err := fr.collection.FindOne(ctx, bson.M{
		"token": token,
	}).Decode(&result)
	if err != nil {
		return forgot_password.Domain{}, err
	}

	return result.ToDomain(), nil
}

/*
Update
*/

func (fr *forgotPasswordRepository) Update(domain *forgot_password.Domain) (forgot_password.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := fr.collection.UpdateOne(ctx, bson.M{
		"_id": domain.Id,
	}, bson.M{
		"$set": FromDomain(domain),
	})
	if err != nil {
		return forgot_password.Domain{}, err
	}

	result, err := fr.GetByID(domain.Id)
	if err != nil {
		return forgot_password.Domain{}, err
	}

	return result, nil
}

/*
Delete
*/

func (fr *forgotPasswordRepository) Delete(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := fr.collection.DeleteOne(ctx, bson.M{
		"_id": id,
	})
	if err != nil {
		return err
	}

	return nil
}
