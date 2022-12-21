package threads

import (
	"charum/business/threads"
	dtoQuery "charum/dto/query"
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

func (tr *threadRepository) GetManyWithPagination(query dtoQuery.Request, domain *threads.Domain) ([]threads.Domain, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	skip64 := int64(query.Skip)
	limit64 := int64(query.Limit)

	var result []Model
	filter := bson.M{}

	if domain.TopicID != primitive.NilObjectID {
		filter["topicId"] = domain.TopicID
	}

	if domain.Title != "" {
		filter["title"] = bson.M{"$regex": domain.Title}
	}

	cursor, err := tr.collection.Find(ctx, filter, &options.FindOptions{
		Skip:  &skip64,
		Limit: &limit64,
		Sort:  bson.M{query.Sort: query.Order},
	})
	if err != nil {
		return []threads.Domain{}, 0, err
	}

	totalData, err := tr.collection.CountDocuments(ctx, filter)
	if err != nil {
		return []threads.Domain{}, 0, err
	}

	if err = cursor.All(ctx, &result); err != nil {
		return []threads.Domain{}, 0, err
	}

	return ToArrayDomain(result), int(totalData), nil
}

func (tr *threadRepository) GetByID(id primitive.ObjectID) (threads.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result Model
	err := tr.collection.FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&result)
	if err != nil {
		return threads.Domain{}, err
	}

	return result.ToDomain(), nil
}

func (tr *threadRepository) GetAllByTopicID(topicID primitive.ObjectID) ([]threads.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result []Model
	cursor, err := tr.collection.Find(ctx, bson.M{
		"topicId": topicID,
	})
	if err != nil {
		return []threads.Domain{}, err
	}

	if err = cursor.All(ctx, &result); err != nil {
		return []threads.Domain{}, err
	}

	return ToArrayDomain(result), nil
}

func (tr *threadRepository) GetAllByUserID(userID primitive.ObjectID) ([]threads.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result []Model
	cursor, err := tr.collection.Find(ctx, bson.M{
		"creatorId": userID,
	})
	if err != nil {
		return []threads.Domain{}, err
	}

	if err = cursor.All(ctx, &result); err != nil {
		return []threads.Domain{}, err
	}

	return ToArrayDomain(result), nil
}

func (tr *threadRepository) GetLikedByUserID(userID primitive.ObjectID) ([]threads.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result []Model
	cursor, err := tr.collection.Find(ctx, bson.M{
		"likes": bson.M{
			"$elemMatch": bson.M{
				"userID": userID,
			},
		},
	})
	if err != nil {
		return []threads.Domain{}, err
	}

	if err = cursor.All(ctx, &result); err != nil {
		return []threads.Domain{}, err
	}

	return ToArrayDomain(result), nil
}

func (tr *threadRepository) CheckLikedByUserID(userID primitive.ObjectID, threadID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result Model
	err := tr.collection.FindOne(ctx, bson.M{
		"_id": threadID,
		"likes": bson.M{
			"$elemMatch": bson.M{
				"userID": userID,
			},
		},
	}).Decode(&result)
	if err != nil {
		return err
	}

	return nil
}

func (tr *threadRepository) GetAll() ([]threads.Domain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var result []Model
	cursor, err := tr.collection.Find(ctx, bson.M{})
	if err != nil {
		return []threads.Domain{}, err
	}

	if err = cursor.All(ctx, &result); err != nil {
		return []threads.Domain{}, err
	}

	return ToArrayDomain(result), nil
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

func (tr *threadRepository) SuspendByUserID(domain *threads.Domain) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := tr.collection.UpdateMany(ctx, bson.M{
		"creatorId": domain.CreatorID,
	}, bson.M{
		"$set": bson.M{
			"suspendStatus": domain.SuspendStatus,
			"suspendDetail": domain.SuspendDetail,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (tr *threadRepository) AppendLike(userID primitive.ObjectID, threadID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := tr.collection.UpdateOne(ctx, bson.M{
		"_id": threadID,
	}, bson.M{
		"$push": bson.M{
			"likes": bson.M{
				"userID":    userID,
				"timestamp": primitive.NewDateTimeFromTime(time.Now()),
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (tr *threadRepository) RemoveLike(userID primitive.ObjectID, threadID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := tr.collection.UpdateOne(ctx, bson.M{
		"_id": threadID,
	}, bson.M{
		"$pull": bson.M{
			"likes": bson.M{
				"userID": userID,
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (tr *threadRepository) RemoveUserFromAllLikes(userID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := tr.collection.UpdateMany(ctx, bson.M{
		"likes.userID": userID,
	}, bson.M{
		"$pull": bson.M{
			"likes": bson.M{
				"userID": userID,
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
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

func (tr *threadRepository) DeleteAllByUserID(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := tr.collection.DeleteMany(ctx, bson.M{
		"creatorId": id,
	})
	if err != nil {
		return err
	}

	return nil
}

func (tr *threadRepository) DeleteAllByTopicID(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := tr.collection.DeleteMany(ctx, bson.M{
		"TopicId": id,
	})
	if err != nil {
		return err
	}

	return nil
}
