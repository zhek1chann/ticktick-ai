no# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Baru is a Go-based e-commerce backend service for **small grocery store delivery orders**. It manages products, categories, shops, and orders for local convenience stores and mini-markets.

### Business Model
- **Target Users**: Small grocery stores (corner shops, mini-markets) that want to accept delivery orders
- **MVP Scope**: Supporting 1-2 stores initially with direct navigation (no shop selection UI)
- **Future Plans**: Multi-store support with maps/location-based shop selection
- **Current Flow**: Users navigate directly to a specific shop's page (e.g., via URL or bookmark)

### MVP Order Flow (Simplified - No Auth, No Backend Cart)
1. **Customer browses catalog** - Frontend fetches products from backend API
2. **Customer adds to cart** - Cart managed entirely in frontend (localStorage/React state)
3. **Customer fills delivery form** - Name, phone number, delivery address
4. **Frontend sends complete order** - Single API call with all order details
5. **Backend saves order** - Store in database with status "pending"
6. **Telegram notification to shop** - Backend sends order details to shop owner via Telegram bot
7. **Shop prepares order** - Shop owner reviews order and starts preparation
8. **Shop marks order as "ready"** - Shop owner updates order status (via Telegram bot or admin)
9. **Backend notifies customer** - "Your order is ready! Payment request sent to your Kaspi account"
10. **Manual Kaspi payment** - Shop owner creates payment request in Kaspi using customer's phone number
11. **Customer pays via Kaspi** - Customer receives and pays the request in Kaspi mobile app
12. **Customer confirms payment** - Customer clicks "I paid" in the order page
13. **Telegram notification to shop** - Backend sends "Payment confirmed" message to shop
14. **Delivery/pickup** - Shop delivers order or customer picks up

**What's NOT in MVP:**
- ❌ User authentication/login
- ❌ Backend cart system (frontend handles cart state)
- ❌ User accounts/profiles
- ❌ Automated payment gateway integration
- ❌ Admin dashboard (shop owner uses Telegram bot)
- ❌ Real-time order tracking
- ❌ Shop selection UI (direct navigation to shop page)

**What IS in MVP:**
- ✅ Product catalog API (browse products by shop)
- ✅ Category filtering
- ✅ Order submission API (guest checkout with phone number)
- ✅ Telegram bot integration (order notifications to shop owner)
- ✅ Order status management (pending → ready → paid → completed)
- ✅ Manual Kaspi payment workflow
- ✅ Customer payment confirmation
- ✅ Order storage in database
- ✅ Inventory management (Excel upload by shop owner)

### Payment Flow (Kaspi Manual Integration)
The payment system is **deliberately manual** to avoid payment gateway fees and complexity in MVP:

1. **Customer submits order** → Order saved with status "pending"
2. **Shop receives Telegram notification** → Reviews order
3. **Shop prepares order** → Marks as "ready" (via Telegram bot or admin endpoint)
4. **Customer gets notification** → "Order ready, payment request sent to your Kaspi"
5. **Shop manually creates Kaspi payment** → Uses customer's phone number in Kaspi app
6. **Customer receives Kaspi request** → Pays in Kaspi mobile app
7. **Customer confirms payment** → Clicks "I paid" button in order page
8. **Shop receives confirmation** → Telegram notification: "Customer confirmed payment"
9. **Shop delivers order** → Manual fulfillment

**Why manual Kaspi payments?**
- ✅ No payment gateway integration fees
- ✅ Direct transfer between customer and shop (no intermediary)
- ✅ Shop has full control over payment requests
- ✅ Familiar process for both shop owners and customers in Kazakhstan
- ✅ Simpler implementation for MVP

