package service

import automationService "idc-finance/internal/modules/automation/service"

func (service *Service) startTask(request automationService.StartTaskRequest) int64 {
	if service == nil || service.tasks == nil {
		return 0
	}
	task := service.tasks.Start(request)
	return task.ID
}

func (service *Service) successTask(taskID int64, message string, result any) {
	if service == nil || service.tasks == nil || taskID == 0 {
		return
	}
	service.tasks.MarkSuccess(taskID, message, result)
}

func (service *Service) failTask(taskID int64, message string, result any) {
	if service == nil || service.tasks == nil || taskID == 0 {
		return
	}
	service.tasks.MarkFailed(taskID, message, result)
}

func (service *Service) blockTask(taskID int64, message string, result any) {
	if service == nil || service.tasks == nil || taskID == 0 {
		return
	}
	service.tasks.MarkBlocked(taskID, message, result)
}
