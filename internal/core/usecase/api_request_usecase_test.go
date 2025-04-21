package usecase

import (
	"errors"
	"net/http"
	"testing"

	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/constants"
	mock_core "github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/adapters/http/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestFetch(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockClient := mock_core.NewMockAPIRequestBuilder(ctrl)

	expectedRequest, _ := http.NewRequest("GET", "https://example.com", nil)

	mockClient.EXPECT().SetMethod(constants.HTTP_GET_METHOD).Return(mockClient)
	mockClient.EXPECT().SetURL("https://example.com").Return(mockClient)
	mockClient.EXPECT().SetHeaders(map[string]string{"Authorization": "Bearer test-token"}).Return(mockClient)
	mockClient.EXPECT().SetBody("").Return(mockClient)
	mockClient.EXPECT().Build().Return(expectedRequest, nil)

	mockUsecase := NewAPIRequestUsecase(mockClient)

	req, err := mockUsecase.Fetch("https://example.com", map[string]string{"Authorization": "Bearer test-token"}, "")

	assert.NoError(t, err)

	assert.NotNil(t, req)

	assert.Equal(t, expectedRequest, req)

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

	req, err := usecase.Fetch("https://invalid-url", map[string]string{}, "")

	assert.Nil(t, req)
	assert.Error(t, err)
	assert.EqualError(t, err, "invalid request")
}