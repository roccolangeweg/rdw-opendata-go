package rdw_opendata_go

import "fmt"

type ApiError struct {
	StatusCode int
	Message    string
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("status code: %d, error: %v", e.StatusCode, e.Message)
}

type AppTokenMissingError struct{}

func (e *AppTokenMissingError) Error() string {
	return "app token is missing"
}

type AppTokenInvalidError struct{}

func (e *AppTokenInvalidError) Error() string {
	return "app token is invalid"
}

type ApiLimitExceededError struct{}

func (e *ApiLimitExceededError) Error() string {
	return "API limit exceeded"
}
