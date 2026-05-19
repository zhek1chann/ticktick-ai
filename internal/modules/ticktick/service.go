package ticktick

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"ticktick-ai/internal/domain"
)

type Client interface {
	CreateTask(ctx context.Context, input CreateTaskInput) (Task, error)
	FindTaskByTitle(ctx context.Context, title string) (Task, error)
	UpdateTask(ctx context.Context, task Task, updates domain.TaskUpdates) (Task, error)
	CompleteTask(ctx context.Context, task Task) error
	ListTasks(ctx context.Context) ([]Task, error)
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

func (s *Service) ExecuteIntent(ctx context.Context, intent domain.ParsedIntent, timezone string) (domain.TaskResult, error) {
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
	case domain.IntentBrief:
		return domain.TaskResult{
			Action:  domain.IntentBrief,
			Message: intent.BriefContent,
		}, nil
	case domain.IntentListTasks:
		return s.listTasks(ctx, intent.ListFilter, timezone)
	default:
		return domain.TaskResult{}, fmt.Errorf("unsupported intent type %q", intent.Type)
	}
}

func (s *Service) listTasks(ctx context.Context, filter string, timezone string) (domain.TaskResult, error) {
	all, err := s.client.ListTasks(ctx)
	if err != nil {
		return domain.TaskResult{}, err
	}

	tasks := all
	if filter == "today" {
		loc, err := time.LoadLocation(timezone)
		if err != nil {
			loc = time.UTC
		}
		today := time.Now().In(loc).Format("2006-01-02")
		var filtered []Task
		for _, t := range all {
			if strings.HasPrefix(t.DueDate, today) {
				filtered = append(filtered, t)
			}
		}
		tasks = filtered
	}

	if len(tasks) == 0 {
		msg := "Задач нет."
		if filter == "today" {
			msg = "На сегодня задач нет."
		}
		return domain.TaskResult{Action: domain.IntentListTasks, Message: msg}, nil
	}

	var sb strings.Builder
	if filter == "today" {
		sb.WriteString("📋 Задачи на сегодня:\n\n")
	} else {
		sb.WriteString("📋 Все задачи:\n\n")
	}
	for _, t := range tasks {
		sb.WriteString("• " + t.Title)
		if t.DueDate != "" {
			due, err := time.Parse(time.RFC3339, t.DueDate)
			if err == nil {
				sb.WriteString(" — " + due.Format("02.01 15:04"))
			}
		}
		if t.Priority != domain.PriorityNone {
			sb.WriteString(" [" + string(t.Priority) + "]")
		}
		sb.WriteString("\n")
	}

	return domain.TaskResult{Action: domain.IntentListTasks, Message: sb.String()}, nil
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
