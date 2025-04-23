# Long Tasks Service

Сервис для асинхронной обработки длительных I/O bound задач с HTTP API интерфейсом.

## Описание

Long Tasks Service - это система обработки длительных операций (3-5 минут и более), которая позволяет асинхронно запускать задачи и получать их результаты через REST API. Сервис разработан с учетом расширяемости и возможности добавления новых типов задач без изменения основной архитектуры.

## Особенности

- Асинхронная обработка длительных задач
- REST API для создания задач и получения результатов
- Расширяемая система типов задач через реестр обработчиков
- Пул воркеров для параллельного выполнения задач
- Отслеживание статуса выполнения задач
- Graceful shutdown для корректного завершения работы

## Системные требования

- Go 1.19 или выше
- Внешние зависимости:
  - github.com/go-chi/chi/v5
  - github.com/google/uuid

## Установка

```bash
# Клонирование репозитория
git clone https://github.com/MR-d3R/LongTaskRunner
cd LongTaskRunner

# Установка зависимостей
go mod download

# Сборка проекта
go build -o longtasksrunner-service
```

## Запуск

```bash
./longtasksrunner-service
```

По умолчанию сервис запускается на порту 8080. Для изменения порта можно использовать переменные окружения (см. раздел Конфигурация).

## Использование API

### Создание новой задачи

```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{"type": "long_computation", "params": {"duration": 120}}'
```

Пример ответа:

```json
{
  "id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
  "status": "pending"
}
```

### Проверка статуса задачи

```bash
curl -X GET http://localhost:8080/api/v1/tasks/f47ac10b-58cc-4372-a567-0e02b2c3d479
```

Пример ответа:

```json
{
  "id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
  "status": "processing"
}
```

### Получение результата задачи

```bash
curl -X GET http://localhost:8080/api/v1/tasks/f47ac10b-58cc-4372-a567-0e02b2c3d479/result
```

Пример ответа для завершенной задачи:

```json
{
  "task_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
  "status": "completed",
  "result": {
    "message": "Long computation completed successfully",
    "duration": 120
  },
  "created_at": "2023-08-15T12:30:45Z",
  "started_at": "2023-08-15T12:30:47Z",
  "ended_at": "2023-08-15T12:32:47Z"
}
```

## Статусы задач

- `pending` - задача поставлена в очередь, но еще не выполняется
- `processing` - задача в данный момент выполняется
- `completed` - задача успешно завершена
- `failed` - при выполнении задачи произошла ошибка

## Архитектура

Сервис состоит из следующих основных компонентов:

1. **HTTP API Server** - обрабатывает входящие HTTP запросы
2. **Task Queue** - очередь задач для асинхронной обработки
3. **Worker Pool** - пул воркеров для выполнения задач
4. **Result Storage** - хранилище результатов выполнения задач
5. **Task Registry** - реестр типов задач и их обработчиков

## Добавление новых типов задач

Для добавления нового типа задачи необходимо:

1. Создать новый обработчик, реализующий интерфейс `TaskHandler`:

```go
type MyNewTaskHandler struct{}

func NewMyNewTaskHandler() *MyNewTaskHandler {
    return &MyNewTaskHandler{}
}

func (h *MyNewTaskHandler) Execute(params map[string]interface{}) (map[string]interface{}, error) {
    // Реализация обработки задачи
    result := map[string]interface{}{
        "data": "Some result data",
    }
    return result, nil
}
```

2. Зарегистрировать обработчик в реестре задач:

```go
// При инициализации сервиса
registry.RegisterHandler("my_new_task_type", tasks.NewMyNewTaskHandler())
```

## Конфигурация

Сервис поддерживает следующие переменные окружения:

- `PORT` - порт для HTTP сервера (по умолчанию 8080)
- `WORKER_COUNT` - количество воркеров (по умолчанию 10)
- `LOG_LEVEL` - уровень логирования (по умолчанию "info")

## Решение проблем

### Ошибка EOF при создании задачи

Если при создании задачи вы получаете ошибку EOF, проверьте:

1. Правильность формата JSON в запросе
2. Наличие заголовка `Content-Type: application/json`
3. Полноту данных в запросе (должны быть указаны и `type`, и `params`)

Пример корректного запроса:

```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{"type": "long_computation", "params": {"duration": 120}}'
```
