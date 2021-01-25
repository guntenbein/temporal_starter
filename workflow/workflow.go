package workflow

import (
	temporal_microservices_workflow "temporal_microservices/domain/workflow"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func StartGenericActivityWorkflow(ctx workflow.Context, activityName, queueName string, req interface{}) (resp interface{}, err error) {
	if req == nil {
		err = temporal_microservices_workflow.BusinessError{Message: "there no activity request to process"}
		return
	}
	ctx = withActivityOptions(ctx, queueName)
	err = workflow.ExecuteActivity(ctx, activityName, req).Get(ctx, &resp)
	return
}

func withActivityOptions(ctx workflow.Context, queue string) workflow.Context {
	ao := workflow.ActivityOptions{
		TaskQueue:              queue,
		ScheduleToStartTimeout: 24 * time.Hour,
		StartToCloseTimeout:    24 * time.Hour,
		HeartbeatTimeout:       time.Second * 5,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:        time.Second,
			BackoffCoefficient:     2.0,
			MaximumInterval:        time.Minute * 5,
			NonRetryableErrorTypes: []string{"BusinessError"},
		},
	}
	ctxOut := workflow.WithActivityOptions(ctx, ao)
	return ctxOut
}
