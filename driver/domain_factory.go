package driver

import (
	topicDomain "charum/business/topics"
	userDomain "charum/business/users"
	topicDB "charum/driver/mongo/topics"
	userDB "charum/driver/mongo/users"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewUserRepository(db *mongo.Database) userDomain.Repository {
	return userDB.NewMongoRepository(db)
}

func NewTopicRepository(db *mongo.Database) topicDomain.Repository {
	return topicDB.NewMongoRepository(db)
}
