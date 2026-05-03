package ticktick

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	tasks "ticktick-ai/internal/modules/ticktick"
)

var ErrTaskNotFound = errors.New("ticktick task not found")

type Client struct {
	httpClient  *http.Client
	baseURL     string
	accessToken string
}

func NewClient(httpClient *http.Client, baseURL string, accessToken string) *Client {
	return &Client{
		httpClient:  httpClient,
		baseURL:     strings.TrimRight(baseURL, "/"),
		accessToken: accessToken,
	}
}

func (c *Client) CreateTask(ctx context.Context, input tasks.CreateTaskInput) (tasks.Task, error) {
	reqBody := taskRequest{
		Title:     input.Title,
		DueDate:   emptyToNil(input.DueDate),
		Priority:  priorityToTickTick(input.Priority),
		ProjectID: emptyToNil(input.ProjectID),
		Tags:      input.Tags,
	}

	var resp taskResponse
	if err := c.doJSON(ctx, http.MethodPost, "/task", reqBody, &resp); err != nil {
		return tasks.Task{}, err
	}
	return resp.toDomain(), nil
}

func (c *Client) FindTaskByTitle(ctx context.Context, title string) (tasks.Task, error) {
	projects, err := c.projects(ctx)
	if err != nil {
		return tasks.Task{}, err
	}

	needle := strings.ToLower(strings.TrimSpace(title))
	for _, project := range projects {
		data, err := c.projectData(ctx, project.ID)
		if err != nil {
			return tasks.Task{}, err
		}

		for _, task := range data.Tasks {
			if strings.ToLower(strings.TrimSpace(task.Title)) == needle {
				task.ProjectID = project.ID
				return task.toDomain(), nil
			}
		}
	}

	return tasks.Task{}, ErrTaskNotFound
}

func (c *Client) UpdateTask(ctx context.Context, task tasks.Task, updates tasksDomainUpdates) (tasks.Task, error) {
	title := task.Title
	if updates.NewTitle != "" {
		title = updates.NewTitle
	}

	dueDate := task.DueDate
	if updates.NewDueDate != "" {
		dueDate = updates.NewDueDate
	}

	priority := task.Priority
	if updates.NewPriority != "" {
		priority = updates.NewPriority
	}

	tags := mergeTags(task.Tags, updates.AddTags, updates.RemoveTags)

	reqBody := taskRequest{
		ID:        task.ID,
		ProjectID: emptyToNil(task.ProjectID),
		Title:     title,
		DueDate:   emptyToNil(dueDate),
		Priority:  priorityToTickTick(priority),
		Tags:      tags,
	}

	var resp taskResponse
	if err := c.doJSON(ctx, http.MethodPost, "/task/"+url.PathEscape(task.ID), reqBody, &resp); err != nil {
		return tasks.Task{}, err
	}
	return resp.toDomain(), nil
}

func (c *Client) CompleteTask(ctx context.Context, task tasks.Task) error {
	if task.ProjectID == "" {
		return errors.New("project id is required to complete task")
	}
	path := "/project/" + url.PathEscape(task.ProjectID) + "/task/" + url.PathEscape(task.ID) + "/complete"
	return c.doJSON(ctx, http.MethodPost, path, nil, nil)
}

func (c *Client) projects(ctx context.Context) ([]projectResponse, error) {
	var resp []projectResponse
	if err := c.doJSON(ctx, http.MethodGet, "/project", nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) projectData(ctx context.Context, projectID string) (projectDataResponse, error) {
	var resp projectDataResponse
	path := "/project/" + url.PathEscape(projectID) + "/data"
	if err := c.doJSON(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return projectDataResponse{}, err
	}
	return resp, nil
}

func (c *Client) doJSON(ctx context.Context, method string, path string, input any, output any) error {
	var body io.Reader
	if input != nil {
		data, err := json.Marshal(input)
		if err != nil {
			return err
		}
		body = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("ticktick api status %d: %s", resp.StatusCode, string(respBody))
	}

	if output == nil || len(respBody) == 0 {
		return nil
	}

	return json.Unmarshal(respBody, output)
}
