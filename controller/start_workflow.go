package controller

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"go.temporal.io/sdk/client"
)

type StartGenericWorkflowRequest struct {
	WorkflowName    string
	QueueName       string
	WorkflowRequest interface{}
}

func MakeWorkflowStarterHandleFunc(temporalClient client.Client) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		var err error
		defer func() {
			if err != nil {
				writeError(rw, err)
			}
		}()
		ctx := context.Background()
		startWorkflowReq, err := getStartWorkflowRequest(req)
		if err != nil {
			return
		}
		ctxWithRequestMetadata := injectFromHeaders(ctx, req)
		output, err := executeGenericWorkflow(ctxWithRequestMetadata, temporalClient, startWorkflowReq)
		if err != nil {
			return
		}
		err = writeOutput(rw, output)
		if err != nil {
			return
		}
	}
}

func getStartWorkflowRequest(req *http.Request) (wr StartGenericWorkflowRequest, err error) {
	defer func() {
		if closeErr := req.Body.Close(); closeErr != nil {
			log.Print("error closing HTTP body: " + closeErr.Error())
		}
	}()
	err = json.NewDecoder(req.Body).Decode(&wr)
	return
}

func executeGenericWorkflow(ctx context.Context,
	temporalClient client.Client, startWorkflowReq StartGenericWorkflowRequest) (resp interface{}, err error) {
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: startWorkflowReq.QueueName,
	}
	workflowRun, err := temporalClient.ExecuteWorkflow(ctx, workflowOptions,
		startWorkflowReq.WorkflowName, startWorkflowReq.WorkflowRequest)
	if err != nil {
		return
	}
	var workflowResp interface{}
	err = workflowRun.Get(ctx, &workflowResp)
	if err != nil {
		return
	}
	return workflowResp, nil
}
