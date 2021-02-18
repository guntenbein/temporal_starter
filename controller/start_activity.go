package controller

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"temporal_starter"
	"temporal_starter/workflow"

	"go.temporal.io/sdk/client"
)

type StartGenericActivityRequest struct {
	ActivityName    string
	QueueName       string
	ActivityRequest interface{}
}

func MakeActivityStarterHandleFunc(temporalClient client.Client) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		var err error
		defer func() {
			if err != nil {
				writeError(rw, err)
			}
		}()
		ctx := context.Background()
		startActivityRequest, err := getStartActivityRequest(req)
		if err != nil {
			return
		}
		ctxWithRequestMetadata := injectFromHeaders(ctx, req)
		output, err := executeSingleActivityWorkflow(ctxWithRequestMetadata, temporalClient, startActivityRequest)
		if err != nil {
			return
		}
		err = writeOutput(rw, output)
		if err != nil {
			return
		}
	}
}

func getStartActivityRequest(req *http.Request) (wr StartGenericActivityRequest, err error) {
	defer func() {
		if closeErr := req.Body.Close(); closeErr != nil {
			log.Print("error closing HTTP body: " + closeErr.Error())
		}
	}()
	err = json.NewDecoder(req.Body).Decode(&wr)
	return
}

func executeSingleActivityWorkflow(ctx context.Context,
	temporalClient client.Client, startActivityReq StartGenericActivityRequest) (resp interface{}, err error) {
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: temporal_starter.WorkflowQueue,
	}
	workflowRun, err := temporalClient.ExecuteWorkflow(ctx, workflowOptions, workflow.StartGenericActivityWorkflow,
		startActivityReq.ActivityName, startActivityReq.QueueName, startActivityReq.ActivityRequest)
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
