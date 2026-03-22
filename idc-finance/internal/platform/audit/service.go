package audit

type Service struct {
	repository Repository
}

func New(repo Repository) *Service {
	return &Service{repository: repo}
}

func (service *Service) Record(entry Entry) Entry {
	return service.repository.Record(entry)
}

func (service *Service) List() []Entry {
	return service.repository.List()
}

func (service *Service) ListByTarget(targetType string, targetID int64) []Entry {
	return service.repository.ListByTarget(targetType, targetID)
}
