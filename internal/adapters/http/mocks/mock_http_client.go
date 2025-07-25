// Code generated by MockGen. DO NOT EDIT.
// Source: internal/core/ports/requests/api_request_builder.go

// Package mock_core is a generated GoMock package.
package mock_core

import (
	http "net/http"
	reflect "reflect"
	time "time"

	core "github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/core/ports/requests"
	gomock "github.com/golang/mock/gomock"
)

// MockAPIRequestBuilder is a mock of APIRequestBuilder interface.
type MockAPIRequestBuilder struct {
	ctrl     *gomock.Controller
	recorder *MockAPIRequestBuilderMockRecorder
}

// MockAPIRequestBuilderMockRecorder is the mock recorder for MockAPIRequestBuilder.
type MockAPIRequestBuilderMockRecorder struct {
	mock *MockAPIRequestBuilder
}

// NewMockAPIRequestBuilder creates a new mock instance.
func NewMockAPIRequestBuilder(ctrl *gomock.Controller) *MockAPIRequestBuilder {
	mock := &MockAPIRequestBuilder{ctrl: ctrl}
	mock.recorder = &MockAPIRequestBuilderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAPIRequestBuilder) EXPECT() *MockAPIRequestBuilderMockRecorder {
	return m.recorder
}

// Build mocks base method.
func (m *MockAPIRequestBuilder) Build() (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Build")
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Build indicates an expected call of Build.
func (mr *MockAPIRequestBuilderMockRecorder) Build() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Build", reflect.TypeOf((*MockAPIRequestBuilder)(nil).Build))
}

// SetBody mocks base method.
func (m *MockAPIRequestBuilder) SetBody(body string) core.APIRequestBuilder {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetBody", body)
	ret0, _ := ret[0].(core.APIRequestBuilder)
	return ret0
}

// SetBody indicates an expected call of SetBody.
func (mr *MockAPIRequestBuilderMockRecorder) SetBody(body interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetBody", reflect.TypeOf((*MockAPIRequestBuilder)(nil).SetBody), body)
}

// SetHeaders mocks base method.
func (m *MockAPIRequestBuilder) SetHeaders(headers map[string]string) core.APIRequestBuilder {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetHeaders", headers)
	ret0, _ := ret[0].(core.APIRequestBuilder)
	return ret0
}

// SetHeaders indicates an expected call of SetHeaders.
func (mr *MockAPIRequestBuilderMockRecorder) SetHeaders(headers interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetHeaders", reflect.TypeOf((*MockAPIRequestBuilder)(nil).SetHeaders), headers)
}

// SetMethod mocks base method.
func (m *MockAPIRequestBuilder) SetMethod(method string) core.APIRequestBuilder {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetMethod", method)
	ret0, _ := ret[0].(core.APIRequestBuilder)
	return ret0
}

// SetMethod indicates an expected call of SetMethod.
func (mr *MockAPIRequestBuilderMockRecorder) SetMethod(method interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetMethod", reflect.TypeOf((*MockAPIRequestBuilder)(nil).SetMethod), method)
}

// SetTimeout mocks base method.
func (m *MockAPIRequestBuilder) SetTimeout(timeout time.Duration) core.APIRequestBuilder {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetTimeout", timeout)
	ret0, _ := ret[0].(core.APIRequestBuilder)
	return ret0
}

// SetTimeout indicates an expected call of SetTimeout.
func (mr *MockAPIRequestBuilderMockRecorder) SetTimeout(timeout interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTimeout", reflect.TypeOf((*MockAPIRequestBuilder)(nil).SetTimeout), timeout)
}

// SetURL mocks base method.
func (m *MockAPIRequestBuilder) SetURL(url string) core.APIRequestBuilder {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetURL", url)
	ret0, _ := ret[0].(core.APIRequestBuilder)
	return ret0
}

// SetURL indicates an expected call of SetURL.
func (mr *MockAPIRequestBuilderMockRecorder) SetURL(url interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetURL", reflect.TypeOf((*MockAPIRequestBuilder)(nil).SetURL), url)
}
