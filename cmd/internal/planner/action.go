package planner

// ActionType represents the type of plan action
type ActionType string

const (
	CreateSchema ActionType = "CREATE_SCHEMA"
	CreateTable  ActionType = "CREATE_TABLE"
	AddColumn    ActionType = "ADD_COLUMN"
)

// Action represents a single planned operation
type Action struct {
	Type        ActionType
	Description string
	SQL         string
}
