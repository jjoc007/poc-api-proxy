package calls_stats_summary

import model "github.com/jjoc007/poc-api-proxy/domain/model/calls_stats_summary"

type Repository interface {
	Save(*model.CallsStatsSummary) error
	GetTotalCalls(*model.CallsStatsSummary) error
}
