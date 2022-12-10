package driver

import (
	commentDomain "charum/business/comments"
	followThreadDomain "charum/business/follow_threads"
	threadDomain "charum/business/threads"
	topicDomain "charum/business/topics"
	userDomain "charum/business/users"
	"charum/driver/cloudinary"
	commentDB "charum/driver/mongo/comments"
	followThreadDB "charum/driver/mongo/follow_threads"
	threadDB "charum/driver/mongo/threads"
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

func NewThreadRepository(db *mongo.Database) threadDomain.Repository {
	return threadDB.NewMongoRepository(db)
}

func NewCommentRepository(db *mongo.Database) commentDomain.Repository {
	return commentDB.NewMongoRepository(db)
}

func NewFollowThreadRepository(db *mongo.Database) followThreadDomain.Repository {
	return followThreadDB.NewMongoRepository(db)
}

func NewCloudinaryRepository() cloudinary.Function {
	return cloudinary.NewCloudinaryRepository()
}
