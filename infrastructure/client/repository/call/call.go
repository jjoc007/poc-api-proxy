package call

import (
	"database/sql"
	"github.com/jjoc007/poc-api-proxy/domain/client/repository/call"
	log "github.com/sirupsen/logrus"
)


func NewCallRepository(dbConnection *sql.DB) call.Repository {
	return &callRepository{
		dbConnection: dbConnection,
	}
}

type callRepository struct {
	dbConnection *sql.DB
}

func (repository *callRepository) LoadRules() {
	log.Info("Start LoadRules")
}