package reports

import "go.mongodb.org/mongo-driver/bson/primitive"

type Domain struct {
	Id           primitive.ObjectID `json:"id"`
	UserID       primitive.ObjectID `json:"userId"`
	ReportedID   primitive.ObjectID `json:"reportedID"`
	ReportType   string             `json:"reportType"`
	ReportDetail string             `json:"reportDetail"`
	CreatedAt    primitive.DateTime `json:"createdAt"`
	UpdatedAt    primitive.DateTime `json:"updatedAt"`
}

type Repository interface {
	// Create
	Create(domain *Domain) (Domain, error)
	// Read
	GetByReportedID(id primitive.ObjectID) ([]Domain, error)
	CheckByUserID(userID primitive.ObjectID, reportedID primitive.ObjectID) (Domain, error)
}

type UseCase interface {
	// Create
	Create(domain *Domain) (Domain, error)
	// Read
	GetByReportedID(id primitive.ObjectID) ([]Domain, error)
}
