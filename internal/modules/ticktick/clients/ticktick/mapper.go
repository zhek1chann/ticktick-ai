package ticktick

import (
	"ticktick-ai/internal/domain"
	tasks "ticktick-ai/internal/modules/ticktick"
)

type tasksDomainUpdates = domain.TaskUpdates

type taskRequest struct {
	ID        string   `json:"id,omitempty"`
	ProjectID *string  `json:"projectId,omitempty"`
	Title     string   `json:"title"`
	DueDate   *string  `json:"dueDate,omitempty"`
	Priority  int      `json:"priority"`
	Tags      []string `json:"tags,omitempty"`
}

type taskResponse struct {
	ID        string   `json:"id"`
	ProjectID string   `json:"projectId"`
	Title     string   `json:"title"`
	DueDate   string   `json:"dueDate"`
	Priority  int      `json:"priority"`
	Tags      []string `json:"tags"`
}

type projectResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type projectDataResponse struct {
	Tasks []taskResponse `json:"tasks"`
}

func (r taskResponse) toDomain() tasks.Task {
	return tasks.Task{
		ID:        r.ID,
		ProjectID: r.ProjectID,
		Title:     r.Title,
		DueDate:   r.DueDate,
		Priority:  priorityFromTickTick(r.Priority),
		Tags:      r.Tags,
	}
}

func priorityToTickTick(priority domain.Priority) int {
	switch priority {
	case domain.PriorityLow:
		return 1
	case domain.PriorityMedium:
		return 3
	case domain.PriorityHigh:
		return 5
	default:
		return 0
	}
}

func priorityFromTickTick(priority int) domain.Priority {
	switch priority {
	case 1:
		return domain.PriorityLow
	case 3:
		return domain.PriorityMedium
	case 5:
		return domain.PriorityHigh
	default:
		return domain.PriorityNone
	}
}

func emptyToNil(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func mergeTags(current []string, add []string, remove []string) []string {
	removed := make(map[string]struct{}, len(remove))
	for _, tag := range remove {
		removed[tag] = struct{}{}
	}

	seen := make(map[string]struct{}, len(current)+len(add))
	var result []string

	for _, tag := range current {
		if _, ok := removed[tag]; ok {
			continue
		}
		if _, ok := seen[tag]; ok {
			continue
		}
		seen[tag] = struct{}{}
		result = append(result, tag)
	}

	for _, tag := range add {
		if _, ok := removed[tag]; ok {
			continue
		}
		if _, ok := seen[tag]; ok {
			continue
		}
		seen[tag] = struct{}{}
		result = append(result, tag)
	}

	return result
}
