package api

import (
	"github.com/jjoc007/poc-api-proxy/config"
	"github.com/jjoc007/poc-api-proxy/domain/service/intercept"
	quotaRepositoryImpl "github.com/jjoc007/poc-api-proxy/infrastructure/client/repository/quota"
	"github.com/jjoc007/poc-api-proxy/infrastructure/db/mysql"
)

type Dependencies struct {
	InterceptService intercept.Service
}

func BuildDependencies() Dependencies {
	dataBaseConnection := mysql.NewMysqlConnection(&config.DataBaseDefinition)
	quotaRepository := quotaRepositoryImpl.NewQuotaRepository(dataBaseConnection.GetReadConnection())
	interceptRequestService := intercept.NewInterceptRequestService(quotaRepository)

	return Dependencies{
		InterceptService: interceptRequestService,
	}
}