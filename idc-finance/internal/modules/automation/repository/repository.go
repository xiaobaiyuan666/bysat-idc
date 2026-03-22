package repository

import "idc-finance/internal/modules/automation/domain"

type Repository interface {
	List(filter domain.TaskFilter) []domain.Task
	GetByID(id int64) (domain.Task, bool)
	Create(task domain.Task) domain.Task
	Update(task domain.Task) (domain.Task, bool)
}
