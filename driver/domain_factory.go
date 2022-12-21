package driver

import (
	bookmarkDomain "charum/business/bookmarks"
	commentDomain "charum/business/comments"
	followThreadDomain "charum/business/follow_threads"
	forgotPasswordDomain "charum/business/forgot_password"
	reportDomain "charum/business/reports"
	threadDomain "charum/business/threads"
	topicDomain "charum/business/topics"
	userDomain "charum/business/users"

	bookmarkDB "charum/driver/mongo/bookmarks"
	commentDB "charum/driver/mongo/comments"
	followThreadDB "charum/driver/mongo/follow_threads"
	forgotPasswordDB "charum/driver/mongo/forgot_password"
	reportDB "charum/driver/mongo/reports"
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

func NewBookmarkRepository(db *mongo.Database) bookmarkDomain.Repository {
	return bookmarkDB.NewMongoRepository(db)
}

func NewForgotPasswordRepository(db *mongo.Database) forgotPasswordDomain.Repository {
	return forgotPasswordDB.NewMongoRepository(db)
}

func NewReportRepository(db *mongo.Database) reportDomain.Repository {
	return reportDB.NewMongoRepository(db)
}
