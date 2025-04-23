package tasks

import (
	"time"
)

// LongComputationHandler обрабатывает длительные вычисления
type LongComputationHandler struct{}

func NewLongComputationHandler() *LongComputationHandler {
	return &LongComputationHandler{}
}

func (h *LongComputationHandler) Execute(params map[string]interface{}) (map[string]interface{}, error) {
	// Здесь реализуется длительная I/O bound задача

	// Эмулируется длительную операцию
	duration, ok := params["duration"].(float64)
	if !ok {
		duration = 180 // Значение по умолчанию - 3 минуты, если не было указано иного
	}

	time.Sleep(time.Duration(duration) * time.Second)

	result := map[string]interface{}{
		"message":  "Long computation completed successfully",
		"duration": duration,
	}

	return result, nil
}
