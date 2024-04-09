package fate

import (
	"context"
	"log/slog"
	"time"

	"go.temporal.io/sdk/workflow"
)

const (
	TaskQueue = "fly-autoscaler-example-queue"
)

func ExampleWorkflow(ctx workflow.Context, duration time.Duration) (*Result, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 120 * time.Second,
	})

	var activity *ExampleActivity
	var result *Result
	if err := workflow.ExecuteActivity(ctx, activity.GetData, duration).Get(ctx, &result); err != nil {
		return result, err
	}
	return result, nil
}

type ExampleActivity struct{}

func (a *ExampleActivity) GetData(ctx context.Context, duration time.Duration) (*Result, error) {
	slog.Info("starting activity", slog.Duration("duration", duration))
	time.Sleep(duration)
	slog.Info("activity complete")
	return &Result{Msg: "Hello World"}, nil
}

type Result struct {
	Msg string
}
