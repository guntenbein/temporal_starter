package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"temporal_microservices"
	"temporal_microservices/context/propagators"
	"temporal_starter"
	"temporal_starter/controller"
	"temporal_starter/workflow"
	"time"

	"github.com/gorilla/mux"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	temporal_workflow "go.temporal.io/sdk/workflow"
)

func main() {
	log.Print("starting TEMPORAL STARTER microservice")

	temporalClient := initTemporalClient()

	worker := initWorkflowWorker(temporalClient)
	defer func() {
		worker.Stop()
	}()

	httpServer := initHTTPServer(controller.MakeStarterHandleFunc(temporalClient))

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-signals

	err := httpServer.Shutdown(context.Background())
	if err != nil {
		log.Fatal("cannot gracefully stop HTTP server: " + err.Error())
	}

	log.Print("closing TEMPORAL STARTER microservice")
}

func initTemporalClient() client.Client {
	temporalClientOptions := client.Options{HostPort: net.JoinHostPort("localhost", "7233"),
		ContextPropagators: []temporal_workflow.ContextPropagator{
			propagators.NewStringMapPropagator([]string{temporal_microservices.ProcessIDContextField}),
			propagators.NewSecretPropagator(propagators.SecretPropagatorConfig{
				Keys:   []string{temporal_microservices.JWTContextField},
				Crypto: propagators.Base64Crypto{},
			}),
		},
	}
	temporalClient, err := client.NewClient(temporalClientOptions)
	if err != nil {
		log.Fatal("cannot start temporal client: " + err.Error())
	}
	return temporalClient
}

func initWorkflowWorker(temporalClient client.Client) worker.Worker {
	workerOptions := worker.Options{
		MaxConcurrentWorkflowTaskExecutionSize: temporal_starter.MaxConcurrentWorkflowSize,
	}
	worker := worker.New(temporalClient, temporal_starter.WorkflowQueue, workerOptions)
	worker.RegisterWorkflow(workflow.StartGenericActivityWorkflow)

	err := worker.Start()
	if err != nil {
		log.Fatal("cannot start temporal worker: " + err.Error())
	}

	return worker
}

func initHTTPServer(singleHandler func(http.ResponseWriter, *http.Request)) *http.Server {
	router := mux.NewRouter()
	router.Methods(http.MethodPost).Path("/").HandlerFunc(singleHandler)
	server := &http.Server{
		Addr:         net.JoinHostPort("localhost", "8080"),
		Handler:      router,
		ReadTimeout:  1 * time.Minute,
		WriteTimeout: 1 * time.Minute,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal("cannot start HTTP server: " + err.Error())
		}
	}()
	return server
}
