# Go Logger

Một thư viện logging mạnh mẽ và linh hoạt cho Go, được xây dựng trên nền tảng [Uber Zap](https://github.com/uber-go/zap) với nhiều tính năng nâng cao như log rotation, cấu hình từ environment variables, và hỗ trợ nhiều định dạng output.

## Tính năng

- ✅ **High Performance**: Dựa trên Uber Zap - một trong những logger nhanh nhất cho Go
- ✅ **Structured Logging**: Hỗ trợ structured logging với JSON và console format
- ✅ **Log Rotation**: Hỗ trợ rotation theo size, time hoặc cả hai
- ✅ **Environment Configuration**: Cấu hình dễ dàng qua environment variables
- ✅ **Multiple Output**: Ghi log đồng thời ra console và file
- ✅ **Flexible Configuration**: Builder pattern cho cấu hình linh hoạt
- ✅ **Production Ready**: Tối ưu cho cả development và production
- ✅ **Interface Based**: Dễ dàng mock và test

## Cài đặt

```bash
go get github.com/csmart-libs/go-logger
```

## Sử dụng cơ bản

### 1. Sử dụng Global Logger (Đơn giản nhất)

```go
package main

import (
    "github.com/csmart-libs/go-logger"
)

func main() {
    // Sử dụng trực tiếp với cấu hình mặc định
    logger.Info("Application started")
    logger.Debug("Debug message")
    logger.Warn("Warning message")
    logger.Error("Error message")

    // Với structured fields
    logger.Info("User logged in",
        logger.String("user_id", "12345"),
        logger.String("ip", "192.168.1.1"),
        logger.Int("attempt", 1),
    )
}
```

### 2. Khởi tạo với cấu hình tùy chỉnh

```go
package main

import (
    "github.com/csmart-libs/go-logger"
)

func main() {
    // Cấu hình cho development
    config := logger.DevelopmentConfig()

    // Hoặc cấu hình cho production
    // config := logger.ProductionConfig()

    // Khởi tạo global logger
    if err := logger.Initialize(config); err != nil {
        panic(err)
    }

    logger.Info("Logger initialized successfully")
}
```

### 3. Tạo logger instance riêng

```go
package main

import (
    "github.com/csmart-libs/go-logger"
)

func main() {
    // Tạo cấu hình
    config := logger.DefaultConfig().
        WithLevel("debug").
        WithEnvironment("development").
        WithFileOutput("logs/app.log")

    // Tạo logger instance
    log, err := logger.NewLogger(config)
    if err != nil {
        panic(err)
    }

    log.Info("This is a custom logger instance")
}
```

## Cấu hình nâng cao

### 1. Cấu hình với Builder Pattern

```go
config := logger.DefaultConfig().
    WithLevel("info").
    WithEnvironment("production").
    WithEncoding("json").
    WithFileOutput("logs/app.log").
    WithFileRotation(100, 30, 10). // maxSize: 100MB, maxAge: 30 days, maxBackups: 10
    WithFileCompression(true).
    WithLocalTime(true).
    WithDailyRotation()

if err := logger.Initialize(config); err != nil {
    panic(err)
}
```

### 2. Cấu hình Log Rotation

#### Size-based Rotation (Mặc định)
```go
config := logger.DefaultConfig().
    WithFileOutput("logs/app.log").
    WithFileRotation(100, 30, 10) // 100MB, 30 days, 10 backups

logger.Initialize(config)
```

#### Time-based Rotation
```go
config := logger.DefaultConfig().
    WithFileOutput("logs/app.log").
    WithDailyRotation() // Rotate hàng ngày

// Hoặc các tùy chọn khác:
// WithHourlyRotation() - Rotate hàng giờ
// WithWeeklyRotation() - Rotate hàng tuần
// WithMonthlyRotation() - Rotate hàng tháng

logger.Initialize(config)
```

#### Combined Rotation (Size + Time)
```go
config := logger.DefaultConfig().
    WithFileOutput("logs/app.log").
    WithBothRotation(100, 30, 10, logger.RotationDaily)

logger.Initialize(config)
```

### 3. Cấu hình từ Environment Variables

```go
// Sử dụng cấu hình từ environment variables
config := logger.ConfigFromEnv()
logger.Initialize(config)
```

Các environment variables được hỗ trợ:

```bash
# Cấu hình cơ bản
export APP_ENV=production          # development, staging, production, test
export LOG_LEVEL=info             # debug, info, warn, error, fatal, panic
export LOG_ENCODING=json          # json, console
export LOG_OUTPUT_PATHS=stdout    # stdout hoặc file paths (phân cách bằng dấu phẩy)

# Cấu hình file
export LOG_FILE=logs/app.log
export LOG_FILE_MAX_SIZE=100      # MB
export LOG_FILE_MAX_AGE=30        # days
export LOG_FILE_MAX_BACKUPS=10
export LOG_FILE_LOCAL_TIME=true
export LOG_FILE_COMPRESS=true
export LOG_FILE_CREATE_DIR=true

# Cấu hình rotation
export LOG_FILE_ROTATION_MODE=size    # size, time, both
export LOG_FILE_TIME_INTERVAL=daily   # hourly, daily, weekly, monthly
export LOG_FILE_TIME_FORMAT=2006-01-02
```

## Các loại cấu hình có sẵn

### 1. Development Config
```go
config := logger.DevelopmentConfig()
// Level: debug
// Environment: development
// Output: stdout
// Encoding: console (có màu sắc)
```

### 2. Production Config
```go
config := logger.ProductionConfig()
// Level: info
// Environment: production
// Output: stdout + file
// Encoding: json
// File rotation: enabled
```

### 3. Test Config
```go
config := logger.TestConfig()
// Level: error
// Environment: test
// Output: stdout only
// Encoding: console
```

## Structured Logging

### Sử dụng các field helpers

```go
logger.Info("User action",
    logger.String("user_id", "12345"),
    logger.String("action", "login"),
    logger.Int("attempt", 1),
    logger.Int64("timestamp", time.Now().Unix()),
    logger.Float64("duration", 1.23),
    logger.Bool("success", true),
    logger.Any("metadata", map[string]interface{}{
        "ip": "192.168.1.1",
        "user_agent": "Mozilla/5.0...",
    }),
)
```

### Error logging
```go
if err != nil {
    logger.Error("Database connection failed",
        logger.Err(err),
        logger.String("database", "postgres"),
        logger.String("host", "localhost"),
    )
}
```

### Child logger với context
```go
// Tạo child logger với context cố định
userLogger := logger.With(
    logger.String("user_id", "12345"),
    logger.String("session_id", "abcdef"),
)

// Sử dụng child logger
userLogger.Info("User performed action", logger.String("action", "purchase"))
userLogger.Warn("User exceeded rate limit")
```

## Dependency Injection

Thư viện cung cấp interface `Logger` để dễ dàng sử dụng với dependency injection:

```go
type UserService struct {
    logger logger.Logger
}

func NewUserService(log logger.Logger) *UserService {
    return &UserService{
        logger: log,
    }
}

func (s *UserService) CreateUser(user User) error {
    s.logger.Info("Creating user", logger.String("email", user.Email))

    // Business logic...

    s.logger.Info("User created successfully", logger.String("user_id", user.ID))
    return nil
}

// Sử dụng
func main() {
    log := logger.GetLogger()
    userService := NewUserService(log)

    // ...
}
```

## Testing

Để test code sử dụng logger, bạn có thể mock interface `Logger`:

```go
type MockLogger struct{}

func (m *MockLogger) Debug(msg string, fields ...zap.Field) {}
func (m *MockLogger) Info(msg string, fields ...zap.Field) {}
func (m *MockLogger) Warn(msg string, fields ...zap.Field) {}
func (m *MockLogger) Error(msg string, fields ...zap.Field) {}
func (m *MockLogger) Fatal(msg string, fields ...zap.Field) {}
func (m *MockLogger) Panic(msg string, fields ...zap.Field) {}
func (m *MockLogger) With(fields ...zap.Field) logger.Logger { return m }
func (m *MockLogger) Sync() error { return nil }

// Sử dụng trong test
func TestUserService(t *testing.T) {
    mockLogger := &MockLogger{}
    userService := NewUserService(mockLogger)

    // Test logic...
}
```

## Ví dụ hoàn chỉnh

```go
package main

import (
    "os"
    "time"

    "github.com/csmart-libs/go-logger"
)

func main() {
    // Cấu hình logger cho production
    config := logger.DefaultConfig().
        WithLevel("info").
        WithEnvironment("production").
        WithEncoding("json").
        WithFileOutput("logs/app.log").
        WithFileRotation(100, 30, 5).
        WithFileCompression(true).
        WithDailyRotation()

    // Khởi tạo global logger
    if err := logger.Initialize(config); err != nil {
        panic(err)
    }

    // Đảm bảo flush logs khi thoát
    defer logger.Sync()

    // Logging cơ bản
    logger.Info("Application started")

    // Structured logging
    logger.Info("Processing request",
        logger.String("method", "GET"),
        logger.String("path", "/api/users"),
        logger.Int("status", 200),
        logger.Float64("duration", 0.123),
    )

    // Error logging
    if err := someOperation(); err != nil {
        logger.Error("Operation failed",
            logger.Err(err),
            logger.String("operation", "database_query"),
        )
    }

    // Child logger với context
    requestLogger := logger.With(
        logger.String("request_id", "req-123"),
        logger.String("user_id", "user-456"),
    )

    requestLogger.Info("Processing user request")
    requestLogger.Warn("Rate limit approaching")

    logger.Info("Application shutting down")
}

func someOperation() error {
    // Simulate some operation that might fail
    return nil
}
```

## API Reference

### Logger Interface

```go
type Logger interface {
    Debug(msg string, fields ...zap.Field)
    Info(msg string, fields ...zap.Field)
    Warn(msg string, fields ...zap.Field)
    Error(msg string, fields ...zap.Field)
    Fatal(msg string, fields ...zap.Field)
    Panic(msg string, fields ...zap.Field)
    With(fields ...zap.Field) Logger
    Sync() error
}
```

### Global Functions

- `Initialize(config Config) error` - Khởi tạo global logger
- `GetLogger() Logger` - Lấy global logger instance
- `NewLogger(config Config) (Logger, error)` - Tạo logger instance mới
- `Debug/Info/Warn/Error/Fatal/Panic(msg string, fields ...zap.Field)` - Global logging functions
- `With(fields ...zap.Field) Logger` - Tạo child logger với context
- `Sync() error` - Flush buffered logs

### Configuration Functions

- `DefaultConfig() Config` - Cấu hình mặc định
- `DevelopmentConfig() Config` - Cấu hình cho development
- `ProductionConfig() Config` - Cấu hình cho production
- `TestConfig() Config` - Cấu hình cho testing
- `ConfigFromEnv() Config` - Cấu hình từ environment variables

### Field Helpers

- `String(key, val string) zap.Field`
- `Int(key string, val int) zap.Field`
- `Int64(key string, val int64) zap.Field`
- `Uint/Uint32/Uint64(key string, val uint) zap.Field`
- `Float64(key string, val float64) zap.Field`
- `Bool(key string, val bool) zap.Field`
- `Any(key string, val any) zap.Field`
- `Err(err error) zap.Field`
- `Duration(key string, val any) zap.Field`

## Requirements

- Go 1.24.4 hoặc cao hơn
- Dependencies:
  - `go.uber.org/zap v1.27.0`
  - `gopkg.in/natefinch/lumberjack.v2 v2.2.1`

## License

MIT License

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.