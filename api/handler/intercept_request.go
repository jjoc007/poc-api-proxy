package handler

import (
	"github.com/jjoc007/poc-api-proxy/domain/service/intercept"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func InterceptRequest(interceptService intercept.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Start InterceptRequest")
		err := interceptService.InterceptRequest(&w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
		}
	}
}
