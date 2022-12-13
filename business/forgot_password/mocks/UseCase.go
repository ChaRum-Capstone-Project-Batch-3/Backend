// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	forgot_password "charum/business/forgot_password"

	mock "github.com/stretchr/testify/mock"
)

// UseCase is an autogenerated mock type for the UseCase type
type UseCase struct {
	mock.Mock
}

// Generate provides a mock function with given fields: domain
func (_m *UseCase) Generate(domain *forgot_password.Domain) (forgot_password.Domain, error) {
	ret := _m.Called(domain)

	var r0 forgot_password.Domain
	if rf, ok := ret.Get(0).(func(*forgot_password.Domain) forgot_password.Domain); ok {
		r0 = rf(domain)
	} else {
		r0 = ret.Get(0).(forgot_password.Domain)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*forgot_password.Domain) error); ok {
		r1 = rf(domain)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByToken provides a mock function with given fields: token
func (_m *UseCase) GetByToken(token string) (forgot_password.Domain, error) {
	ret := _m.Called(token)

	var r0 forgot_password.Domain
	if rf, ok := ret.Get(0).(func(string) forgot_password.Domain); ok {
		r0 = rf(token)
	} else {
		r0 = ret.Get(0).(forgot_password.Domain)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdatePassword provides a mock function with given fields: domain
func (_m *UseCase) UpdatePassword(domain *forgot_password.Domain) (forgot_password.Domain, error) {
	ret := _m.Called(domain)

	var r0 forgot_password.Domain
	if rf, ok := ret.Get(0).(func(*forgot_password.Domain) forgot_password.Domain); ok {
		r0 = rf(domain)
	} else {
		r0 = ret.Get(0).(forgot_password.Domain)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*forgot_password.Domain) error); ok {
		r1 = rf(domain)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ValidateToken provides a mock function with given fields: token
func (_m *UseCase) ValidateToken(token string) (forgot_password.Domain, error) {
	ret := _m.Called(token)

	var r0 forgot_password.Domain
	if rf, ok := ret.Get(0).(func(string) forgot_password.Domain); ok {
		r0 = rf(token)
	} else {
		r0 = ret.Get(0).(forgot_password.Domain)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(token)
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
