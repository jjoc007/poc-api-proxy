package api

import (
	"github.com/jjoc007/poc-api-proxy/config"
	"github.com/jjoc007/poc-api-proxy/domain/service/intercept"
	callRepositoryImpl "github.com/jjoc007/poc-api-proxy/infrastructure/client/repository/call"
	callsStatsSummaryRepositoryImpl "github.com/jjoc007/poc-api-proxy/infrastructure/client/repository/calls_stats_summary"
	quotaRepositoryImpl "github.com/jjoc007/poc-api-proxy/infrastructure/client/repository/quota"
	"github.com/jjoc007/poc-api-proxy/infrastructure/db/mysql"
)

type Dependencies struct {
	InterceptService intercept.Service
}

func BuildDependencies() Dependencies {
	dataBaseConnection := mysql.NewMysqlConnection(&config.DataBaseDefinition)
	quotaRepository := quotaRepositoryImpl.NewQuotaRepository(dataBaseConnection.GetReadConnection(), dataBaseConnection.GetWriteConnection())
	callsStatsSummaryRepository := callsStatsSummaryRepositoryImpl.NewCallsStatsSummaryRepository(dataBaseConnection.GetReadConnection(), dataBaseConnection.GetWriteConnection())
	callRepository := callRepositoryImpl.NewCallRepository(dataBaseConnection.GetReadConnection(), dataBaseConnection.GetWriteConnection())
	interceptRequestService := intercept.NewInterceptRequestService(
		quotaRepository,
		callRepository,
		callsStatsSummaryRepository,
	)

	return Dependencies{
		InterceptService: interceptRequestService,
	}
}
