# TickTick AI Bot Architecture

## Goal

Telegram bot that accepts text or voice messages, extracts a task-management intent with Gemini, executes the action in TickTick, and returns a short confirmation back to Telegram.

## Stack

- Language: Go
- Telegram transport: `gopkg.in/telebot.v3`
- AI layer: `google.golang.org/api/generativeai/v1`
- TickTick integration: `github.com/slavkluev/go-ticktick`
- Optional state: PostgreSQL with `goqu`, if sessions, OAuth tokens, or request logs become necessary

## High-Level Flow

```text
Telegram user
    |
    | text or voice message
    v
Telegram Bot Layer
    |
    | normalized input: text or audio bytes
    v
Gemini Layer
    |
    | structured intent
    v
Application Service
    |
    | command mapped to TickTick API call
    v
TickTick Layer
    |
    | success or error
    v
Telegram Bot Layer
    |
    | short human-readable response
    v
Telegram user
```

## Layers

### Telegram Layer

Responsibilities:

- Start long polling or webhook handling.
- Receive Telegram updates.
- Detect supported message types:
  - text messages
  - voice messages
- For voice messages:
  - read `FileID`
  - request file metadata from Telegram
  - download `.ogg` file bytes into memory
- Convert Telegram input into an internal request object.
- Send final user-facing replies.

This layer should not know TickTick details and should not parse natural language itself.

### Gemini Layer

Responsibilities:

- Accept either plain text or audio bytes.
- Send the input to Gemini with a strict system instruction.
- Use function calling / tools to force structured output.
- Return a normalized intent object to the application service.

The model instruction should be strict:

```text
Do not answer with free-form text.
Extract the user's task-management intent and return only a structured function call.
Supported intents: create task, update task, complete task.
Extract task title, due date, priority, project, tags, and update fields when available.
If required data is missing, return a clarification-required result.
```

### Application Service

Responsibilities:

- Own the business workflow.
- Validate Gemini output.
- Decide whether the command can be executed immediately.
- Map the parsed intent to a TickTick operation.
- Build the final Telegram response message.
- Keep infrastructure-specific details behind interfaces.

Examples:

- `CreateTask` -> create TickTick task.
- `UpdateTask` -> search or resolve existing task, then update fields.
- `CompleteTask` -> search or resolve existing task, then mark completed.
- `ClarificationRequired` -> ask the user a short follow-up question.

### TickTick Layer

Responsibilities:

- Authenticate with TickTick.
- Call the official TickTick API V1 through `github.com/slavkluev/go-ticktick`.
- Hide client-specific request and response details from the application service.
- Return domain-level results and errors.

If the wrapper does not support required fields later, add a small adapter for the missing API surface. Unofficial V2 should be isolated behind the same TickTick interface.

### Optional State Layer

Use PostgreSQL only when there is a real need to persist data.

Possible future data:

- Telegram user to TickTick account mapping
- OAuth tokens
- pending clarification sessions
- request logs
- audit trail of executed task changes

The first version can stay stateless if a single TickTick account is configured through environment variables.

## Core Domain Model

```go
type IntentType string

const (
    IntentCreateTask IntentType = "create_task"
    IntentUpdateTask IntentType = "update_task"
    IntentCompleteTask IntentType = "complete_task"
    IntentClarify IntentType = "clarification_required"
)

type Priority string

const (
    PriorityNone Priority = "none"
    PriorityLow Priority = "low"
    PriorityMedium Priority = "medium"
    PriorityHigh Priority = "high"
)

type ParsedIntent struct {
    Type IntentType `json:"type"`
    TaskTitle string `json:"task_title,omitempty"`
    DueDate string `json:"due_date,omitempty"`
    Priority Priority `json:"priority,omitempty"`
    Project string `json:"project,omitempty"`
    Tags []string `json:"tags,omitempty"`
    Updates TaskUpdates `json:"updates,omitempty"`
    ClarificationQuestion string `json:"clarification_question,omitempty"`
}

type TaskUpdates struct {
    NewTitle string `json:"new_title,omitempty"`
    NewDueDate string `json:"new_due_date,omitempty"`
    NewPriority Priority `json:"new_priority,omitempty"`
    AddTags []string `json:"add_tags,omitempty"`
    RemoveTags []string `json:"remove_tags,omitempty"`
}
```

Dates should be normalized before calling TickTick. Relative expressions like `tomorrow`, `next Monday`, or `tonight` should be resolved by Gemini using the configured user timezone.

## Main Interfaces

```go
type TelegramBot interface {
    Run(ctx context.Context) error
    SendMessage(ctx context.Context, chatID int64, text string) error
}

type InputParser interface {
    ParseText(ctx context.Context, text string, timezone string) (ParsedIntent, error)
    ParseAudio(ctx context.Context, audio []byte, mimeType string, timezone string) (ParsedIntent, error)
}

type TaskManager interface {
    CreateTask(ctx context.Context, input CreateTaskInput) (TaskResult, error)
    UpdateTask(ctx context.Context, input UpdateTaskInput) (TaskResult, error)
    CompleteTask(ctx context.Context, input CompleteTaskInput) (TaskResult, error)
}
```

## Configuration

Required environment variables for the first version:

```text
TELEGRAM_BOT_TOKEN=
GEMINI_API_KEY=
TICKTICK_CLIENT_ID=
TICKTICK_CLIENT_SECRET=
TICKTICK_ACCESS_TOKEN=
USER_TIMEZONE=Asia/Almaty
```

Optional future variables:

```text
DATABASE_URL=
APP_ENV=local
LOG_LEVEL=info
```

## Error Handling

User-facing errors should be short and actionable.

Examples:

