// Package sqlf generates SQL statements in Go, sprintf style.
package sqlf

import (
	"fmt"
	"io"
)

// SQL stores a SQL query and arguments for passing on to database/sql.Query
// or gorp.SqlExecutor.
type SQL struct {
	fmt  string
	args []interface{}
}

// Sprintf generates a SQL struct the format arguments escaped
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

func (e *SQL) Query(binder SQLBindVar) string {
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