### Key Features
- **Product Catalog**: Shop-specific inventory with prices and stock levels
- **Inventory Management**: Shop owners upload Excel files to update prices/stock via barcode matching
- **Product Synchronization**: Integration with external systems (like Umag) for product data import
- **Multi-format Barcode Support**: Handles JSON arrays, comma-separated, and semicolon-separated barcodes in Excel files
- **Guest Checkout**: Accept orders without user authentication (phone number as identifier)
- **Telegram Bot Integration**: Bi-directional communication with shop owners
  - Order notifications (new order placed)
  - Status updates (order ready, payment confirmed)
  - Shop owner can mark orders as ready via Telegram bot
- **Order Management**: Track order lifecycle (pending → ready → paid → completed)
- **Manual Kaspi Payments**: Shop owner creates payment requests, customer confirms

### Technical Stack
- **Backend**: Go with Gin framework
- **Database**: PostgreSQL for persistent storage, SQLite for catalog data migration
- **Architecture**: Clean architecture with handler/service/repository layers
- **Authentication**: Not implemented in MVP (guest checkout only)
- **Cart**: Frontend-only (no backend cart state management)
- **Notifications**: Telegram Bot API for shop owner communication
- **Payment**: Manual Kaspi payment creation by shop owner (no gateway integration)
- **File Storage**: S3-compatible storage (MinIO) for product images
- **Inventory Upload**: Excel file parsing with barcode-based product matching

## Development Commands

### Running the Application
```bash
go run cmd/main.go
```

### Building
```bash
go build -o bin/app cmd/main.go
```

### Running Tests
```bash
# Run all tests
go test ./...

# Run tests for a specific package
go test ./pkg/db/prettier

# Run tests with verbose output
go test -v ./...
```

### Database Operations
```bash
# Start PostgreSQL via Docker Compose
docker-compose up pg -d

# Run migrations (via Docker)
docker-compose up migrator

# Stop services
docker-compose down
```

### Docker Build
```bash
docker build -t baru -f dockerfile .
```

## Architecture

### Module Structure
The codebase follows a clean architecture pattern organized by feature modules under `internal/modules/`:

- **product**: Product catalog management (categories, products, brands, barcodes)
- **shop**: Shop management (shops, settings, delivery zones via shop_houses)
- **catalog**: Catalog synchronization with SQLite database
- **orders**: Order processing and management

Each module follows this internal structure:
- `api/`: HTTP handlers (Gin framework)
- `model/`: Domain models and errors
- `repository/`: Database access layer with row mappings
- `service/`: Business logic layer
- `routes.go`: Route registration

### Dependency Injection Pattern
The application uses a service provider pattern (`internal/app/provider.go`) with lazy initialization. Dependencies are constructed on first access and cached for subsequent requests. This includes:
- Database clients (PostgreSQL via pgx/v4)
- Transaction managers
- S3 clients (AWS SDK v2)
- Module-specific repositories, services, and handlers

### Database Layer (`pkg/db/`)
- **Abstraction**: `db.Client` interface wraps the database connection
- **Implementation**: PostgreSQL via `pgx/v4` in `pkg/db/pg/`
- **Query Building**: Uses goqu for type-safe SQL query construction
- **Query Logging**: `pkg/db/prettier/` prettifies SQL queries for logging
- **Transactions**: `pkg/db/transaction/` manages transactions with support for nested transactions via context

### Transaction Pattern
Transactions are context-based. The transaction manager checks for an existing transaction in context before starting a new one:
```go
// Start a transaction with ReadCommitted isolation
txManager.ReadCommitted(ctx, func(ctx context.Context) error {
    // All repository calls within this function share the same transaction
    return service.doWork(ctx)
})
```

### Configuration (`internal/config/`)
Configuration is loaded from environment variables (via `.env` file using godotenv):
- `PG_*`: PostgreSQL connection settings
- `APP_*`: Application settings (mode, debug, log level)
- `HTTP_*`: HTTP server configuration
- `S3_*`: S3/MinIO configuration for file storage

