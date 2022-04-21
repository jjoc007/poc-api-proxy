package api

import (
	"net/http"

	"github.com/jjoc007/poc-api-proxy/api/handler"
	log "github.com/sirupsen/logrus"
)

func StartApp() {
	dependencies := BuildDependencies()
	http.HandleFunc("/", handler.InterceptRequest(dependencies.InterceptService))
	log.Fatal(http.ListenAndServe(":8082", nil))
}
