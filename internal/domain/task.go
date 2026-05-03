package domain

type CreateTaskInput struct {
	Title     string
	DueDate   string
	Priority  Priority
	ProjectID string
	Tags      []string
}

type UpdateTaskInput struct {
	Title   string
	Updates TaskUpdates
}

type CompleteTaskInput struct {
	Title string
}

type TaskResult struct {
	ID      string
	Title   string
	Action  IntentType
	Message string
}
