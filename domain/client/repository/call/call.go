package call

import model "github.com/jjoc007/poc-api-proxy/domain/model/call"

type Repository interface {
	SaveCall(*model.Call) error
}
