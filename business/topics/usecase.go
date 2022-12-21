package topics

import (
	dtoPagination "charum/dto/pagination"
	dtoQuery "charum/dto/query"
	_cloudinary "charum/helper/cloudinary"
	"charum/util"
	"errors"
	"math"
	"mime/multipart"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TopicUseCase struct {
	topicsRepository Repository
	cloudinary       _cloudinary.Function
}

func NewTopicUseCase(tr Repository, cld _cloudinary.Function) UseCase {
	return &TopicUseCase{
		topicsRepository: tr,
		cloudinary:       cld,
	}
}

/*
Create
*/

func (tu *TopicUseCase) Create(domain *Domain, image *multipart.FileHeader) (Domain, error) {
	_, err := tu.topicsRepository.GetByTopic(domain.Topic)
	if err == nil {
		return Domain{}, errors.New("topic already exist")
	}

	if image != nil {
		cloudinaryURL, err := tu.cloudinary.Upload("topic", image, util.GenerateUUID())
		if err != nil {
			return Domain{}, errors.New("failed to upload image")
		}

		domain.ImageURL = cloudinaryURL
	}

	domain.Id = primitive.NewObjectID()
	domain.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	domain.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	result, err := tu.topicsRepository.Create(domain)
	if err != nil {
		return Domain{}, errors.New("failed to create topic")
	}
	return result, nil
}

/*
Read
*/

func (tu *TopicUseCase) GetByID(id primitive.ObjectID) (Domain, error) {
	result, err := tu.topicsRepository.GetByID(id)
	if err != nil {
		return Domain{}, errors.New("failed to get topic")
	}
	return result, nil
}

func (tu *TopicUseCase) GetManyWithPagination(pagination dtoPagination.Request, domain *Domain) ([]Domain, int, int, error) {
	skip := pagination.Limit * (pagination.Page - 1)
	var orderInMongo int

	if pagination.Order == "asc" {
		orderInMongo = 1
	} else {
		orderInMongo = -1
	}

	query := dtoQuery.Request{
		Skip:  skip,
		Limit: pagination.Limit,
		Order: orderInMongo,
		Sort:  pagination.Sort,
	}

	users, totalData, err := tu.topicsRepository.GetManyWithPagination(query, domain)
	if err != nil {
		return []Domain{}, 0, 0, errors.New("failed to get topics")
	}

	totalPage := math.Ceil(float64(totalData) / float64(pagination.Limit))
	return users, int(totalPage), totalData, nil
}

func (tu *TopicUseCase) GetByTopic(topic string) (Domain, error) {
	result, err := tu.topicsRepository.GetByTopic(topic)
	if err != nil {
		return Domain{}, errors.New("failed to get topic")
	}
	return result, nil
}

/*
Update
*/

func (tu *TopicUseCase) Update(domain *Domain, image *multipart.FileHeader) (Domain, error) {
	result, err := tu.topicsRepository.GetByID(domain.Id)
	if err != nil {
		return Domain{}, errors.New("failed to get topic")
	}

	if domain.Topic != result.Topic {
		_, err := tu.topicsRepository.GetByTopic(domain.Topic)
		if err == nil {
			return Domain{}, errors.New("topic already exist")
		}
	}

	if image != nil {
		if result.ImageURL != "" {
			err = tu.cloudinary.Delete("topic", util.GetFilenameWithoutExtension(result.ImageURL))
			if err != nil {
				return Domain{}, errors.New("failed to delete image")
			}
		}

		cloudinaryURL, err := tu.cloudinary.Upload("topic", image, util.GenerateUUID())
		if err != nil {
			return Domain{}, err
		}

		result.ImageURL = cloudinaryURL
	}

	result.Topic = domain.Topic
	result.Description = domain.Description
	result.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	updatedResult, err := tu.topicsRepository.Update(&result)
	if err != nil {
		return Domain{}, errors.New("failed to update topic")
	}
	return updatedResult, nil
}

/*
Delete
*/

func (tu *TopicUseCase) Delete(id primitive.ObjectID) (Domain, error) {
	result, err := tu.topicsRepository.GetByID(id)
	if err != nil {
		return Domain{}, errors.New("failed to get topic")
	}

	if result.ImageURL != "" {
		err = tu.cloudinary.Delete("topic", util.GetFilenameWithoutExtension(result.ImageURL))
		if err != nil {
			return Domain{}, errors.New("failed to delete image")
		}
	}

	err = tu.topicsRepository.Delete(id)
	if err != nil {
		return Domain{}, errors.New("failed to delete topic")
	}

	return result, nil
}