### HTTP Layer
- **Framework**: Gin
- **Documentation**: Swagger/OpenAPI via swaggo (accessible at `/swagger/*`)
- **CORS**: Configured to allow all origins (for development)
- **Logging**: Custom Gin logger middleware in `pkg/logger/gin.go`
- **Routes**: Mounted under `/api` prefix

### External Integrations
- **Umag Client** (`internal/modules/shop/umag/`): HTTP client for syncing product data from Umag API
- **S3 Storage** (`pkg/s3/`): File upload/download via AWS SDK v2 (compatible with MinIO)

### Migrations
Database migrations use goose format (SQL files with `-- +goose Up` and `-- +goose Down` directives) located in `migrations/`:
1. `1_base_db.sql`: Core product catalog schema (categories, brands, products, barcodes)
2. `2_add_shops_table.sql`: Shop management tables
3. `3_add_product_shop.sql`: Product-shop junction table for pricing and inventory

The migrator service (`data/migrator.go`) handles data migration from SQLite to PostgreSQL.

## Key Design Patterns

The **product module** is the reference implementation. All new modules MUST follow these patterns exactly.

### Module Organization (ETALON: `internal/modules/product/`)

Each module follows this structure:
```
internal/modules/{module_name}/
├── api/
│   ├── handler.go           # Main handler struct with service composition
│   ├── {entity}_h.go        # Entity-specific handler methods (e.g., category_h.go)
│   └── dto/
│       └── {entity}_dto.go  # Request/response DTOs and domain-to-DTO converters
├── model/
│   └── {entity}.go          # Domain models (pure business entities)
├── repository/
│   ├── repo.go              # Main Repo struct wrapping db.Client
│   ├── {entity}_r.go        # Entity-specific repository methods (e.g., category_r.go)
│   └── row/
│       └── {entity}.go      # Database row structs and row↔model converters
├── service/
│   ├── service.go           # Main Service struct with composed interfaces
│   └── {entity}_svc.go      # Entity-specific service methods (e.g., category_svc.go)
└── routes.go                # Route registration function
```

### Naming Conventions

**Files:**
- Handlers: `{entity}_h.go` (e.g., `category_h.go`)
- Services: `{entity}_svc.go` (e.g., `category_svc.go`)
- Repositories: `{entity}_r.go` (e.g., `category_r.go`)
- Main structs: `handler.go`, `service.go`, `repo.go`

**Structs and Functions:**
- Handler: `Handler{module}` (e.g., `Handlerproduct`)
- Service: `Service` (generic, lives in service package)
- Repository: `Repo` (generic, lives in repository package)
- Constructor: `NewHandler{module}`, `NewService`, `New{Module}Repo`

**Tables and Columns (Repository):**
```go
// {table_name} — table and columns
var categoryNodesTable = goqu.T("category_nodes")

const (
    cnIDC       = "id"
    cnNameC     = "name"
    cnParentIDC = "parent_id"
)
```
- Tables use `goqu.T()` as a `var`
- Columns are string `const` with a `C` suffix (e.g., `cnNameC`)
- Column prefix = short table alias (e.g., `cn` for category_nodes, `ps` for product_shop)
- Group columns directly under their table

### Layer-by-Layer Patterns

#### 1. API Layer (Handler Pattern)

**File: `api/handler.go`**
```go
// Compose service interfaces
type Serviceproduct interface {
    categoryService
    productService
}

type Handlerproduct struct {
    svcProduct Serviceproduct
    debug      bool
}

func NewHandlerproduct(svcProduct Serviceproduct, debug bool) *Handlerproduct {
    return &Handlerproduct{
        svcProduct: svcProduct,
        debug:      debug,
    }
}
```

