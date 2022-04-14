package db

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
)

type Connection interface {
	GetWriteConnection() *sql.DB
	GetReadConnection() *sql.DB
}

func CloseConnections(err error, tx *sql.Tx, stmt *sql.Stmt, rows *sql.Rows) {
	if tx != nil {
		switch err {
		case nil:
			err := tx.Commit()
			if err != nil {
				log.Error(err.Error(), err)
			}
		default:
			err = tx.Rollback()
			if err != nil {
				log.Error(err.Error(), err)
			}
		}
	}

	if stmt != nil {
		err = stmt.Close()
		if err != nil {
			log.Error(err.Error(), err)
		}
	}

	if rows != nil {
		err = rows.Close()
		if err != nil {
			log.Error(err.Error(), err)
		}
	}
}
