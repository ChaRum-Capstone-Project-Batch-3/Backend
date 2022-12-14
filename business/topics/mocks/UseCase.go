// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	multipart "mime/multipart"

	mock "github.com/stretchr/testify/mock"

	pagination "charum/dto/pagination"

	primitive "go.mongodb.org/mongo-driver/bson/primitive"

	topics "charum/business/topics"
)

// UseCase is an autogenerated mock type for the UseCase type
type UseCase struct {
	mock.Mock
}

// Create provides a mock function with given fields: domain, image
func (_m *UseCase) Create(domain *topics.Domain, image *multipart.FileHeader) (topics.Domain, error) {
	ret := _m.Called(domain, image)

	var r0 topics.Domain
	if rf, ok := ret.Get(0).(func(*topics.Domain, *multipart.FileHeader) topics.Domain); ok {
		r0 = rf(domain, image)
	} else {
		r0 = ret.Get(0).(topics.Domain)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*topics.Domain, *multipart.FileHeader) error); ok {
		r1 = rf(domain, image)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: id
func (_m *UseCase) Delete(id primitive.ObjectID) (topics.Domain, error) {
	ret := _m.Called(id)

	var r0 topics.Domain
	if rf, ok := ret.Get(0).(func(primitive.ObjectID) topics.Domain); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(topics.Domain)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(primitive.ObjectID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByID provides a mock function with given fields: id
func (_m *UseCase) GetByID(id primitive.ObjectID) (topics.Domain, error) {
	ret := _m.Called(id)

	var r0 topics.Domain
	if rf, ok := ret.Get(0).(func(primitive.ObjectID) topics.Domain); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(topics.Domain)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(primitive.ObjectID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByTopic provides a mock function with given fields: topic
func (_m *UseCase) GetByTopic(topic string) (topics.Domain, error) {
	ret := _m.Called(topic)

	var r0 topics.Domain
	if rf, ok := ret.Get(0).(func(string) topics.Domain); ok {
		r0 = rf(topic)
	} else {
		r0 = ret.Get(0).(topics.Domain)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(topic)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetManyWithPagination provides a mock function with given fields: _a0, domain
func (_m *UseCase) GetManyWithPagination(_a0 pagination.Request, domain *topics.Domain) ([]topics.Domain, int, int, error) {
	ret := _m.Called(_a0, domain)

	var r0 []topics.Domain
	if rf, ok := ret.Get(0).(func(pagination.Request, *topics.Domain) []topics.Domain); ok {
		r0 = rf(_a0, domain)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]topics.Domain)
		}
	}

	var r1 int
	if rf, ok := ret.Get(1).(func(pagination.Request, *topics.Domain) int); ok {
		r1 = rf(_a0, domain)
	} else {
		r1 = ret.Get(1).(int)
	}

	var r2 int
	if rf, ok := ret.Get(2).(func(pagination.Request, *topics.Domain) int); ok {
		r2 = rf(_a0, domain)
	} else {
		r2 = ret.Get(2).(int)
	}

	var r3 error
	if rf, ok := ret.Get(3).(func(pagination.Request, *topics.Domain) error); ok {
		r3 = rf(_a0, domain)
	} else {
		r3 = ret.Error(3)
	}

	return r0, r1, r2, r3
}

// Update provides a mock function with given fields: domain, image
func (_m *UseCase) Update(domain *topics.Domain, image *multipart.FileHeader) (topics.Domain, error) {
	ret := _m.Called(domain, image)

	var r0 topics.Domain
	if rf, ok := ret.Get(0).(func(*topics.Domain, *multipart.FileHeader) topics.Domain); ok {
		r0 = rf(domain, image)
	} else {
		r0 = ret.Get(0).(topics.Domain)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*topics.Domain, *multipart.FileHeader) error); ok {
		r1 = rf(domain, image)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewUseCase interface {
	mock.TestingT
	Cleanup(func())
}

// NewUseCase creates a new instance of UseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUseCase(t mockConstructorTestingTNewUseCase) *UseCase {
	mock := &UseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