- Gemini cannot extract intent: `Не понял задачу. Напиши, что нужно сделать с задачей.`
- Missing task title: `Какую задачу нужно изменить?`
- TickTick task not found: `Не нашел задачу с таким названием в TickTick.`
- TickTick API error: `TickTick сейчас не принял запрос. Попробуй еще раз.`

Internal errors should be logged with enough context:

- Telegram chat ID
- message type
- Gemini request type
- parsed intent
- TickTick operation
- external API status or error code

## Package Shape

Target structure should follow the `olx-parser` style:

- `internal/app/provider.go` wires dependencies lazily.
- `internal/app/app.go` starts the Telegram bot and any background workers.
- Telegram code lives in `internal/modules/tg-bot`.
- Business logic lives in separate modules and is consumed by `tg-bot/service` through interfaces.
- Handlers stay thin and only translate Telegram events into service calls.

Suggested structure:

```text
cmd/
  main.go
internal/
  app/
    app.go
    provider.go
  config/
    config.go
    tg.go
    gemini.go
    ticktick.go
  modules/
    tg-bot/
      handler/
        handler.go
        start.go
        text_h.go
        voice_h.go
      middleware/
        guard.go
      model/
        models.go
      service/
        service.go
        message_svc.go
      repository/
        repository.go
        session.go
    ai/
      service.go
      gemini/
        client.go
        tools.go
        prompt.go
        mapper.go
    ticktick/
      service.go
      clients/
        ticktick/
          client.go
          mapper.go
      repository/
        repository.go
  domain/
    intent.go
    task.go
    session.go
```

For the MVP without database-backed sessions, `tg-bot/repository` and `domain/session.go` can be skipped. If clarification state is needed later, add them using the same approach as `olx-parser/internal/modules/tg-bot/repository`.

## Module Responsibilities

The module split should mirror `olx-parser`.

### `internal/modules/tg-bot`

```text
tg-bot/
  handler/      Telegram routes, callbacks, text and voice handlers
  middleware/   admin/user guard if needed
  model/        view models and Telegram-specific DTOs
  service/      bot workflow orchestration
  repository/   optional session/user persistence
```

The handler registers Telegram routes and message handlers:

```go
type Handler struct {
    svc *tgbot.BotService
}

func (h *Handler) RegisterRoutes(bot *tgbotapi.BotAPI) {
    // long-polling update dispatch will call text and voice handlers
}
```

The service owns the bot workflow and depends on interfaces:

```go
type IAIService interface {
    ParseText(ctx context.Context, text string, timezone string) (domain.ParsedIntent, error)
    ParseAudio(ctx context.Context, audio []byte, mimeType string, timezone string) (domain.ParsedIntent, error)
}

type ITaskService interface {
    ExecuteIntent(ctx context.Context, intent domain.ParsedIntent) (domain.TaskResult, error)
}

type BotService struct {
    ai        IAIService
    tasks     ITaskService
    timezone  string
}
```

### `internal/modules/ai`

This module wraps Gemini and returns domain intents. It should not know Telegram or TickTick details.

```text
ai/
  service.go
  gemini/
    client.go
    tools.go
    prompt.go
    mapper.go
```

### `internal/modules/ticktick`

This module owns task execution and TickTick API integration.

```text
ticktick/
  service.go
  clients/
    ticktick/
      client.go
      mapper.go
  repository/
    repository.go
```

`ticktick/service.go` receives `domain.ParsedIntent`, validates it, and calls the client adapter. If the wrapper is missing methods, only `clients/ticktick` should change.

## Provider Wiring

`internal/app/provider.go` should look close to `olx-parser`:

```go
type serviceProvider struct {
    ctx context.Context
    cfg *config.Config

    tgBot *tgbotapi.BotAPI

    aiService *ai.Service
    ticktickService *ticktick.Service
    tgBotService *tgbot.BotService
}

func (s *serviceProvider) TgBot(ctx context.Context) *tgbotapi.BotAPI {
    if s.tgBot == nil {
        bot, err := tgbotapi.NewBotAPI(s.Config().Tg().Token())
        if err != nil {
            log.Fatalf("failed to create telegram bot: %s", err.Error())
        }
        s.tgBot = bot
    }
    return s.tgBot
}

func (s *serviceProvider) TgBotService(ctx context.Context) *tgbot.BotService {
    if s.tgBotService == nil {
        s.tgBotService = tgbot.NewBotService(
            s.AIService(ctx),
            s.TickTickService(ctx),
            s.Config().App().Timezone(),
        )
    }
    return s.tgBotService
}
```

`internal/app/app.go` should start the bot similarly to `olx-parser`, using `telebot` long polling:

```go
func (a *App) Run(ctx context.Context) error {
    handler := tgbotHandler.New(a.serviceProvider.TgBotService(ctx))
    bot := a.serviceProvider.TgBot(ctx)

    handler.RegisterRoutes(bot)
    go bot.Start()
    closer.Wait()
    return nil
}
```

## MVP Scope

Version 1 should support:

- Telegram long polling.
- Text messages.
- Voice messages downloaded from Telegram by `FileID`.
- Gemini parsing into structured intent.
- TickTick task creation.
- Basic task update by title.
- Basic task completion by title.
- Short Telegram confirmation.
- Structured logs for failed external calls.

Out of scope for the first version:

- multi-user OAuth
- persistent sessions
- advanced task search disambiguation
- unofficial TickTick V2 adapter
- database-backed audit log

## Implementation Order

1. Add config for Telegram, Gemini, TickTick, and user timezone.
2. Implement Telegram update receiver and text handling.
3. Add voice download by Telegram `FileID`.
4. Define domain structs for parsed intents and task commands.
5. Implement Gemini parser with function calling.
6. Implement TickTick client adapter behind `TaskManager`.
7. Connect the application service.
8. Add happy-path and validation tests.
9. Add Docker and environment examples if needed.
