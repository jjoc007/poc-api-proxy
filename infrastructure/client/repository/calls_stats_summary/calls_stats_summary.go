package calls_stats_summary

import (
	"database/sql"
	apierrors "github.com/jjoc007/poc-api-proxy/api/error"
	"github.com/jjoc007/poc-api-proxy/domain/client/repository/calls_stats_summary"
	model "github.com/jjoc007/poc-api-proxy/domain/model/calls_stats_summary"
	"github.com/jjoc007/poc-api-proxy/infrastructure/db"
	"github.com/jjoc007/poc-api-proxy/utils/errors"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	dataBaseBeginTransactionError = "an error happened trying begin transaction with the data base"
	errorInsertCallsStatsSummary  = "error when trying to insert calls stats summary"
)

func NewCallsStatsSummaryRepository(rConnection, wConnection *sql.DB) calls_stats_summary.Repository {
	return &callsStatsSummaryRepository{
		rConnection: rConnection,
		wConnection: wConnection,
	}
}

type callsStatsSummaryRepository struct {
	rConnection *sql.DB
	wConnection *sql.DB
}

func (repository *callsStatsSummaryRepository) Save(callsStatsSummary *model.CallsStatsSummary) (err error) {
	log.Info("Start Save calls stats summary")
	var transaction *sql.Tx
	year, month, _ := time.Now().Date()

	defer func() { db.CloseConnections(err, transaction, nil, nil) }()

	if transaction, err = repository.getTransaction(); err != nil {
		return err
	}

	_, errIns := transaction.Exec(insertCallsStatsSummary, db.CheckDatabaseNullString(callsStatsSummary.SourceIP),
		db.CheckDatabaseNullString(callsStatsSummary.TargetPath), db.CheckDatabaseNullInt(year), db.CheckDatabaseNullInt(int(month)))

	if errIns != nil {
		log.Error(errorInsertCallsStatsSummary, errIns)
		err = errors.NewInternalServerError(errorInsertCallsStatsSummary)
	}

	return err
}

func (repository *callsStatsSummaryRepository) GetTotalCalls(callsStatsSummary *model.CallsStatsSummary) (err error) {
	year, month, _ := time.Now().Date()
	queryParams := []interface{}{callsStatsSummary.SourceIP, callsStatsSummary.TargetPath, year, month}
	err = repository.rConnection.QueryRow(selectCallsStatsSummary, queryParams...).Scan(&callsStatsSummary.TotalCalls)
	if err != nil && err.Error() != sql.ErrNoRows.Error() {
		return
	}
	return
}

func (repository *callsStatsSummaryRepository) getTransaction() (transaction *sql.Tx, err error) {
	if transaction, err = repository.wConnection.Begin(); err != nil {
		log.Error(dataBaseBeginTransactionError, err)
		err = apierrors.NewInternalServerApiError(dataBaseBeginTransactionError, nil)
		return nil, err
	}

	return transaction, nil
}
