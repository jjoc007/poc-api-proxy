package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jjoc007/poc-api-proxy/config"
	database "github.com/jjoc007/poc-api-proxy/infrastructure/db"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	connectionUrl            = "%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true"
	connectionErrorMessage   = "error when trying to connect with mysql database"
	notAvailableErrorMessage = "the database is not available or accessible"
	driverTypeMysql          = "mysql"
	maxOpenConns             = 30
	maxIdleConns             = 20
	connMaxLifetime          = 250
)

func NewMysqlConnection(connectionData *config.ConnectionData) database.Connection {
	log.Info("Start bootstrap app DB objects")
	writeConnString := fmt.Sprintf(connectionUrl,
		connectionData.WriteUser,
		connectionData.WritePwd,
		connectionData.Host,
		connectionData.Schema)

	writeConnection, err := getConnection(writeConnString)
	if err != nil {
		log.Error(connectionErrorMessage, err)
	}

	readConnString := fmt.Sprintf(connectionUrl,
		connectionData.ReadUser,
		connectionData.ReadPwd,
		connectionData.Host,
		connectionData.Schema)

	readConnection, err := getConnection(readConnString)
	if err != nil {
		log.Error(connectionErrorMessage, err)
	}

	log.Info("End bootstrap app DB objects")

	return &mysqlConnection{
		writeConnection: writeConnection,
		readConnection:  readConnection,
	}

}

type mysqlConnection struct {
	writeConnection *sql.DB
	readConnection  *sql.DB
}

// GetWriteConnection database connection client
func (m *mysqlConnection) GetWriteConnection() *sql.DB {
	return m.writeConnection
}

// GetReadConnection database connection client
func (m *mysqlConnection) GetReadConnection() *sql.DB {
	return m.readConnection
}

func getConnection(connectionString string) (sqlDB *sql.DB, err error) {
	sqlDB, err = sql.Open(driverTypeMysql, connectionString)
	if err != nil {
		log.Error(connectionErrorMessage, err)
		return
	}
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetConnMaxLifetime(connMaxLifetime * time.Millisecond)
	err = sqlDB.Ping()
	if err != nil {
		log.Error(notAvailableErrorMessage, err)
		return
	}
	log.Infof("Connection stats: %+v", sqlDB.Stats())
	return
}
