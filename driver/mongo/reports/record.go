package reports

import (
	"charum/business/reports"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Model struct {
	Id           primitive.ObjectID `json:"_id" bson:"_id"`
	ReportedID   primitive.ObjectID `json:"reportedId" bson:"reportedId"`
	UserID       primitive.ObjectID `json:"userId" bson:"userId"`
	ReportedType string             `json:"reportedType" bson:"reportedType"`
	ReportDetail string             `json:"reportDetail" bson:"reportDetail"`
	CreatedAt    primitive.DateTime `json:"createdAt" bson:"createdAt"`
	UpdatedAt    primitive.DateTime `json:"updatedAt" bson:"updatedAt"`
}

func FromDomain(domain *reports.Domain) *Model {
	return &Model{
		Id:           domain.Id,
		ReportedID:   domain.ReportedID,
		UserID:       domain.UserID,
		ReportedType: domain.ReportedType,
		ReportDetail: domain.ReportDetail,
		CreatedAt:    domain.CreatedAt,
		UpdatedAt:    domain.UpdatedAt,
	}
}

func (user *Model) ToDomain() reports.Domain {
	return reports.Domain{
		Id:           user.Id,
		ReportedID:   user.ReportedID,
		UserID:       user.UserID,
		ReportedType: user.ReportedType,
		ReportDetail: user.ReportDetail,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}

func ToDomainArray(data []Model) []reports.Domain {
	var result []reports.Domain
	for _, v := range data {
		result = append(result, v.ToDomain())
	}
	return result
}
