module temporal_starter

go 1.14

require (
	github.com/golangci/golangci-lint v1.34.1 // indirect
	github.com/gorilla/mux v1.8.0
	go.temporal.io/api v1.4.0
	go.temporal.io/sdk v1.3.0
	temporal_microservices v0.0.0
)

replace temporal_microservices v0.0.0 => github.com/guntenbein/temporal_microservices v0.0.0-20210121102849-3d275f5cf5d1
