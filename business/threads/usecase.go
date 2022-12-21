package threads

import (
	"charum/business/topics"
	"charum/business/users"
	dtoPagination "charum/dto/pagination"
	dtoQuery "charum/dto/query"
	dtoThread "charum/dto/threads"
	"charum/helper/cloudinary"
	"charum/util"
	"errors"
	"math"
	"mime/multipart"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ThreadUseCase struct {
	threadRepository Repository
	topicRepository  topics.Repository
	userRepository   users.Repository
	cloudinary       cloudinary.Function
}

func NewThreadUseCase(thr Repository, tor topics.Repository, ur users.Repository, c cloudinary.Function) UseCase {
	return &ThreadUseCase{
		threadRepository: thr,
		topicRepository:  tor,
		userRepository:   ur,
		cloudinary:       c,
	}
}

/*
Create
*/

func (tu *ThreadUseCase) Create(domain *Domain, image *multipart.FileHeader) (Domain, error) {
	_, err := tu.topicRepository.GetByID(domain.TopicID)
	if err != nil {
		return Domain{}, errors.New("failed to get topic")
	}

	if image != nil {
		cloudinaryURL, err := tu.cloudinary.Upload("thread", image, util.GenerateUUID())
		if err != nil {
			return Domain{}, errors.New("failed to upload image")
		}

		domain.ImageURL = cloudinaryURL
	}

	domain.Id = primitive.NewObjectID()
	domain.Likes = []Like{}
	domain.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	domain.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	thread, err := tu.threadRepository.Create(domain)
	if err != nil {
		if domain.ImageURL != "" {
			delErr := tu.cloudinary.Delete("thread", util.GetFilenameWithoutExtension(domain.ImageURL))
			if delErr != nil {
				return Domain{}, errors.New("failed to delete image")
			}
		}

		return Domain{}, errors.New("failed to create thread")
	}

	return thread, nil
}

/*
Read
*/

func (tu *ThreadUseCase) GetManyWithPagination(pagination dtoPagination.Request, domain *Domain) ([]Domain, int, int, error) {
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

	if domain.TopicID != primitive.NilObjectID {
		_, err := tu.topicRepository.GetByID(domain.TopicID)
		if err != nil {
			return []Domain{}, 0, 0, errors.New("failed to get topic")
		}
	}

	threads, totalData, err := tu.threadRepository.GetManyWithPagination(query, domain)
	if err != nil {
		return []Domain{}, 0, 0, errors.New("failed to get threads")
	}

	totalPage := math.Ceil(float64(totalData) / float64(pagination.Limit))

	return threads, int(totalPage), totalData, nil
}

func (tu *ThreadUseCase) GetByID(id primitive.ObjectID) (Domain, error) {
	thread, err := tu.threadRepository.GetByID(id)
	if err != nil {
		return Domain{}, errors.New("failed to get thread")
	}

	return thread, nil
}

func (tu *ThreadUseCase) GetAllByTopicID(topicID primitive.ObjectID) ([]Domain, error) {
	threads, err := tu.threadRepository.GetAllByTopicID(topicID)
	if err != nil {
		return []Domain{}, errors.New("failed to get threads")
	}

	return threads, nil
}

func (tu *ThreadUseCase) GetAllByUserID(userID primitive.ObjectID) ([]Domain, error) {
	threads, err := tu.threadRepository.GetAllByUserID(userID)
	if err != nil {
		return []Domain{}, errors.New("failed to get threads")
	}

	return threads, nil
}

func (tu *ThreadUseCase) GetLikedByUserID(userID primitive.ObjectID) ([]Domain, error) {
	threads, err := tu.threadRepository.GetLikedByUserID(userID)
	if err != nil {
		return []Domain{}, errors.New("failed to get liked threads")
	}

	return threads, nil
}

// get all
func (tu *ThreadUseCase) GetAll() (int, error) {
	threads, err := tu.threadRepository.GetAll()
	if err != nil {
		return 0, errors.New("failed to get threads")
	}

	return len(threads), nil
}

func (tu *ThreadUseCase) DomainToResponse(domain Domain, userID primitive.ObjectID) (dtoThread.Response, error) {
	creator, err := tu.userRepository.GetByID(domain.CreatorID)
	if err != nil {
		return dtoThread.Response{}, errors.New("failed to get creator")
	}

	topic, err := tu.topicRepository.GetByID(domain.TopicID)
	if err != nil {
		return dtoThread.Response{}, errors.New("failed to get topic")
	}

	likes := []dtoThread.Like{}
	for _, like := range domain.Likes {
		user, err := tu.userRepository.GetByID(like.UserID)
		if err != nil {
			return dtoThread.Response{}, err
		}

		likes = append(likes, dtoThread.Like{
			User:      user,
			Timestamp: like.Timestamp,
		})
	}

	isLiked := false

	if userID != primitive.NilObjectID {
		for _, like := range domain.Likes {
			if like.UserID == userID {
				isLiked = true
				break
			}
		}
	}

	return dtoThread.Response{
		Id:            domain.Id,
		Topic:         topic,
		Creator:       creator,
		Title:         domain.Title,
		Description:   domain.Description,
		Likes:         likes,
		TotalLike:     len(domain.Likes),
		IsLiked:       isLiked,
		IsBookmarked:  false,
		IsFollowed:    false,
		ImageURL:      domain.ImageURL,
		SuspendStatus: domain.SuspendStatus,
		SuspendDetail: domain.SuspendDetail,
		CreatedAt:     domain.CreatedAt,
		UpdatedAt:     domain.UpdatedAt,
	}, nil
}

func (tu *ThreadUseCase) DomainsToResponseArray(domains []Domain, userID primitive.ObjectID) ([]dtoThread.Response, error) {
	var responses []dtoThread.Response
	for _, domain := range domains {
		response, err := tu.DomainToResponse(domain, userID)
		if err != nil {
			return []dtoThread.Response{}, errors.New("failed to get thread")
		}

		responses = append(responses, response)
	}

	return responses, nil
}

/*
Update
*/

func (tu *ThreadUseCase) UserUpdate(domain *Domain, image *multipart.FileHeader) (Domain, error) {
	_, err := tu.topicRepository.GetByID(domain.TopicID)
	if err != nil {
		return Domain{}, errors.New("failed to get topic")
	}

	thread, err := tu.threadRepository.GetByID(domain.Id)
	if err != nil {
		return Domain{}, errors.New("failed to get thread")
	}

	if thread.CreatorID != domain.CreatorID {
		return Domain{}, errors.New("user are not the thread creator")
	}

	if image != nil {
		if thread.ImageURL != "" {
			delErr := tu.cloudinary.Delete("thread", util.GetFilenameWithoutExtension(thread.ImageURL))
			if delErr != nil {
				return Domain{}, errors.New("failed to delete image")
			}
		}

		cloudinaryURL, err := tu.cloudinary.Upload("thread", image, util.GenerateUUID())
		if err != nil {
			return Domain{}, err
		}

		thread.ImageURL = cloudinaryURL
	}

	thread.TopicID = domain.TopicID
	thread.Title = domain.Title
	thread.Description = domain.Description
	thread.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	updatedThread, err := tu.threadRepository.Update(&thread)
	if err != nil {
		return Domain{}, errors.New("failed to update thread")
	}

	return updatedThread, nil
}

func (tu *ThreadUseCase) AdminUpdate(domain *Domain, image *multipart.FileHeader) (Domain, error) {
	_, err := tu.topicRepository.GetByID(domain.TopicID)
	if err != nil {
		return Domain{}, errors.New("failed to get topic")
	}

	thread, err := tu.threadRepository.GetByID(domain.Id)
	if err != nil {
		return Domain{}, errors.New("failed to get thread")
	}

	if image != nil {
		if thread.ImageURL != "" {
			delErr := tu.cloudinary.Delete("thread", util.GetFilenameWithoutExtension(thread.ImageURL))
			if delErr != nil {
				return Domain{}, errors.New("failed to delete image")
			}
		}

		cloudinaryURL, err := tu.cloudinary.Upload("thread", image, util.GenerateUUID())
		if err != nil {
			return Domain{}, err
		}

		thread.ImageURL = cloudinaryURL
	}

	thread.TopicID = domain.TopicID
	thread.Title = domain.Title
	thread.Description = domain.Description
	thread.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	updatedThread, err := tu.threadRepository.Update(&thread)
	if err != nil {
		return Domain{}, errors.New("failed to update thread")
	}

	return updatedThread, nil
}

func (tu *ThreadUseCase) SuspendByUserID(userID primitive.ObjectID) error {
	domain := Domain{
		CreatorID:     userID,
		SuspendStatus: "user suspend",
		SuspendDetail: "user is violate the rules",
	}

	err := tu.threadRepository.SuspendByUserID(&domain)
	if err != nil {
		return errors.New("failed to suspend user threads")
	}

	return nil
}

func (tu *ThreadUseCase) Like(userID primitive.ObjectID, threadID primitive.ObjectID) error {
	_, err := tu.threadRepository.GetByID(threadID)
	if err != nil {
		return errors.New("failed to get thread")
	}

	err = tu.threadRepository.CheckLikedByUserID(userID, threadID)
	if err == nil {
		return errors.New("user already like this thread")
	}

	err = tu.threadRepository.AppendLike(userID, threadID)
	if err != nil {
		return errors.New("failed to like thread")
	}

	return nil
}

func (tu *ThreadUseCase) Unlike(userID primitive.ObjectID, threadID primitive.ObjectID) error {
	_, err := tu.threadRepository.GetByID(threadID)
	if err != nil {
		return errors.New("failed to get thread")
	}

	err = tu.threadRepository.CheckLikedByUserID(userID, threadID)
	if err != nil {
		return errors.New("user not like this thread")
	}

	err = tu.threadRepository.RemoveLike(userID, threadID)
	if err != nil {
		return errors.New("failed to unlike thread")
	}

	return nil
}

func (tu *ThreadUseCase) RemoveUserFromAllLikes(userID primitive.ObjectID) error {
	err := tu.threadRepository.RemoveUserFromAllLikes(userID)
	if err != nil {
		return errors.New("failed to remove user from all likes")
	}

	return nil
}

/*
Delete
*/

func (tu *ThreadUseCase) Delete(userID primitive.ObjectID, threadID primitive.ObjectID) (Domain, error) {
	thread, err := tu.threadRepository.GetByID(threadID)
	if err != nil {
		return Domain{}, errors.New("failed to get thread")
	}

	if thread.CreatorID != userID {
		return Domain{}, errors.New("user are not the thread creator")
	}

	if thread.ImageURL != "" {
		delErr := tu.cloudinary.Delete("thread", util.GetFilenameWithoutExtension(thread.ImageURL))
		if delErr != nil {
			return Domain{}, errors.New("failed to delete image")
		}
	}

	err = tu.threadRepository.Delete(threadID)
	if err != nil {
		return Domain{}, errors.New("failed to delete thread")
	}

	return thread, nil
}

func (tu *ThreadUseCase) DeleteAllByUserID(userID primitive.ObjectID) error {
	threads, err := tu.threadRepository.GetAllByUserID(userID)
	if err != nil {
		return errors.New("failed to get user threads")
	}

	for _, thread := range threads {
		if thread.ImageURL != "" {
			delErr := tu.cloudinary.Delete("thread", util.GetFilenameWithoutExtension(thread.ImageURL))
			if delErr != nil {
				return errors.New("failed to delete image")
			}
		}
	}

	err = tu.threadRepository.DeleteAllByUserID(userID)
	if err != nil {
		return errors.New("failed to delete user threads")
	}

	return nil
}

func (tu *ThreadUseCase) DeleteByThreadID(threadID primitive.ObjectID) error {
	thread, err := tu.threadRepository.GetByID(threadID)
	if err != nil {
		return errors.New("failed to get thread")
	}

	if thread.ImageURL != "" {
		delErr := tu.cloudinary.Delete("thread", util.GetFilenameWithoutExtension(thread.ImageURL))
		if delErr != nil {
			return errors.New("failed to delete image")
		}
	}

	err = tu.threadRepository.Delete(threadID)
	if err != nil {
		return errors.New("failed to delete thread")
	}

	return nil
}

func (tu *ThreadUseCase) AdminDelete(threadID primitive.ObjectID) (Domain, error) {
	thread, err := tu.threadRepository.GetByID(threadID)
	if err != nil {
		return Domain{}, errors.New("failed to get thread")
	}

	if thread.ImageURL != "" {
		err = tu.cloudinary.Delete("thread", util.GetFilenameWithoutExtension(thread.ImageURL))
		if err != nil {
			return Domain{}, errors.New("failed to delete image")
		}
	}

	err = tu.threadRepository.Delete(threadID)
	if err != nil {
		return Domain{}, errors.New("failed to delete thread")
	}

	return thread, nil
}
