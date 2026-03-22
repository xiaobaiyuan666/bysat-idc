package service

import (
	"fmt"

	"idc-finance/internal/modules/automation/domain"
)

type RetryExecutor func(task domain.Task) error

type RetryResult struct {
	SourceTask    domain.Task  `json:"sourceTask"`
	TriggeredTask *domain.Task `json:"triggeredTask,omitempty"`
	Message       string       `json:"message"`
}

func (service *Service) Retry(id int64, executor RetryExecutor) (RetryResult, error) {
	if service == nil || service.repository == nil {
		return RetryResult{}, fmt.Errorf("自动化任务服务未初始化")
	}
	if executor == nil {
		return RetryResult{}, fmt.Errorf("当前环境未配置任务重试执行器")
	}

	task, ok := service.repository.GetByID(id)
	if !ok {
		return RetryResult{}, fmt.Errorf("自动化任务不存在")
	}
	if task.Status == domain.TaskStatusRunning {
		return RetryResult{}, fmt.Errorf("任务仍在执行中，暂不支持重复提交")
	}

	beforeID := service.latestTaskID()
	if err := executor(task); err != nil {
		return RetryResult{}, err
	}

	result := RetryResult{
		SourceTask: task,
		Message:    "已重新提交自动化任务",
	}
	if latest, ok := service.latestTaskAfter(beforeID); ok {
		result.TriggeredTask = &latest
	}
	return result, nil
}

func (service *Service) latestTaskID() int64 {
	if service == nil || service.repository == nil {
		return 0
	}
	items := service.repository.List(domain.TaskFilter{Limit: 1})
	if len(items) == 0 {
		return 0
	}
	return items[0].ID
}

func (service *Service) latestTaskAfter(id int64) (domain.Task, bool) {
	if service == nil || service.repository == nil {
		return domain.Task{}, false
	}
	items := service.repository.List(domain.TaskFilter{Limit: 1})
	if len(items) == 0 || items[0].ID <= id {
		return domain.Task{}, false
	}
	return items[0], true
}
