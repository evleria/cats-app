// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package repository

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	entities "github.com/evleria/cats-app/internal/repository/entities"

	uuid "github.com/google/uuid"
)

// MockCats is an autogenerated mock type for the Cats type
type MockCats struct {
	mock.Mock
}

// Delete provides a mock function with given fields: ctx, id
func (_m *MockCats) Delete(ctx context.Context, id uuid.UUID) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAll provides a mock function with given fields: ctx
func (_m *MockCats) GetAll(ctx context.Context) ([]entities.Cat, error) {
	ret := _m.Called(ctx)

	var r0 []entities.Cat
	if rf, ok := ret.Get(0).(func(context.Context) []entities.Cat); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entities.Cat)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOne provides a mock function with given fields: ctx, id
func (_m *MockCats) GetOne(ctx context.Context, id uuid.UUID) (entities.Cat, error) {
	ret := _m.Called(ctx, id)

	var r0 entities.Cat
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) entities.Cat); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(entities.Cat)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Insert provides a mock function with given fields: ctx, name, color, age, price
func (_m *MockCats) Insert(ctx context.Context, name string, color string, age int, price float64) (uuid.UUID, error) {
	ret := _m.Called(ctx, name, color, age, price)

	var r0 uuid.UUID
	if rf, ok := ret.Get(0).(func(context.Context, string, string, int, float64) uuid.UUID); ok {
		r0 = rf(ctx, name, color, age, price)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(uuid.UUID)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, int, float64) error); ok {
		r1 = rf(ctx, name, color, age, price)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdatePrice provides a mock function with given fields: ctx, id, price
func (_m *MockCats) UpdatePrice(ctx context.Context, id uuid.UUID, price float64) error {
	ret := _m.Called(ctx, id, price)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, float64) error); ok {
		r0 = rf(ctx, id, price)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
