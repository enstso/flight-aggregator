package health

import "net/http"

type Health struct {
	Status  int
	Message string
}

func NewHealth(status int, message string) *Health {
	return &Health{
		Status:  status,
		Message: message,
	}
}

// Ok creates and returns a pointer to a Health instance with status 200 and message "Health Ok".
func Ok() *Health {
	health := NewHealth(http.StatusOK, "Health Ok")
	return health
}

// NOk creates and returns a pointer to a Health instance with status 503 and message "Health Not Ok".
func NOk() *Health {
	health := NewHealth(http.StatusServiceUnavailable, "Health Not Ok")
	return health
}
