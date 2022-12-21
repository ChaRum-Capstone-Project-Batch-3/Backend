package users

import (
	"charum/business/users"
	dtoQuery "charum/dto/query"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) users.Repository {
	return &userRepository{
		collection: db.Collection("users"),
	}
}

/*
Create
*/

func (ur *userRepository) Create(domain *users.Domain) (users.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	res, err := ur.collection.InsertOne(ctx, FromDomain(domain))
	if err != nil {
		return users.Domain{}, err
	}

	result, err := ur.GetByID(res.InsertedID.(primitive.ObjectID))
	if err != nil {
		return users.Domain{}, err
	}

	return result, err
}

/*
Read
*/

func (ur *userRepository) GetByID(id primitive.ObjectID) (users.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result Model
	err := ur.collection.FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&result)
	if err != nil {
		return users.Domain{}, err
	}

	return result.ToDomain(), nil
}

func (ur *userRepository) GetByEmail(email string) (users.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result Model
	err := ur.collection.FindOne(ctx, bson.M{
		"email": email,
	}).Decode(&result)

	return result.ToDomain(), err
}

func (ur *userRepository) GetByUsername(username string) (users.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result Model
	err := ur.collection.FindOne(ctx, bson.M{
		"userName": username,
	}).Decode(&result)

	return result.ToDomain(), err
}

func (ur *userRepository) GetManyWithPagination(query dtoQuery.Request, domain *users.Domain) ([]users.Domain, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	skip64 := int64(query.Skip)
	limit64 := int64(query.Limit)

	var result []Model
	filter := bson.M{}

	if domain.Email != "" {
		filter["email"] = bson.M{"$regex": domain.Email}
	}

	if domain.UserName != "" {
		filter["userName"] = bson.M{"$regex": domain.UserName}
	}

	if domain.DisplayName != "" {
		filter["displayName"] = bson.M{"$regex": domain.DisplayName}
	}

	cursor, err := ur.collection.Find(ctx, filter, &options.FindOptions{
		Skip:  &skip64,
		Limit: &limit64,
		Sort:  bson.M{query.Sort: query.Order},
	})
	if err != nil {
		return []users.Domain{}, 0, err
	}

	// count total data in collection
	totalData, err := ur.collection.CountDocuments(ctx, filter)
	if err != nil {
		return []users.Domain{}, 0, err
	}

	if err = cursor.All(ctx, &result); err != nil {
		return []users.Domain{}, 0, err
	}

	return ToArrayDomain(result), int(totalData), nil
}

func (ur *userRepository) GetAll() ([]users.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result []Model
	cursor, err := ur.collection.Find(ctx, bson.M{})
	if err != nil {
		return []users.Domain{}, err
	}

	if err = cursor.All(ctx, &result); err != nil {
		return []users.Domain{}, err
	}

	return ToArrayDomain(result), nil
}

/*
Update
*/

func (ur *userRepository) Update(domain *users.Domain) (users.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := ur.collection.UpdateOne(ctx, bson.M{
		"_id": domain.Id,
	}, bson.M{
		"$set": FromDomain(domain),
	})
	if err != nil {
		return users.Domain{}, err
	}

	result, err := ur.GetByID(domain.Id)
	if err != nil {
		return users.Domain{}, err
	}

	return result, nil
}

func (ur *userRepository) UpdatePassword(domain *users.Domain) (users.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := ur.collection.UpdateOne(ctx, bson.M{
		"_id": domain.Id,
	}, bson.M{
		"$set": bson.M{
			"password":  domain.Password,
			"updatedAt": domain.UpdatedAt,
		},
	})
	if err != nil {
		return users.Domain{}, err
	}

	result, err := ur.GetByID(domain.Id)
	if err != nil {
		return users.Domain{}, err
	}

	return result, nil
}

/*
Delete
*/

func (ur *userRepository) Delete(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := ur.collection.DeleteOne(ctx, bson.M{
		"_id": id,
	})
	if err != nil {
		return err
	}

	return nil
}
