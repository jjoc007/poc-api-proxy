package call

import (
	"database/sql"
	"github.com/google/uuid"
	apierrors "github.com/jjoc007/poc-api-proxy/api/error"
	"github.com/jjoc007/poc-api-proxy/domain/client/repository/call"
	model "github.com/jjoc007/poc-api-proxy/domain/model/call"
	"github.com/jjoc007/poc-api-proxy/infrastructure/db"
	"github.com/jjoc007/poc-api-proxy/utils/errors"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	dataBaseBeginTransactionError = "an error happened trying begin transaction with the data base"
	errorInsertCall               = "error when trying to insert call"
)

func NewCallRepository(rConnection, wConnection *sql.DB) call.Repository {
	return &callRepository{
		rConnection: rConnection,
		wConnection: wConnection,
	}
}

type callRepository struct {
	rConnection *sql.DB
	wConnection *sql.DB
}

func (repository *callRepository) SaveCall(call *model.Call) (err error) {
	log.Info("Start SaveCall")
	var transaction *sql.Tx
	guid := uuid.New()

	defer func() { db.CloseConnections(err, transaction, nil, nil) }()

	if transaction, err = repository.getTransaction(); err != nil {
		return err
	}

	_, errIns := transaction.Exec(insertCall, guid.String(), db.CheckDatabaseNullString(call.SourceIP),
		db.CheckDatabaseNullString(call.TargetPath), db.CheckDatabaseNullUInt64(call.Duration),
		time.Now())

	if errIns != nil {
		log.Error(errorInsertCall, errIns)
		err = errors.NewInternalServerError(errorInsertCall)
	}

	return err
}

func (repository *callRepository) getTransaction() (transaction *sql.Tx, err error) {
	if transaction, err = repository.wConnection.Begin(); err != nil {
		log.Error(dataBaseBeginTransactionError, err)
		err = apierrors.NewInternalServerApiError(dataBaseBeginTransactionError, nil)
		return nil, err
	}

	return transaction, nil
}
