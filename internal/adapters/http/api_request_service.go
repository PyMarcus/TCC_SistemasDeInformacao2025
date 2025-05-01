package adapters

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/core/domain"
	core "github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/core/ports/requests"
)

type APIRequestService struct {
	request *domain.APIRequestModel
}

func NewApiRequestService() core.APIRequestBuilder {
	return &APIRequestService{
		request: &domain.APIRequestModel{},
	}
}

func (ars *APIRequestService) SetMethod(method string) core.APIRequestBuilder {
	ars.request.Method = method
	return ars
}

func (ars *APIRequestService) SetURL(url string) core.APIRequestBuilder {
	ars.request.Url = url
	return ars
}

func (ars *APIRequestService) SetHeaders(headers map[string]string) core.APIRequestBuilder {
	ars.request.Headers = headers
	return ars
}

func (ars *APIRequestService) SetBody(body string) core.APIRequestBuilder {
	ars.request.Body = body
	return ars
}

func (ars *APIRequestService) SetTimeout(timeout time.Duration) core.APIRequestBuilder {
	ars.request.Timeout = timeout
	return ars
}

func (ars *APIRequestService) Build() (*http.Response, error) {
	req, err := http.NewRequest(
		ars.request.Method,
		ars.request.Url,
		bytes.NewBufferString(ars.request.Body),
	)

	if err != nil {
		return nil, err
	}

	for key, value := range ars.request.Headers {
		req.Header.Set(key, value)
	}
	/*
	client := &http.Client{
		Timeout: ars.request.Timeout,
	}*/

	client := &http.Client{}

	response, err := client.Do(req)

	if err != nil {
		return response, fmt.Errorf("error making API request: %v", err)
	}

	if response.StatusCode != http.StatusOK {
		resp, err  := io.ReadAll(response.Body)
		if err != nil {
			return response, fmt.Errorf("error to read body after status code not 200: %v", err)
		}
		return response, fmt.Errorf("API request failed with status code %d and response: %s", response.StatusCode, string(resp))
	}

	return response, nil
}
