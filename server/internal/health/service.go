package health

import "net/http"

// HealthOk creates and returns a pointer to a Health instance with status 200 and message "Health Ok".
func HealthOk() *Health {
	health := NewHealth(http.StatusOK, "Health Ok")
	return health
}

// HealthNOk creates and returns a pointer to a Health instance with status 503 and message "Health Not Ok".
func HealthNOk() *Health {
	health := NewHealth(http.StatusServiceUnavailable, "Health Not Ok")
	return health
}
