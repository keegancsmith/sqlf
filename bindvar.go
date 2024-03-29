package sqlf

import (
	"strconv"
)

// BindVar is used to take a format string and convert it into query
// string that a gorp.SqlExecutor or sql.Query can use. It is the
// same as BindVar from gorp.Dialect.
type BindVar interface {
	// BindVar binds a variable string to use when forming SQL statements
	// in many dbs it is "?", but Postgres appears to use $1
	//
	// i is a zero based index of the bind variable in this statement
	//
	BindVar(i int) string
}

// SimpleBindVar is the BindVar format used by SQLite, MySQL, SQLServer(mssql driver)
var SimpleBindVar = simpleBindVar{}

// SQLServerBindVar is the BindVar format used by SQL Server(sqlserver driver)
// https://github.com/denisenkom/go-mssqldb#parameters
// https://github.com/denisenkom/go-mssqldb#deprecated
var SQLServerBindVar = sqlServerBindVar{}

// PostgresBindVar is the BindVar format used by PostgreSQL
var PostgresBindVar = postgresBindVar{}

// OracleBindVar is the BindVar format used by Oracle Database
var OracleBindVar = oracleBindVar{}

type simpleBindVar struct{}

// Returns "?"
func (d simpleBindVar) BindVar(i int) string {
	return "?"
}

type postgresBindVar struct{}

// Returns "$(i+1)"
func (d postgresBindVar) BindVar(i int) string {
	return "$" + strconv.Itoa(i+1)
}

type sqlServerBindVar struct{}

// Returns "@p(i+1)"
func (d sqlServerBindVar) BindVar(i int) string {
	return "@p" + strconv.Itoa(i+1)
}

type oracleBindVar struct{}

// Returns ":(i+1)"
func (d oracleBindVar) BindVar(i int) string {
	return ":" + strconv.Itoa(i+1)
}
