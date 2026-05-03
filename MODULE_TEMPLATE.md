# Module Template

This file documents the standard module structure used in this project.
Copy this pattern when adding new modules (e.g. `user`, `shop`, `admin`).

---

## Directory Structure

```
internal/modules/{module_name}/
├── api/
│   ├── handler.go           # Handler struct + constructor
│   ├── {entity}_h.go        # HTTP handler methods
│   └── dto/
│       └── {entity}_dto.go  # Request/response DTOs + converters
├── model/
│   └── {entity}.go          # Domain models (no JSON/DB tags)
├── repository/
│   ├── repo.go              # Repo struct + constructor
│   ├── {entity}_r.go        # SQL queries
│   └── row/
│       └── {entity}.go      # DB row structs + To/From converters
├── service/
│   ├── service.go           # Service struct + constructor
│   └── {entity}_svc.go      # Business logic
└── routes.go                # Route registration
```

---

## Layer-by-Layer Examples

### `api/handler.go`

```go
package api

type Service{Module} interface {
    {entity}Service
}

type Handler{Module} struct {
    svc   Service{Module}
    debug bool
}

func NewHandler{Module}(svc Service{Module}, debug bool) *Handler{Module} {
    return &Handler{Module}{svc: svc, debug: debug}
}
```

### `api/{entity}_h.go`

```go
package api

type {entity}Service interface {
    Create{Entity}(ctx context.Context, ...) ({model}, error)
    Get{Entity}(ctx context.Context, id int) ({model}, error)
}

func (h *Handler{Module}) Create{Entity}(c *gin.Context) {
    var req dto.Create{Entity}Request
    if !validation.BindAndValidate(c, h.validator, &req) {
        return
    }

    result, err := h.svc.Create{Entity}(c.Request.Context(), ...)
    if err != nil {
        response.ErrorResponse(c, http.StatusInternalServerError, "internal error", err)
        return
    }

    response.SuccessResponse(c, http.StatusOK, "created", dto.From{Entity}(result))
}
```

### `api/dto/{entity}_dto.go`

```go
package dto

type Create{Entity}Request struct {
    Name string `json:"name" validate:"required"`
}

type {Entity}Response struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

func From{Entity}(m model.{Entity}) {Entity}Response {
    return {Entity}Response{ID: m.ID, Name: m.Name}
}

func From{Entities}(ms []model.{Entity}) []{Entity}Response {
    result := make([]{Entity}Response, len(ms))
    for i, m := range ms {
        result[i] = From{Entity}(m)
    }
    return result
}
```

### `model/{entity}.go`

```go
package model

// No JSON tags, no DB tags — pure business entity
type {Entity} struct {
    ID   int
    Name string
}
```

### `service/service.go`

```go
package service

type {Module}Repo interface {
    {entity}Repo
}

type Service struct {
    txManager db.TxManager
    repo      {Module}Repo
}

func NewService(txManager db.TxManager, repo {Module}Repo) *Service {
    return &Service{txManager: txManager, repo: repo}
}
```

### `service/{entity}_svc.go`

```go
package service

type {entity}Repo interface {
    Create{Entity}(ctx context.Context, e model.{Entity}) (int, error)
    {Entity}ByID(ctx context.Context, id int) (model.{Entity}, error)
}

func (s *Service) Create{Entity}(ctx context.Context, name string) (model.{Entity}, error) {
    var result model.{Entity}

    err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
        id, errTx := s.repo.Create{Entity}(ctx, model.{Entity}{Name: name})
        if errTx != nil {
            return errTx
        }
        var errTx2 error
        result, errTx2 = s.repo.{Entity}ByID(ctx, id)
        return errTx2
    })

    if err != nil {
        slog.ErrorContext(ctx, "Create{Entity}", "error", err)
        return model.{Entity}{}, err
    }

    return result, nil
}
```

### `repository/repo.go`

```go
package repository

type Repo struct {
    db db.Client
}

func New{Module}Repo(db db.Client) *Repo {
    return &Repo{db: db}
}
```

