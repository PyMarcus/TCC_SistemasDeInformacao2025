package usecase

import (
	"net/http"

	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/constants"
	core "github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/core/ports/requests"
)

type APIRequestUsecase struct{
	Client core.APIRequestBuilder
}

func NewAPIRequestUsecase(clientService core.APIRequestBuilder) *APIRequestUsecase{
	return &APIRequestUsecase{
		Client: clientService,
	}
}

func (aru *APIRequestUsecase) Fetch(url string, headers map[string]string, body string) (*http.Response, error){
	response, err := aru.Client.
						SetMethod(constants.HTTP_GET_METHOD).
						SetURL(url). 
						SetHeaders(headers). 
						SetBody(body).
						Build()

	return response, err
}

func (aru *APIRequestUsecase) Post(url string, headers map[string]string, body string) (*http.Response, error){
	response, err := aru.Client.
						SetMethod(constants.HTTP_POST_METHOD).
						SetURL(url). 
						SetHeaders(headers). 
						SetBody(body).
						Build()

	return response, err
}