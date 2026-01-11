package executor

import (
	"context"
	"fmt"

	"grimleytk/internal/planner"
)

// Executor defines a contract for executing a plan
type Executor interface {
	// Execute executes a list of planned actions inside a transaction
	Execute(ctx context.Context, actions []planner.Action) error
}

// ExecutionError represents a failure during plan execution
type ExecutionError struct {
	Action planner.Action
	Err    error
}

func (e *ExecutionError) Error() string {
	return fmt.Sprintf(
		"execution failed for action [%s]: %s\nSQL:\n%s",
		e.Action.Type,
		e.Err.Error(),
		e.Action.SQL,
	)
}
