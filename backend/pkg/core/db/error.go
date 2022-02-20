package db

import (
	"github.com/VividCortex/mysqlerr"
	mysql "github.com/go-sql-driver/mysql"
	"github.com/iancoleman/strcase"
	"regexp"
	"strings"
)

// CollidesWithUniqueIndex returns true if the error is about
// colliding with a unique index when inserting/updating data
func CollidesWithUniqueIndex(err error) bool {
	if mysqlError, ok := err.(*mysql.MySQLError); ok {
		return mysqlError.Number == mysqlerr.ER_DUP_ENTRY
	}

	return false
}

func ReferenceDoesNotExist(err error) bool {
	if mysqlError, ok := err.(*mysql.MySQLError); ok {
		return mysqlError.Number == mysqlerr.ER_NO_REFERENCED_ROW_2
	}

	return false
}

func DeadLockFound(err error) bool {
	if mysqlError, ok := err.(*mysql.MySQLError); ok {
		return mysqlError.Number == mysqlerr.ER_LOCK_DEADLOCK
	}

	return false
}

func CantAggregate2Collations(err error) bool {
	if mysqlError, ok := err.(*mysql.MySQLError); ok {
		return mysqlError.Number == mysqlerr.ER_CANT_AGGREGATE_2COLLATIONS
	}

	return false
}

func QueryTimeout(err error) bool {
	if mysqlError, ok := err.(*mysql.MySQLError); ok {
		return mysqlError.Number == mysqlerr.ER_QUERY_TIMEOUT
	}

	return false
}

func SelectTooBig(err error) bool {
	if mysqlError, ok := err.(*mysql.MySQLError); ok {
		return mysqlError.Number == mysqlerr.ER_TOO_BIG_SELECT
	}

	return false
}

type TableNameExtractErrType int

const (
	TABLE_NAME_EXTRACT_FROM_ERR_REF TableNameExtractErrType = 0
)

// GetTableName looks into a given sql error and, depending on TableNameExtractErrType, extracts the column name
// defaultRet is the table name to fallback on  when the table name can't be extracted
// BE CAUTIOUS not to return table names from database to clients when using this!
func GetTableName(errType TableNameExtractErrType, err error, defaultRet string) string {
	var regexTableName string
	switch errType {
	case TABLE_NAME_EXTRACT_FROM_ERR_REF:
		regexTableName = "REFERENCES `.*?`"
		r := regexp.MustCompile(regexTableName)
		rMatch := r.FindString(err.Error())

		if rMatch == "" {
			return defaultRet
		}

		// go doesn't support lookaround, so we can't extract the table name from the regex directly sadly. We have to trim around the table name
		tableName := strings.TrimPrefix(strings.TrimSuffix(rMatch, "`"), "REFERENCES `")

		return tableName
	}

	return defaultRet
}

// GetTableNameCamelCase calls GetTableName and parses the string to camelCase
// defaultRet is the table name to fallback on  when the table name can't be extracted. defaultRet doesn't get converted to camelCase
func GetTableNameCamelCase(errType TableNameExtractErrType, err error, defaultRet string) string {
	tableName := GetTableName(errType, err, defaultRet)
	if tableName == defaultRet {
		return defaultRet
	}
	return strcase.ToLowerCamel(tableName)
}
