package quota

import (
	"database/sql"
	"github.com/jjoc007/poc-api-proxy/domain/client/repository/quota"
	log "github.com/sirupsen/logrus"
)


func NewQuotaRepository(dbConnection *sql.DB) quota.Repository {
	return &quotaRepository{
		dbConnection: dbConnection,
	}
}

type quotaRepository struct {
	dbConnection *sql.DB
}

func (repository *quotaRepository) LoadRules() {
	log.Info("Start LoadRules")
}