### `repository/{entity}_r.go`

```go
package repository

// {table_name} — table and columns
var {entity}Table = goqu.T("{table_name}")

const (
    {e}IDC   = "id"
    {e}NameC = "name"
)

var dialect = goqu.Dialect("postgres")

func (r *Repo) {Entity}ByID(ctx context.Context, id int) (model.{Entity}, error) {
    query, args, err := dialect.
        Select({e}IDC, {e}NameC).
        From({entity}Table).
        Where(goqu.C({e}IDC).Eq(id)).
        Prepared(true).
        ToSQL()
    if err != nil {
        return model.{Entity}{}, err
    }

    var rowData row.{Entity}
    err = r.db.DB().QueryRowContext(ctx, db.Query{Name: "{Entity}ByID", QueryRaw: query}, args...).
        Scan(&rowData.ID, &rowData.Name)
    if err != nil {
        return model.{Entity}{}, domain.MapError(err)
    }

    return row.To{Entity}(rowData), nil
}
```

### `repository/row/{entity}.go`

```go
package row

type {Entity} struct {
    ID       int
    Name     string
    ParentID *int  // nullable → pointer
}

func To{Entity}(r {Entity}) model.{Entity} {
    return model.{Entity}{ID: r.ID, Name: r.Name}
}

func To{Entities}(rows []{Entity}) []model.{Entity} {
    result := make([]model.{Entity}, len(rows))
    for i, r := range rows {
        result[i] = To{Entity}(r)
    }
    return result
}
```

### `routes.go`

```go
package {module}

func RegisterRoutes(router *gin.RouterGroup, h *api.Handler{Module}, jwt *middleware.JWTMiddleware) {
    g := router.Group("/{module}")
    {
        g.GET("/:id", h.Get{Entity})
        g.POST("/", jwt.RequireAuth(), h.Create{Entity})
    }
}
```

---

## Wiring a New Module in `provider.go`

```go
// Fields
{module}Handler *{module}Api.Handler{Module}
{module}Service *{module}Svc.Service
{module}Repo    *{module}Repo.Repo

// Methods
func (s *serviceProvider) {Module}Handler(ctx context.Context) *{module}Api.Handler{Module} {
    if s.{module}Handler == nil {
        s.{module}Handler = {module}Api.NewHandler{Module}(s.{Module}Service(ctx), s.Config().App().Debug())
    }
    return s.{module}Handler
}

func (s *serviceProvider) {Module}Service(ctx context.Context) *{module}Svc.Service {
    if s.{module}Service == nil {
        s.{module}Service = {module}Svc.NewService(s.TxManager(ctx), s.{Module}Repo(ctx))
    }
    return s.{module}Service
}

func (s *serviceProvider) {Module}Repo(ctx context.Context) *{module}Repo.Repo {
    if s.{module}Repo == nil {
        s.{module}Repo = {module}Repo.New{Module}Repo(s.DBClient(ctx))
    }
    return s.{module}Repo
}
```

Then in `app.go`:

```go
{module}Handler := a.serviceProvider.{Module}Handler(ctx)
{module}.RegisterRoutes(apiGroup, {module}Handler, jwtMiddleware)
```

---

## Roles & Auth Context

Roles are defined in `internal/domain/role.go`:

```go
RoleUser       Role = "user"
RoleShopOwner  Role = "shop_owner"
RoleSuperAdmin Role = "super_admin"
```

Access control in handlers:

```go
// Require any authenticated user
jwt.RequireAuth()

// Require specific role
jwt.RequireRole(domain.RoleSuperAdmin.String())
jwt.RequireRole(domain.RoleShopOwner.String())
```

Get current user in handler:

```go
userID := domain.GetUserIDFromContext(c)        // int64
role, _ := c.Get(domain.CtxKeyRole)             // string
```

---

## Migration Template

```sql
-- +goose Up

CREATE TABLE {table_name} (
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down

DROP TABLE {table_name};
```
