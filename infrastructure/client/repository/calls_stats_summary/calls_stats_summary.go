package calls_stats_summary

import (
	"database/sql"
	"github.com/jjoc007/poc-api-proxy/domain/client/repository/calls_stats_summary"
	log "github.com/sirupsen/logrus"
)


func NewCallsStatsSummaryRepository(dbConnection *sql.DB) calls_stats_summary.Repository {
	return &callsStatsSummaryRepository{
		dbConnection: dbConnection,
	}
}

type callsStatsSummaryRepository struct {
	dbConnection *sql.DB
}

func (repository *callsStatsSummaryRepository) LoadRules() {
	log.Info("Start LoadRules")
}