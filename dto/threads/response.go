package threads

import (
	"charum/business/topics"
	"charum/business/users"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Response struct {
	Id            primitive.ObjectID `json:"_id"`
	Topic         topics.Domain      `json:"topic"`
	Creator       users.Domain       `json:"creator"`
	Title         string             `json:"title"`
	Description   string             `json:"description"`
	Likes         []Like             `json:"likes"`
	ImageURL      string             `json:"imageURL"`
	IsLiked       bool               `json:"isLiked"`
	IsBookmarked  bool               `json:"isBookmarked"`
	IsFollowed    bool               `json:"isFollowed"`
	TotalLike     int                `json:"totalLike"`
	TotalFollow   int                `json:"totalFollow"`
	TotalComment  int                `json:"totalComment"`
	TotalBookmark int                `json:"totalBookmark"`
	TotalReported int                `json:"totalReported"`
	SuspendStatus string             `json:"suspendStatus,omitempty"`
	SuspendDetail string             `json:"suspendDetail,omitempty"`
	CreatedAt     primitive.DateTime `json:"createdAt"`
	UpdatedAt     primitive.DateTime `json:"updatedAt"`
}

type Like struct {
	User      users.Domain       `json:"user"`
	Timestamp primitive.DateTime `json:"timestamp" bson:"timestamp"`
}
