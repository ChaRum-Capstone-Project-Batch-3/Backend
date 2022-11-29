package driver

import (
	bookmarkDomain "charum/business/bookmarks"
	commentDomain "charum/business/comments"
	threadDomain "charum/business/threads"
	topicDomain "charum/business/topics"
	userDomain "charum/business/users"

	bookmarkDB "charum/driver/mongo/bookmark"
	commentDB "charum/driver/mongo/comments"
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

func NewBookmarkRepository(db *mongo.Database) bookmarkDomain.Repository {
	return bookmarkDB.NewMongoRepository(db)
}

func NewCommentRepository(db *mongo.Database) commentDomain.Repository {
	return commentDB.NewMongoRepository(db)
}
