package db

import (
	"errors"
	"github.com/VividCortex/mysqlerr"
	"testing"

	"github.com/go-sql-driver/mysql"

	"github.com/stretchr/testify/assert"
)

var (
	refErrUser             = errors.New("SQL Error [1452] [23000]: Cannot add or update a child row: a foreign key constraint fails (`mydb`.`item`, CONSTRAINT `fk_item_created_by` FOREIGN KEY (`created_by`) REFERENCES `user` (`id`))")
	refErrEmployeeFunction = errors.New("SQL Error [1452] [23000]: Cannot add or update a child row: a foreign key constraint fails (`mydb`.`employee_access_type_general`, CONSTRAINT `fk_employee_access_type_general_employee_function_id` FOREIGN KEY (`employee_function_id`) REFERENCES `employee_function` (`id`))")
)

func TestCollidesWithUniqueIndex(t *testing.T) {
	assert.False(t, CollidesWithUniqueIndex(nil))
	assert.False(t, CollidesWithUniqueIndex(errors.New("boom")))
	assert.False(t, CollidesWithUniqueIndex(&mysql.MySQLError{}))
	assert.False(t, CollidesWithUniqueIndex(&mysql.MySQLError{Number: 1}))
	assert.True(t, CollidesWithUniqueIndex(&mysql.MySQLError{Number: mysqlerr.ER_DUP_ENTRY}))
}

func TestDeadLockFound(t *testing.T) {
	assert.False(t, DeadLockFound(nil))
	assert.False(t, DeadLockFound(errors.New("boom")))
	assert.False(t, DeadLockFound(&mysql.MySQLError{}))
	assert.False(t, DeadLockFound(&mysql.MySQLError{Number: 1}))
	assert.True(t, DeadLockFound(&mysql.MySQLError{Number: mysqlerr.ER_LOCK_DEADLOCK}))
}

func TestCantAggregate2Collations(t *testing.T) {
	assert.False(t, CantAggregate2Collations(nil))
	assert.False(t, CantAggregate2Collations(errors.New("boom")))
	assert.False(t, CantAggregate2Collations(&mysql.MySQLError{}))
	assert.False(t, CantAggregate2Collations(&mysql.MySQLError{Number: 1}))
	assert.True(t, CantAggregate2Collations(&mysql.MySQLError{Number: mysqlerr.ER_CANT_AGGREGATE_2COLLATIONS}))
}

func TestQueryTimeout(t *testing.T) {
	assert.False(t, QueryTimeout(nil))
	assert.False(t, QueryTimeout(errors.New("boom")))
	assert.False(t, QueryTimeout(&mysql.MySQLError{}))
	assert.False(t, QueryTimeout(&mysql.MySQLError{Number: 1}))
	assert.True(t, QueryTimeout(&mysql.MySQLError{Number: mysqlerr.ER_QUERY_TIMEOUT}))
}

func TestSelectTooBig(t *testing.T) {
	assert.False(t, SelectTooBig(nil))
	assert.False(t, SelectTooBig(errors.New("boom")))
	assert.False(t, SelectTooBig(&mysql.MySQLError{}))
	assert.False(t, SelectTooBig(&mysql.MySQLError{Number: 1}))
	assert.True(t, SelectTooBig(&mysql.MySQLError{Number: mysqlerr.ER_TOO_BIG_SELECT}))
}

func TestReferenceDoesNotExist(t *testing.T) {
	assert.False(t, ReferenceDoesNotExist(nil))
	assert.False(t, ReferenceDoesNotExist(errors.New("boom")))
	assert.False(t, ReferenceDoesNotExist(&mysql.MySQLError{}))
	assert.False(t, ReferenceDoesNotExist(&mysql.MySQLError{Number: 1}))
	assert.True(t, ReferenceDoesNotExist(&mysql.MySQLError{Number: mysqlerr.ER_NO_REFERENCED_ROW_2}))
}

func TestGetTableNameFromErr(t *testing.T) {
	t.Run("with TABLE_NAME_EXTRACT_FROM_ERR_REF", func(t *testing.T) {
		t.Run("returns table name with valid errors", func(t *testing.T) {
			var tableName string
			tableName = GetTableName(TABLE_NAME_EXTRACT_FROM_ERR_REF, refErrUser, "")
			assert.Equal(t, "user", tableName)

			tableName = GetTableName(TABLE_NAME_EXTRACT_FROM_ERR_REF, refErrEmployeeFunction, "")
			assert.Equal(t, "employee_function", tableName)
		})

		t.Run("returns defaultRet on error content not sql", func(t *testing.T) {
			tableName := GetTableName(TABLE_NAME_EXTRACT_FROM_ERR_REF, errors.New("some unknown sql error"), "unknown")
			assert.Equal(t, "unknown", tableName)
		})
	})
	t.Run("with invalid TableNameExtractErrType", func(t *testing.T) {
		invalidTableNameExtractErrType := TableNameExtractErrType(123456)
		tableName := GetTableName(invalidTableNameExtractErrType, errors.New("sql error"), "unknown")
		assert.Equal(t, "unknown", tableName)
	})
}

func TestGetTableNameFromErrCamelCase(t *testing.T) {
	t.Run("with valid TableNameExtractErrType", func(t *testing.T) {
		var tableName string
		tableName = GetTableNameCamelCase(TABLE_NAME_EXTRACT_FROM_ERR_REF, refErrUser, "")
		assert.Equal(t, "user", tableName)

		tableName = GetTableNameCamelCase(TABLE_NAME_EXTRACT_FROM_ERR_REF, refErrEmployeeFunction, "")
		assert.Equal(t, "employeeFunction", tableName)
	})
	t.Run("returns defaultRet without parsing whe GetTableName returns defaultRet", func(t *testing.T) {
		invalidTableNameExtractErrType := TableNameExtractErrType(123456)
		tableName := GetTableNameCamelCase(invalidTableNameExtractErrType, refErrUser, "my unknown entity")
		assert.Equal(t, "my unknown entity", tableName)
	})
}
