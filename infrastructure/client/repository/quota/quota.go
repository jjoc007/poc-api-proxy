package quota

import (
	"database/sql"
	"fmt"
	"github.com/jjoc007/poc-api-proxy/domain/client/repository/quota"
	model "github.com/jjoc007/poc-api-proxy/domain/model/quota"
	"github.com/jjoc007/poc-api-proxy/utils/errors"
)

func NewQuotaRepository(rConnection, wConnection *sql.DB) quota.Repository {
	return &quotaRepository{
		rConnection: rConnection,
		wConnection: wConnection,
	}
}

type quotaRepository struct {
	rConnection *sql.DB
	wConnection *sql.DB
}

func (repository *quotaRepository) GetLimitCalls(quota *model.Quota) (err error) {
	queryParams := []interface{}{quota.SourceIP, quota.TargetPath}
	err = repository.rConnection.QueryRow(selectQuota, queryParams...).Scan(&quota.LimitCalls)
	if err != nil && err.Error() == sql.ErrNoRows.Error() {
		return errors.NewBadRequestError(fmt.Sprintf("Quota for IP: %s Path: %s Not Found", quota.SourceIP, quota.TargetPath))
	}
	return
}
