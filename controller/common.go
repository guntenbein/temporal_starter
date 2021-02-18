package controller

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"temporal_microservices"
)

const bearerPrefix = "bearer"

func injectFromHeaders(ctx context.Context, req *http.Request) context.Context {
	jwt := extractTokenFromHeaders(req, temporal_microservices.AuthorizationHTTPHeader)
	processId := req.Header.Get(temporal_microservices.ProcessIDHTTPHeader)
	//nolint:staticcheck
	ctx1 := context.WithValue(ctx, temporal_microservices.ProcessIDContextField, processId)
	//nolint:staticcheck
	ctxOut := context.WithValue(ctx1, temporal_microservices.JWTContextField, jwt)
	return ctxOut
}

func extractTokenFromHeaders(req *http.Request, tokenHeaderName string) string {
	bearer := req.Header.Get(tokenHeaderName)
	authHeaderParts := strings.Split(bearer, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != bearerPrefix {
		return ""
	}
	return authHeaderParts[1]
}

func writeOutput(rw http.ResponseWriter, output interface{}) (err error) {
	body, err := json.Marshal(output)
	if err != nil {
		return
	}
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
	_, err = rw.Write(body)
	if err != nil {
		return
	}
	return
}

func writeError(rw http.ResponseWriter, err error) {
	log.Print(err.Error())
	rw.WriteHeader(http.StatusInternalServerError)
	if _, errWrite := rw.Write([]byte(err.Error())); errWrite != nil {
		log.Print("error writing the HTTP response: " + errWrite.Error())
	}
}
