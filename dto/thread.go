package dto

import (
	"charum/business/topics"
	"charum/business/users"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ResponseThread struct {
	Id            primitive.ObjectID `json:"_id"`
	Topic         topics.Domain      `json:"topic"`
	Creator       users.Domain       `json:"creator"`
	Title         string             `json:"title"`
	Description   string             `json:"description"`
	Likes         []Like             `json:"likes"`
	TotalLike     int                `json:"totalLike"`
	TotalFollow   int                `json:"totalFollow"`
	TotalComment  int                `json:"totalComment"`
	SuspendStatus string             `json:"suspendStatus,omitempty"`
	SuspendDetail string             `json:"suspendDetail,omitempty"`
	CreatedAt     primitive.DateTime `json:"createdAt"`
	UpdatedAt     primitive.DateTime `json:"updatedAt"`
}

type Like struct {
	User      users.Domain       `json:"user"`
	CreatedAt primitive.DateTime `json:"createdAt"`
}
