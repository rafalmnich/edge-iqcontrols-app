// Code generated by mockery v2.12.1. DO NOT EDIT.

package mocks

import (
	fimpgo "github.com/futurehomeno/fimpgo"
	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// DeviceMapper is an autogenerated mock type for the DeviceMapper type
type DeviceMapper struct {
	mock.Mock
}

// Device provides a mock function with given fields: msg
func (_m *DeviceMapper) Device(msg *fimpgo.Message) (string, string, error) {
	ret := _m.Called(msg)

	var r0 string
	if rf, ok := ret.Get(0).(func(*fimpgo.Message) string); ok {
		r0 = rf(msg)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 string
	if rf, ok := ret.Get(1).(func(*fimpgo.Message) string); ok {
		r1 = rf(msg)
	} else {
		r1 = ret.Get(1).(string)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(*fimpgo.Message) error); ok {
		r2 = rf(msg)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// NewDeviceMapper creates a new instance of DeviceMapper. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewDeviceMapper(t testing.TB) *DeviceMapper {
	mock := &DeviceMapper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
