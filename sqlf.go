package sqlf

import (
	"fmt"
	"io"
)

// SQLBinder is just BindVar from gorp.Dialect. It is used to take a format
// string and convert it into query string that a gorp.SqlExecutor can use.
type SQLBinder interface {
	// BindVar binds a variable string to use when forming SQL statements
	// in many dbs it is "?", but Postgres appears to use $1
	//
	// i is a zero based index of the bind variable in this statement
	//
	BindVar(i int) string
}

type SQL struct {
	fmt  string
	args []interface{}
}

// Sprintf generates a SQL struct the format arguments escaped
func Sprintf(format string, a ...interface{}) *SQL {
	return &SQL{format, a}
}

func (e *SQL) Query(binder SQLBinder) string {
	a := make([]interface{}, len(e.args))
	for i := range a {
		a[i] = ignoreFormat{binder.BindVar(i)}
	}
	return fmt.Sprintf(e.fmt, a...)
}

func (e *SQL) Args() []interface{} {
	return e.args
}

type ignoreFormat struct{ s string }

func (e ignoreFormat) Format(f fmt.State, c rune) {
	io.WriteString(f, e.s)
}
