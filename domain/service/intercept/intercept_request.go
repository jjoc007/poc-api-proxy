package intercept

import (
	"bytes"
	"fmt"
	"github.com/jjoc007/poc-api-proxy/config"
	"github.com/jjoc007/poc-api-proxy/domain/client/repository/call"
	"github.com/jjoc007/poc-api-proxy/domain/client/repository/calls_stats_summary"
	"github.com/jjoc007/poc-api-proxy/domain/client/repository/quota"
	model "github.com/jjoc007/poc-api-proxy/domain/model/call"
	model2 "github.com/jjoc007/poc-api-proxy/domain/model/calls_stats_summary"
	model3 "github.com/jjoc007/poc-api-proxy/domain/model/quota"
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

func NewInterceptRequestService(quotaRepository quota.Repository,
	callRepository call.Repository,
	callsStatsSummaryRepository calls_stats_summary.Repository,
) Service {
	return &interceptRequestService{
		quotaRepository:             quotaRepository,
		callRepository:              callRepository,
		callsStatsSummaryRepository: callsStatsSummaryRepository,
		httpClient:                  &http.Client{},
	}
}

type interceptRequestService struct {
	quotaRepository             quota.Repository
	callRepository              call.Repository
	callsStatsSummaryRepository calls_stats_summary.Repository
	httpClient                  *http.Client
}

func (service *interceptRequestService) InterceptRequest(w *http.ResponseWriter, r *http.Request) error {
	body, err := ioutil.ReadAll(r.Body)
	u, _ := url2.Parse(r.RequestURI)
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

	// validate
	quota := &model3.Quota{
		SourceIP:   r.Header.Get(config.XSourceIPHeader),
		TargetPath: u.Path,
	}
	err = service.quotaRepository.GetLimitCalls(quota)
	if err != nil {
		http.Error(*w, err.Error(), http.StatusBadGateway)
		return nil
	}

	// current calls stats summary
	callsStatsSummary := &model2.CallsStatsSummary{
		SourceIP:   r.Header.Get(config.XSourceIPHeader),
		TargetPath: u.Path,
	}

	err = service.callsStatsSummaryRepository.GetTotalCalls(callsStatsSummary)
	if err != nil {
		http.Error(*w, err.Error(), http.StatusBadGateway)
		return nil
	}

	if callsStatsSummary.TotalCalls >= quota.LimitCalls {
		http.Error(*w, fmt.Sprintf("exceeded the total limit of calls ip: %s, path: %s", r.Header.Get(config.XSourceIPHeader), u.Path), http.StatusBadGateway)
		return nil
	}

	endpointURL := fmt.Sprintf("%s%s", config.BaseURLBackendServer, u.Path)
	proxyReq, err := http.NewRequest(r.Method, endpointURL, bytes.NewReader(body))
	proxyReq.Header = make(http.Header)
	for h, val := range r.Header {
		proxyReq.Header[h] = val
	}

	start := time.Now()
	resp, err := service.httpClient.Do(proxyReq)
	timeTaken := time.Since(start).Nanoseconds()
	log.Info("Everything:", timeTaken)
	defer resp.Body.Close()
	if err != nil {
		http.Error(*w, err.Error(), http.StatusBadGateway)
		return nil
	}

	callModel := model.Call{
		SourceIP:   r.Header.Get(config.XSourceIPHeader),
		TargetPath: u.Path,
		Duration:   uint64(timeTaken),
	}

	copyHeader((*w).Header(), resp.Header)
	(*w).WriteHeader(resp.StatusCode)
	io.Copy(*w, resp.Body)

	//insert call
	err = service.callRepository.SaveCall(&callModel)
	if err != nil {
		return err
	}

	callStatsSummary := model2.CallsStatsSummary{
		SourceIP:   r.Header.Get(config.XSourceIPHeader),
		TargetPath: u.Path,
	}

	//insert calls stats summary
	err = service.callsStatsSummaryRepository.Save(&callStatsSummary)
	if err != nil {
		return err
	}

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
			if strings.EqualFold(h, header) {
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