**File: `api/category_h.go`**
```go
// Define minimal interface for this handler file
type categoryService interface {
    CategoriesTree(ctx context.Context, depthLvl int) ([]model.Category, error)
    Category(ctx context.Context, categoryID int) (model.Category, error)
}

// Handler methods follow this pattern:
func (h *Handlerproduct) GetCategoriesTree(c *gin.Context) {
    // 1. Parse and validate input
    depthLv, err := strconv.Atoi(c.Query("depthLvl"))
    if err != nil {
        c.JSON(http.StatusBadRequest, response.Response[any]{
            Status: "invalid depthLvl parameter",
            Data:   nil,
        })
        return
    }

    // 2. Call service layer
    categories, err := h.svcProduct.CategoriesTree(c.Request.Context(), depthLv)
    if err != nil {
        c.JSON(http.StatusInternalServerError, response.Response[any]{
            Status: "internal error",
            Data:   nil,
        })
        return
    }

    // 3. Convert to DTO and return
    c.JSON(http.StatusOK, response.Response[dto.GetCategoriesTreeResponse]{
        Status: "ok",
        Data: dto.GetCategoriesTreeResponse{
            Categories: dto.FromCategories(categories),
        },
    })
}
```

**Handler Responsibilities:**
- Parse request parameters (query, path, body)
- Validate input format (NOT business validation)
- Call service layer with `c.Request.Context()`
- Convert domain models to DTOs
- Return appropriate HTTP status codes
- Handle errors generically (no error details leaked to client)

#### 2. DTO Layer Pattern

**File: `api/dto/product_dto.go`**
```go
// Swagger wrapper types (for documentation only)
type GetCategoriesTreeSwaggerResponse struct {
    Status string                    `json:"status"`
    Data   GetCategoriesTreeResponse `json:"data"`
}

// Actual request/response types
type GetCategoriesTreeResponse struct {
    Categories []Category `json:"categories"`
}

type Category struct {
    ID       int        `json:"id"`
    Name     string     `json:"name"`
    Children []Category `json:"children"`
}

// Converter functions (domain → DTO)
func FromCategories(categories []model.Category) []Category {
    if len(categories) == 0 {
        return nil
    }
    result := make([]Category, len(categories))
    for i, category := range categories {
        result[i] = FromCategory(category)
    }
    return result
}

func FromCategory(category model.Category) Category {
    return Category{
        ID:       category.ID,
        Name:     category.Name,
        Children: FromCategories(category.Children),
    }
}
```

**DTO Responsibilities:**
- Define JSON structure for API contracts
- Provide `From{Entity}` converters (domain → DTO)
- Provide `To{Entity}` converters (DTO → domain) for requests
- Keep DTOs flat and JSON-serializable
- No business logic in DTOs

#### 3. Service Layer Pattern

**File: `service/service.go`**
```go
// Compose repository interfaces
type ProductRepo interface {
    productRepo
    brandRepo
    categoryRepo
}

type Service struct {
    txManager db.TxManager
    repo      ProductRepo
}

func NewService(txManager db.TxManager, repo ProductRepo) *Service {
    return &Service{txManager: txManager, repo: repo}
}
```

**File: `service/category_svc.go`**
```go
// Define minimal interface for this service file
type categoryRepo interface {
    Category(ctx context.Context, categoryID int) (model.Category, error)
    CategoryChildren(ctx context.Context, categoryID int) ([]model.Category, error)
    MainCategories(ctx context.Context) ([]model.Category, error)
}

func (s *Service) CategoriesTree(ctx context.Context, depthLvl int) ([]model.Category, error) {
    var categories []model.Category
    err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
        var errTx error
        categories, errTx = s.repo.MainCategories(ctx)
        if errTx != nil {
            return errTx
        }
        for i := range categories {
            if depthLvl-1 > 0 {
                children, errTx := s.CategoryChildren(ctx, depthLvl-1, categories[i].ID)
                if errTx != nil {
                    return errTx
                }
                categories[i].Children = children
            }
        }
        return nil
    })

    if err != nil {
        slog.ErrorContext(ctx, "CategoriesTree", "error", err)
        return nil, err
    }
    return categories, nil
}
```

