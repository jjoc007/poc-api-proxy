package quota

import model "github.com/jjoc007/poc-api-proxy/domain/model/quota"

type Repository interface {
	GetLimitCalls(*model.Quota) error
}
