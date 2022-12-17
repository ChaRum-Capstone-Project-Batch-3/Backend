package reports

import (
	"charum/business/threads"
	"charum/business/users"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ReportUseCase struct {
	reportRepository Repository
	userRepository   users.Repository
	threadRepository threads.Repository
}

func NewReportUseCase(rr Repository, ur users.Repository, tr threads.Repository) UseCase {
	return &ReportUseCase{
		reportRepository: rr,
		userRepository:   ur,
		threadRepository: tr,
	}
}

/*
Create
*/

func (ru *ReportUseCase) Create(domain *Domain) (Domain, error) {
	// check ReportedID if exist in users or threads ID
	reportType, err := ru.CheckID(domain.ReportedID)
	if err != nil {
		return Domain{}, errors.New("ID not found")
	}
	_, err = ru.reportRepository.CheckByUserID(domain.UserID, domain.ReportedID)
	if err == nil {
		return Domain{}, errors.New("already reported")
	}

	domain.Id = primitive.NewObjectID()
	domain.ReportType = reportType
	domain.ReportDetail = "Inappropriate content or behavior"
	domain.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	domain.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	report, err := ru.reportRepository.Create(domain)
	if err != nil {
		return Domain{}, err
	}

	return report, nil
}

/*
Read
*/

func (ru *ReportUseCase) CheckID(id primitive.ObjectID) (string, error) {
	_, err := ru.userRepository.GetByID(id)
	if err == nil {
		return "user", nil
	} else {
		_, err := ru.threadRepository.GetByID(id)
		if err == nil {
			return "thread", nil
		} else {
			return "", err
		}
	}
}

func (ru *ReportUseCase) GetByReportedID(id primitive.ObjectID) ([]Domain, error) {
	report, err := ru.reportRepository.GetByReportedID(id)
	if err != nil {
		return []Domain{}, err
	}

	return report, nil
}