**Service Responsibilities:**
- Contain ALL business logic
- Wrap repository calls in transactions via `txManager`
- Orchestrate multiple repository calls
- Return domain models (never DTOs or rows)
- Log errors with `slog.ErrorContext`
- Use named return values for clarity in transaction closures

**Transaction Pattern:**
- ALWAYS use `txManager.ReadCommitted(ctx, func(ctx context.Context) error {...})`
- Use closure with named return variables captured from outer scope
- Return `errTx` (not `err`) inside transaction closure
- Log errors AFTER transaction completes, not inside

#### 4. Repository Layer Pattern

**File: `repository/repo.go`**
```go
type Repo struct {
    db db.Client
}

func NewProductRepo(db db.Client) *Repo {
    return &Repo{db: db}
}
```

**File: `repository/category_r.go`**
```go
// category_nodes — table and columns
var categoryNodesTable = goqu.T("category_nodes")

const (
    cnIDC       = "id"
    cnNameC     = "name"
    cnParentIDC = "parent_id"
)

var dialect = goqu.Dialect("postgres")

func (r *Repo) Category(ctx context.Context, id int) (model.Category, error) {
    query, args, err := dialect.
        Select(cnIDC, cnNameC, cnParentIDC).
        From(categoryNodesTable).
        Where(goqu.C(cnIDC).Eq(id)).
        Prepared(true).
        ToSQL()
    if err != nil {
        return model.Category{}, err
    }

    q := db.Query{
        Name:     "Category",
        QueryRaw: query,
    }

    var categoryRow row.CategoryNode
    queryRow := r.db.DB().QueryRowContext(ctx, q, args...)
    err = queryRow.Scan(&categoryRow.ID, &categoryRow.Name, &categoryRow.ParentID)
    if err != nil {
        return model.Category{}, err
    }

    return row.ToCategory(categoryRow), nil
}

func (r *Repo) CategoryChildren(ctx context.Context, id int) ([]model.Category, error) {
    query, args, err := dialect.
        Select(cnIDC, cnNameC, cnParentIDC).
        From(categoryNodesTable).
        Where(goqu.C(cnParentIDC).Eq(id)).
        Prepared(true).
        ToSQL()
    if err != nil {
        return nil, err
    }

    q := db.Query{
        Name:     "CategoryChildren",
        QueryRaw: query,
    }

    rows, err := r.db.DB().QueryContext(ctx, q, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var categoryRows []row.CategoryNode

    for rows.Next() {
        var categoryRow row.CategoryNode
        err = rows.Scan(&categoryRow.ID, &categoryRow.Name, &categoryRow.ParentID)
        if err != nil {
            return nil, err
        }
        categoryRows = append(categoryRows, categoryRow)
    }

    return row.ToCategories(categoryRows), nil
}
```

**Repository Responsibilities:**
- Group table (`goqu.T`) and column constants (`C` suffix) at top of file
- Build SQL queries using goqu with `dialect.Select/Insert/Update/Delete`
- Use `.Prepared(true)` for parameterized queries
- Use `goqu.C(columnC).Eq(value)` for WHERE expressions
- Use string constants directly in `Select()`, `Cols()`, `goqu.Record{}`, `alias.Col()`
- Wrap queries in `db.Query{Name: "MethodName", QueryRaw: query}`
- Scan results into `row` structs (NEVER directly into domain models)
- Convert rows to domain models using `row.To{Entity}` functions
- Return domain models to service layer
- NO business logic, NO transactions (transaction is in context)

#### 5. Row Mapping Pattern

