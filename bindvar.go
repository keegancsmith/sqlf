package sqlf

import "fmt"

// SQLBindVar is just BindVar from gorp.Dialect. It is used to take a format
// string and convert it into query string that a gorp.SqlExecutor can use.
type SQLBindVar interface {
	// BindVar binds a variable string to use when forming SQL statements
	// in many dbs it is "?", but Postgres appears to use $1
	//
	// i is a zero based index of the bind variable in this statement
	//
	BindVar(i int) string
}

var (
	// SimpleBindVar is the BindVar format used by SQLite, MySQL, SQLServer
	SimpleBindVar = simpleBindVar{}

	// PostgresBindVar is the BindVar format used by PostgreSQL
	PostgresBindVar = postgresBindVar{}

	// OracleBindVar is the BindVar format used by Oracle Database
	OracleBindVar = oracleBindVar{}
)

type simpleBindVar struct{}

// Returns "?"
func (d simpleBindVar) BindVar(i int) string {
	return "?"
}

type postgresBindVar struct{}

// Returns "$(i+1)"
func (d postgresBindVar) BindVar(i int) string {
	return fmt.Sprintf("$%d", i+1)
}

type oracleBindVar struct{}

// Returns ":(i+1)"
func (d oracleBindVar) BindVar(i int) string {
	return fmt.Sprintf(":%d", i+1)
}
