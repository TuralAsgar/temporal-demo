package main

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func LeadCallingWorkflow(ctx workflow.Context, lead Lead) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("Workflow started", "leadID", lead.ID)
	var callID string
	var callResult string

	// activity options
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 5 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 4.0,
			MaximumInterval:    1 * time.Minute,
			MaximumAttempts:    4,
			NonRetryableErrorTypes: []string{
				"ValidationError",
			},
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	err := workflow.ExecuteActivity(ctx, PickLeadActivity, lead).Get(ctx, &callID)
	if err != nil {
		return err
	}

	err = workflow.ExecuteActivity(ctx, CallCustomerActivity, lead).Get(ctx, &callID)
	if err != nil {
		return err
	}

	workflow.GetSignalChannel(ctx, "call-result").Receive(ctx, &callResult)

	err = workflow.ExecuteActivity(ctx, UpdateLeadActivity, lead.ID, callResult).Get(ctx, nil)
	if err != nil {
		return err
	}

	err = workflow.ExecuteActivity(ctx, CreateInteractionActivity, lead.ID, callResult).Get(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}
