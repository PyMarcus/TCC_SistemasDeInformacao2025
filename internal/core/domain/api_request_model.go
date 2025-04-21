package domain

import "time"

type APIRequestModel struct{
	Method  string
	Url     string
	Headers map[string]string
	Body    string
	Timeout time.Duration
}