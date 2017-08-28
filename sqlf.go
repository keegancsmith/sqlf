// Package sqlf generates SQL statements in Go, sprintf style.
package sqlf

import (
	"fmt"
	"io"
)

// SQL stores a SQL query and arguments for passing on to
// database/sql/db.Query or gorp.SqlExecutor.
type SQL struct {
	fmt  string
	args []interface{}
}

// Sprintf formats according to a format specifier and returns the resulting
// SQL struct.
func Sprintf(format string, args ...interface{}) *SQL {
	f := make([]interface{}, len(args))
	a := make([]interface{}, 0, len(args))
	for i, arg := range args {
		if sql, ok := arg.(*SQL); ok {
			f[i] = ignoreFormat{sql.fmt}
			a = append(a, sql.args...)
		} else {
			f[i] = ignoreFormat{"%s"}
			a = append(a, arg)
		}
	}
	return &SQL{
		fmt:  fmt.Sprintf(format, f...),
		args: a,
	}
}

// Query returns a string for use in database/sql/db.Query. binder is used to
// update the format specifiers with the relevant BindVar format
func (e *SQL) Query(binder BindVar) string {
	a := make([]interface{}, len(e.args))
	for i := range a {
		a[i] = ignoreFormat{binder.BindVar(i)}
	}
	return fmt.Sprintf(e.fmt, a...)
}

// Args returns the args for use in database/sql/db.Query along with
// SQL.Query()
func (e *SQL) Args() []interface{} {
	return e.args
}

type ignoreFormat struct{ s string }

func (e ignoreFormat) Format(f fmt.State, c rune) {
	io.WriteString(f, e.s)
}
