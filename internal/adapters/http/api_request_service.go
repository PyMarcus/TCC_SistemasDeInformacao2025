package adapters

import (
	"bytes"
	"fmt"
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

func (ars *APIRequestService) Build() (*http.Request, error) {
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

	client := &http.Client{
		Timeout: ars.request.Timeout,
	}

	response, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("error making API request: %v", err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code %d", response.StatusCode)
	}

	return req, nil
}
