package executor

import (
	"encoding/json"
	"fmt"

	"github.com/harshit3011/Job-Queue-Task-Scheduler-Go/internal/job"
)

type Handler func(json.RawMessage) error

type Executor struct{
	handlers map[string]Handler
}

func NewExecutor() *Executor{
	return &Executor{
		handlers: make(map[string]Handler),
	}
}

func (e Executor) RegisterExecutor(jobType string, handler Handler) error{
	if jobType==""{
		return fmt.Errorf("Invalid type")
	}
	e.handlers[jobType]=handler
	return nil
}
func (e Executor) Execute(job job.Job) error{
	handler, exists:=e.handlers[job.Type]
	if !exists{
		return fmt.Errorf("Invalid job type")
	}
	return handler(job.Payload)
}

// executorInstance := executor.NewExecutor()

// executorInstance.Register("PRINT", func(payload json.RawMessage) error {
// 	var data struct {
// 		Message string `json:"message"`
// 	}

// 	if err := json.Unmarshal(payload, &data); err != nil {
// 		return err
// 	}

// 	fmt.Println(data.Message)
// 	return nil
// })