package main

import (
	"context"
	"time"

	logsuppress "github.com/Dendromeda/suppresserrors"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

var TaskQueue = "SUPPRESS-LOGS-TQ"

func main() {

	ctx := context.Background()

	temporalClient, err := client.Dial(client.Options{
		Logger: logsuppress.NewLoggerWithSuppressedTypes(ErrSuppressed),
	})
	if err != nil {
		panic(err)
	}
	err = StartWorker(temporalClient)
	if err != nil {
		panic(err)
	}

	wf, err := temporalClient.ExecuteWorkflow(ctx, client.StartWorkflowOptions{
		TaskQueue:                TaskQueue,
		WorkflowExecutionTimeout: time.Second * 10,
	}, Workflow)
	if err != nil {
		panic(err)
	}

	err = wf.Get(ctx, nil)
	if err != nil {
		panic(err)
	}

}

func StartWorker(temporalClient client.Client) error {

	worker := worker.New(temporalClient, TaskQueue, worker.Options{})

	worker.RegisterActivity(Activity)
	worker.RegisterWorkflow(Workflow)

	return worker.Start()

}

func Workflow(ctx workflow.Context) error {

	err := workflow.ExecuteActivity(workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		ScheduleToCloseTimeout: time.Second * 10,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts:    7,
			BackoffCoefficient: 1.0,
			InitialInterval:    time.Millisecond * 100,
		},
	}), Activity).Get(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}

var count = 0
var ErrSuppressed = "SuppressedError"
var ErrLogged = "LoggedError"

func Activity(ctx context.Context) error {
	count++
	if count%2 == 0 {
		return temporal.NewApplicationError("suppress this!", ErrSuppressed)
	} else {
		return temporal.NewApplicationError("log this!", ErrLogged)
	}

}
