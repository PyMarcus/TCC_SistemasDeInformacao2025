package usecase

import (
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/constants"
	mock_core "github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/adapters/http/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestPost(t *testing.T){
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockClient := mock_core.NewMockAPIRequestBuilder(ctrl)

	expectedResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body: nil,
		Request: &http.Request{
			Method: "POST",
			URL: &url.URL{Scheme: "https", Host: "example.com", Path: "/"},
		},
	}
	mockClient.EXPECT().SetMethod(constants.HTTP_POST_METHOD).Return(mockClient)
	mockClient.EXPECT().SetURL("https://example.com").Return(mockClient)
	mockClient.EXPECT().SetHeaders(map[string]string{"Authorization": "Bearer test-token"}).Return(mockClient)
	mockClient.EXPECT().SetBody("").Return(mockClient)
	mockClient.EXPECT().Build().Return(expectedResponse, nil)

	mockUsecase := NewAPIRequestUsecase(mockClient)

	response, err := mockUsecase.Post("https://example.com", map[string]string{"Authorization": "Bearer test-token"}, "")
	assert.NoError(t, err)

	assert.NotNil(t, response)

	assert.Equal(t, "https://example.com/", response.Request.URL.String())
}

func TestPost_BuildError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_core.NewMockAPIRequestBuilder(ctrl)

	mockClient.EXPECT().SetMethod(constants.HTTP_POST_METHOD).Return(mockClient)
	mockClient.EXPECT().SetURL("https://invalid-url").Return(mockClient)
	mockClient.EXPECT().SetHeaders(map[string]string{}).Return(mockClient)
	mockClient.EXPECT().SetBody(`{"key":"value"}`).Return(mockClient)
	mockClient.EXPECT().Build().Return(nil, errors.New("invalid request"))

	mockUsecase := NewAPIRequestUsecase(mockClient)

	response, err := mockUsecase.Post("https://invalid-url", map[string]string{}, `{"key":"value"}`)

	assert.Nil(t, response)
	assert.Error(t, err)
	assert.EqualError(t, err, "invalid request")
}

func TestFetch(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockClient := mock_core.NewMockAPIRequestBuilder(ctrl)

	
	expectedResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       nil, // Para o teste, o corpo pode ser nil
		Request: &http.Request{
			Method: "GET",
			URL:    &url.URL{Scheme: "https", Host: "example.com", Path: "/"},
		},
	}
	mockClient.EXPECT().SetMethod(constants.HTTP_GET_METHOD).Return(mockClient)
	mockClient.EXPECT().SetURL("https://example.com").Return(mockClient)
	mockClient.EXPECT().SetHeaders(map[string]string{"Authorization": "Bearer test-token"}).Return(mockClient)
	mockClient.EXPECT().SetBody("").Return(mockClient)
	mockClient.EXPECT().Build().Return(expectedResponse, nil)

	mockUsecase := NewAPIRequestUsecase(mockClient)

	response, err := mockUsecase.Fetch("https://example.com", map[string]string{"Authorization": "Bearer test-token"}, "")

	assert.NoError(t, err)

	assert.NotNil(t, response)

	assert.Equal(t, "https://example.com/", response.Request.URL.String())

}

func TestFetch_BuildError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_core.NewMockAPIRequestBuilder(ctrl)

	mockClient.EXPECT().SetMethod(constants.HTTP_GET_METHOD).Return(mockClient)
	mockClient.EXPECT().SetURL("https://invalid-url").Return(mockClient)
	mockClient.EXPECT().SetHeaders(map[string]string{}).Return(mockClient)
	mockClient.EXPECT().SetBody("").Return(mockClient)
	mockClient.EXPECT().Build().Return(nil, errors.New("invalid request"))

	usecase := NewAPIRequestUsecase(mockClient)

	response, err := usecase.Fetch("https://invalid-url", map[string]string{}, "")

	assert.Nil(t, response)
	assert.Error(t, err)
	assert.EqualError(t, err, "invalid request")
}