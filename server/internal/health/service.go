package health

import "net/http"

func HealthOk() *Health {
	health := NewHealth(http.StatusOK, "Health Ok")
	return health
}

func HealthNOk() *Health {
	health := NewHealth(http.StatusInternalServerError, "Health Not Ok")
	return health
}
