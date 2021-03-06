// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/baetyl/baetyl-cloud/v2/plugin (interfaces: MessageQueue)

// Package plugin is a generated GoMock package.
package plugin

import (
	mq "github.com/baetyl/baetyl-go/v2/mq"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockMessageQueue is a mock of MessageQueue interface
type MockMessageQueue struct {
	ctrl     *gomock.Controller
	recorder *MockMessageQueueMockRecorder
}

// MockMessageQueueMockRecorder is the mock recorder for MockMessageQueue
type MockMessageQueueMockRecorder struct {
	mock *MockMessageQueue
}

// NewMockMessageQueue creates a new mock instance
func NewMockMessageQueue(ctrl *gomock.Controller) *MockMessageQueue {
	mock := &MockMessageQueue{ctrl: ctrl}
	mock.recorder = &MockMessageQueueMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMessageQueue) EXPECT() *MockMessageQueueMockRecorder {
	return m.recorder
}

// Close mocks base method
func (m *MockMessageQueue) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockMessageQueueMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockMessageQueue)(nil).Close))
}

// Publish mocks base method
func (m *MockMessageQueue) Publish(arg0 string, arg1 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Publish", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Publish indicates an expected call of Publish
func (mr *MockMessageQueueMockRecorder) Publish(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockMessageQueue)(nil).Publish), arg0, arg1)
}

// Subscribe mocks base method
func (m *MockMessageQueue) Subscribe(arg0 string, arg1 mq.MQHandler) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Subscribe", arg0, arg1)
}

// Subscribe indicates an expected call of Subscribe
func (mr *MockMessageQueueMockRecorder) Subscribe(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subscribe", reflect.TypeOf((*MockMessageQueue)(nil).Subscribe), arg0, arg1)
}

// Unsubscribe mocks base method
func (m *MockMessageQueue) Unsubscribe(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Unsubscribe", arg0)
}

// Unsubscribe indicates an expected call of Unsubscribe
func (mr *MockMessageQueueMockRecorder) Unsubscribe(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unsubscribe", reflect.TypeOf((*MockMessageQueue)(nil).Unsubscribe), arg0)
}
