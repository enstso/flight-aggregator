package health

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
