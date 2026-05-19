# TickTick AI

Telegram bot that accepts text or voice messages, parses the user's intent with Gemini, and creates, updates, or completes tasks in TickTick.

## Run Locally

```bash
cp example.env .env
make run
```

Required variables:

- `TG_TOKEN`
- `GEMINI_API_KEY`
- `TICKTICK_ACCESS_TOKEN`

Useful optional variables:

- `USER_TIMEZONE`, default `Asia/Almaty`
- `GEMINI_MODEL`, default `gemini-1.5-flash`
- `TICKTICK_DEFAULT_PROJECT_ID`

## Structure

```text
internal/app                 dependency wiring and startup
internal/config              env config
internal/domain              shared intent/task models
internal/modules/tg-bot      Telegram handlers and bot workflow
internal/modules/ai          Gemini parser service
internal/modules/ticktick    TickTick task execution service and client
```

The Telegram module follows the same shape as `olx-parser`: `handler`, `service`, `model`, and `middleware`.

## CI/CD

GitHub Actions runs on pull requests and pushes to `main`:

- `go test ./...`
- formatting check with `gofmt`
- binary build
- Docker image build
- Docker image push to `ghcr.io/<owner>/<repo>` on `main`

For deployment, pull the published image on your server and run it with `.env` values for `TG_TOKEN`, `GEMINI_API_KEY`, and `TICKTICK_ACCESS_TOKEN`.
