package erro

import "fmt"

type ProviderError struct {
	error
	Status  string
	Message string
}

func NewProviderError(status string, message string) *ProviderError {
	return &ProviderError{Status: status, Message: message}
}

func (e *ProviderError) Error() string {
	return fmt.Sprintf("message %s: status %s", e.Message, e.Status)
}
