package healthcheck

type Service interface {
	Status() HealthStatus
}

type service struct {
}

func New() Service {
	return &service{}
}

type HealthStatus struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (s *service) Status() HealthStatus {
	data := HealthStatus{
		Status:  "ok",
		Message: "Service is healthy",
	}
	return data
}
