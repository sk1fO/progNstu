package business

import "monitoring/data"

type AlertService struct {
	dataAccess data.Database
}

func NewAlertService(db data.Database) *AlertService {
	return &AlertService{dataAccess: db}
}

func (s *AlertService) HandleAlert(buttonID string) {
	// Валидация и бизнес-логика
	if buttonID == "" {
		panic("Пустой ID кнопки!")
	}
	s.dataAccess.SaveAlert(buttonID)
}
