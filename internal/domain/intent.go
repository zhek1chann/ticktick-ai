package domain

type IntentType string

const (
	IntentCreateTask   IntentType = "create_task"
	IntentUpdateTask   IntentType = "update_task"
	IntentCompleteTask IntentType = "complete_task"
	IntentClarify      IntentType = "clarification_required"
	IntentBrief        IntentType = "brief"
	IntentListTasks    IntentType = "list_tasks"
)

type Priority string

const (
	PriorityNone   Priority = "none"
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

type ParsedIntent struct {
	Type                  IntentType  `json:"type"`
	TaskTitle             string      `json:"task_title,omitempty"`
	DueDate               string      `json:"due_date,omitempty"`
	Priority              Priority    `json:"priority,omitempty"`
	Project               string      `json:"project,omitempty"`
	Tags                  []string    `json:"tags,omitempty"`
	Updates               TaskUpdates `json:"updates,omitempty"`
	ClarificationQuestion string      `json:"clarification_question,omitempty"`
	BriefContent          string      `json:"brief_content,omitempty"`
	ListFilter            string      `json:"list_filter,omitempty"`
}

type TaskUpdates struct {
	NewTitle    string   `json:"new_title,omitempty"`
	NewDueDate  string   `json:"new_due_date,omitempty"`
	NewPriority Priority `json:"new_priority,omitempty"`
	AddTags     []string `json:"add_tags,omitempty"`
	RemoveTags  []string `json:"remove_tags,omitempty"`
}
