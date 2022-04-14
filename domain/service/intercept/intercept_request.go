package intercept

import (
	"bytes"
	"fmt"
	"github.com/jjoc007/poc-api-proxy/config"
	"github.com/jjoc007/poc-api-proxy/domain/client/repository/quota"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	url2 "net/url"
	"strings"
	"time"
)

var RequiredHeaders = []string{
	config.XSourceIPHeader,
}


type Service interface {
	InterceptRequest(*http.ResponseWriter, *http.Request) error
}

func NewInterceptRequestService(quotaRepository quota.Repository) Service {
	return &interceptRequestService{
		quotaRepository: quotaRepository,
		httpClient: &http.Client{},
	}
}

type interceptRequestService struct {
	quotaRepository quota.Repository
	httpClient              *http.Client
}

func (service *interceptRequestService) InterceptRequest(w *http.ResponseWriter, r *http.Request) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(*w, err.Error(), http.StatusInternalServerError)
		return nil
	}

	r.Body = ioutil.NopCloser(bytes.NewReader(body))

	// validators
	valid, headersNotFound := validateRequiredHeaders(r.Header)
	if !valid {
		http.Error(*w, fmt.Sprintf("Headers Required: [%s]", strings.Join(headersNotFound, ",")), http.StatusBadGateway)
		return nil
	}


	u, _ :=url2.Parse(r.RequestURI)
	endpointURL := fmt.Sprintf("%s%s", config.BaseURLBackendServer, u.RawPath)
	proxyReq, err := http.NewRequest(r.Method, endpointURL, bytes.NewReader(body))
	proxyReq.Header = make(http.Header)
	for h, val := range r.Header {
		proxyReq.Header[h] = val
	}

	start := time.Now()
	resp, err := service.httpClient.Do(proxyReq)
	log.Info("Everything:", time.Since(start).Nanoseconds())
	defer resp.Body.Close()
	if err != nil {
		http.Error(*w, err.Error(), http.StatusBadGateway)
		return nil
	}

	// tomar tiempo de ejecucion
	// enviar metricas a APM o tabla
	// enviar contador de reglas ip origen y path destino

	copyHeader((*w).Header(), resp.Header)
	(*w).WriteHeader(resp.StatusCode)
	io.Copy(*w, resp.Body)

	return nil
}


func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func validateRequiredHeaders(headers http.Header) (bool, []string) {
	headersNotFound := make([]string, 0)
	for _, header := range RequiredHeaders {
		found := false
		for h, _ := range headers {
			if strings.EqualFold(h,header) {
				found = true
				break
			}
		}
		if !found {
			headersNotFound = append(headersNotFound, header)
		}
	}

	return len(headersNotFound) == 0, headersNotFound
}