**File: `repository/row/product.go`**
```go
// CategoryNode represents category data in database
type CategoryNode struct {
    ID       int
    Name     string
    ParentID *int    // Nullable fields use pointers
    Path     string
}

// FromCategory converts domain category to database row
func FromCategory(category model.Category) CategoryNode {
    return CategoryNode{
        ID:       category.ID,
        Name:     category.Name,
        ParentID: &category.Parent.ID,
    }
}

// ToCategory converts database row to domain category
func ToCategory(row CategoryNode) model.Category {
    if row.ParentID == nil {
        return model.Category{
            ID:   row.ID,
            Name: row.Name,
        }
    }

    return model.Category{
        ID:     row.ID,
        Name:   row.Name,
        Parent: &model.Category{ID: *row.ParentID},
    }
}

// ToCategories converts database rows to domain categories
func ToCategories(rows []CategoryNode) []model.Category {
    categories := make([]model.Category, len(rows))
    for i, row := range rows {
        categories[i] = ToCategory(row)
    }
    return categories
}
```

**Row Mapping Responsibilities:**
- Mirror database schema exactly (nullable columns = pointer fields)
- Provide `From{Entity}` (domain → row) and `To{Entity}` (row → domain) converters
- Handle null values using pointers
- Provide batch converters (e.g., `ToCategories` for `[]CategoryNode`)
- Keep conversion logic simple and pure (no queries, no business logic)

#### 6. Model Layer Pattern

**File: `model/category.go`**
```go
type Category struct {
    ID       int
    Name     string
    Parent   *Category
    Children []Category
}
```

**Model Responsibilities:**
- Pure business entities
- NO JSON tags (use DTOs for serialization)
- NO database tags (use row structs for persistence)
- Relations use pointers for optional fields
- Collections use slices (not pointers to slices)

#### 7. Routes Pattern

**File: `routes.go`**
```go
func RegisterRoutes(router *gin.RouterGroup, h *api.Handlerproduct) {
    productRoutes := router.Group("/product")
    {
        productRoutes.GET("/category/tree", h.GetCategoriesTree)
        productRoutes.GET("/category/:id", h.GetCategory)
    }
}
```

### Interface Segregation Pattern

Each layer defines minimal interfaces for its dependencies:

1. **Handler** defines service interfaces (per handler file)
2. **Service** defines repository interfaces (per service file)
3. **Main service.go** composes all repo interfaces into `ProductRepo`
4. **Main handler.go** composes all service interfaces into `Serviceproduct`

This pattern enables:
- Clear dependency boundaries
- Easy testing (mock minimal interfaces)
- Compile-time verification
- Self-documenting code

### Error Handling Pattern

**Handler Layer:**
```go
if err != nil {
    c.JSON(http.StatusInternalServerError, response.Response[any]{
        Status: "internal error",
        Data:   nil,
    })
    return
}
```

**Service Layer:**
```go
if err != nil {
    slog.ErrorContext(ctx, "MethodName", "error", err)
    return nil, err
}
```

**Repository Layer:**
```go
if err != nil {
    return model.Entity{}, err  // or nil, err for slices
}
```

### Context-First
All service and repository methods accept `context.Context` as the first parameter for cancellation, timeouts, and transaction propagation.

### Graceful Shutdown
The `pkg/closer/` package manages graceful shutdown by collecting cleanup functions (e.g., database connection close) that execute on application termination.

## Testing

Tests are written using `testify` for assertions and `gofakeit` for generating test data. See `pkg/db/prettier/query_prettier_test.go` for examples.

## Common Gotcalls

### Database Transactions
Always use the transaction manager (`db.TxManager`) instead of starting transactions directly. The manager handles nested transactions correctly by reusing the transaction from context.

### Query Placeholders
PostgreSQL queries use dollar sign placeholders (`$1`, `$2`, etc.). Set `PlaceholderFormat(sq.Dollar)` when using Squirrel query builder.

### Import Paths
All imports use the module name `baru` as defined in `go.mod`. Use absolute imports like `baru/internal/config` rather than relative imports.

### Logging
Use structured logging via `log/slog`. Error contexts can be enriched using `logger.ErrorCtx(ctx, err)`.