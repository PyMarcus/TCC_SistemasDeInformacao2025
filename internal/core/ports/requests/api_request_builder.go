package core

import (
	"net/http"
	"time"
)


type APIRequestBuilder interface{
	SetMethod(method string) APIRequestBuilder
	SetURL(url string) APIRequestBuilder
	SetHeaders(headers map[string]string) APIRequestBuilder
	SetBody(body string) APIRequestBuilder
	SetTimeout(timeout time.Duration) APIRequestBuilder
	Build() (*http.Request, error)
}