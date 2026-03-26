package taskapp

type EventType string

const (
	EventTaskCompleted EventType = "task_completed"
)

type Event struct {
	Type   EventType
	TaskID int
}
