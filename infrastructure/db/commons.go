package db

import (
	"database/sql"
	"time"
)

const emptyLength = 0

func CheckDatabaseNullTimeStamp(value time.Time) interface{} {
	if (value.Equal(time.Time{})) {
		return sql.NullTime{}
	}
	return value
}

func CheckDatabaseNullString(value string) interface{} {
	if len(value) == emptyLength {
		return sql.NullString{}
	}
	return value
}

func CheckDatabaseNullUInt64(value uint64) interface{} {
	if value == emptyLength {
		return sql.NullInt64{}
	}
	return value
}

func CheckDatabaseNullInt(value int) interface{} {
	if value == emptyLength {
		return sql.NullInt64{}
	}
	return value
}
