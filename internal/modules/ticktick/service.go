package ticktick

import (
	"context"
	"errors"
	"fmt"
	"ticktick-ai/internal/domain"
)

type Client interface {
	CreateTask(ctx context.Context, input CreateTaskInput) (Task, error)
	FindTaskByTitle(ctx context.Context, title string) (Task, error)
	UpdateTask(ctx context.Context, task Task, updates domain.TaskUpdates) (Task, error)
	CompleteTask(ctx context.Context, task Task) error
}

type Service struct {
	client           Client
	defaultProjectID string
}

func NewService(client Client, defaultProjectID string) *Service {
	return &Service{
		client:           client,
		defaultProjectID: defaultProjectID,
	}
}

func (s *Service) ExecuteIntent(ctx context.Context, intent domain.ParsedIntent) (domain.TaskResult, error) {
	switch intent.Type {
	case domain.IntentCreateTask:
		return s.createTask(ctx, intent)
	case domain.IntentUpdateTask:
		return s.updateTask(ctx, intent)
	case domain.IntentCompleteTask:
		return s.completeTask(ctx, intent)
	case domain.IntentClarify:
		return domain.TaskResult{
			Action:  domain.IntentClarify,
			Message: intent.ClarificationQuestion,
		}, nil
	default:
		return domain.TaskResult{}, fmt.Errorf("unsupported intent type %q", intent.Type)
	}
}

func (s *Service) createTask(ctx context.Context, intent domain.ParsedIntent) (domain.TaskResult, error) {
	if intent.TaskTitle == "" {
		return domain.TaskResult{}, errors.New("task title is required")
	}

	task, err := s.client.CreateTask(ctx, CreateTaskInput{
		Title:     intent.TaskTitle,
		DueDate:   intent.DueDate,
		Priority:  intent.Priority,
		ProjectID: s.defaultProjectID,
		Tags:      intent.Tags,
	})
	if err != nil {
		return domain.TaskResult{}, err
	}

	return domain.TaskResult{
		ID:      task.ID,
		Title:   task.Title,
		Action:  domain.IntentCreateTask,
		Message: fmt.Sprintf("✅ Задача «%s» создана", task.Title),
	}, nil
}

func (s *Service) updateTask(ctx context.Context, intent domain.ParsedIntent) (domain.TaskResult, error) {
	if intent.TaskTitle == "" {
		return domain.TaskResult{}, errors.New("task title is required")
	}

	task, err := s.client.FindTaskByTitle(ctx, intent.TaskTitle)
	if err != nil {
		return domain.TaskResult{}, err
	}

	updated, err := s.client.UpdateTask(ctx, task, intent.Updates)
	if err != nil {
		return domain.TaskResult{}, err
	}

	return domain.TaskResult{
		ID:      updated.ID,
		Title:   updated.Title,
		Action:  domain.IntentUpdateTask,
		Message: fmt.Sprintf("✅ Задача «%s» обновлена", updated.Title),
	}, nil
}

func (s *Service) completeTask(ctx context.Context, intent domain.ParsedIntent) (domain.TaskResult, error) {
	if intent.TaskTitle == "" {
		return domain.TaskResult{}, errors.New("task title is required")
	}

	task, err := s.client.FindTaskByTitle(ctx, intent.TaskTitle)
	if err != nil {
		return domain.TaskResult{}, err
	}

	if err := s.client.CompleteTask(ctx, task); err != nil {
		return domain.TaskResult{}, err
	}

	return domain.TaskResult{
		ID:      task.ID,
		Title:   task.Title,
		Action:  domain.IntentCompleteTask,
		Message: fmt.Sprintf("✅ Задача «%s» завершена", task.Title),
	}, nil
}

type CreateTaskInput struct {
	Title     string
	DueDate   string
	Priority  domain.Priority
	ProjectID string
	Tags      []string
}

type Task struct {
	ID        string
	ProjectID string
	Title     string
	DueDate   string
	Priority  domain.Priority
	Tags      []string
}
