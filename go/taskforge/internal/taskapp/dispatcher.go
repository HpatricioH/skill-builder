package taskapp

import "context"

type Dispatcher interface {
	Dispatch(ctx context.Context, events []Event)
}
