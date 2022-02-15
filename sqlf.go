// Package sqlf generates parameterized SQL statements in Go, sprintf style.
//
// A simple example:
//
//   q := sqlf.Sprintf("SELECT * FROM users WHERE country = %s AND age > %d", "US", 27);
//   rows, err := db.Query(q.Query(sqlf.SimpleBindVar), q.Args()...) // db is a database/sql.DB
//
// sqlf.Sprintf does not return a string. It returns *sqlf.Query which has
// methods for a parameterized SQL query and arguments. You then pass that to
// db.Query, db.Exec, etc. This is not like using fmt.Sprintf, which could
// expose you to malformed SQL or SQL injection attacks.
//
// sqlf.Query can be passed as an argument to sqlf.Sprintf. It will "flatten"
// the query string, while preserving the correct variable binding. This
// allows you to easily compose and build SQL queries. See the below examples
// to find out more.
package sqlf

import (
	"fmt"
	"io"
	"strings"

	"github.com/lib/pq"
)

// Query stores a SQL expression and arguments for passing on to
// database/sql/db.Query or gorp.SqlExecutor.
type Query struct {
	fmt  string
	args []interface{}
}

// Sprintf formats according to a format specifier and returns the resulting
// Query.
func Sprintf(format string, args ...interface{}) *Query {
	f := make([]interface{}, len(args))
	a := make([]interface{}, 0, len(args))
	for i, arg := range args {
		if q, ok := arg.(*Query); ok {
			f[i] = ignoreFormat{q.fmt}
			a = append(a, q.args...)
		} else {
			f[i] = ignoreFormat{"%s"}
			a = append(a, arg)
		}
	}
	// The format string below goes through fmt.Sprintf, which would reduce `%%` (a literal `%`
	// according to fmt format specifier semantics), but it would also go through fmt.Sprintf
	// again at Query.Query(binder) time - so we need to make sure `%%` remains as `%%` in our
	// format string. See the literal_percent_operator test.
	format = strings.Replace(format, "%%", "%%%%", -1)
	return &Query{
		fmt:  fmt.Sprintf(format, f...),
		args: a,
	}
}

// Query returns a string for use in database/sql/db.Query. binder is used to
// update the format specifiers with the relevant BindVar format
func (q *Query) Query(binder BindVar) string {
	a := make([]interface{}, len(q.args))
	for i := range a {
		a[i] = ignoreFormat{binder.BindVar(i)}
	}
	return fmt.Sprintf(q.fmt, a...)
}

// Args returns the args for use in database/sql/db.Query along with
// q.Query()
func (q *Query) Args() []interface{} {
	return q.args
}

// Render returns a string that can be executed directly as a SQL query.
func (q *Query) Rneder(binder BindVar) (string, error) {
	s := q.Query(binder)
	for i, arg := range q.Args() {
		argString, err := argToString(arg)
		if err != nil {
			return "", err
		}

		// TODO there's a bug where if the arg contains valid binders,
		// they too will be unintentionally replaced by subsequent arguments.

		// TODO support other BindVars such as `?`.

		// TODO there must be a better way to do this, perhaps using `q.fmt`.
		s = strings.ReplaceAll(s, fmt.Sprintf("$%d", i+1), argString)
	}
	return s, nil
}

func argToString(arg interface{}) (string, error) {
	switch arg := arg.(type) {
	case string:
		return fmt.Sprintf("'%s'", escapeQuotes(arg)), nil
	case *pq.StringArray:
		value, err := arg.Value()
		if err != nil {
			return "", err
		}
		switch value := value.(type) {
		case string:
			return fmt.Sprintf("'%s'", escapeQuotes(value)), nil
		default:
			return "", fmt.Errorf("unrecognized array element type %T", value)
		}
	case int:
		return fmt.Sprintf("%d", arg), nil
	default:
		return "", fmt.Errorf("unrecognized type %T", arg)
	}
}

func escapeQuotes(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}

// Join concatenates the elements of queries to create a single Query. The
// separator string sep is placed between elements in the resulting Query.
//
// This is commonly used to join clauses in a WHERE query. As such sep is
// usually "AND" or "OR".
func Join(queries []*Query, sep string) *Query {
	f := make([]string, 0, len(queries))
	var a []interface{}
	for _, q := range queries {
		f = append(f, q.fmt)
		a = append(a, q.args...)
	}
	return &Query{
		fmt:  strings.Join(f, " "+sep+" "),
		args: a,
	}
}

type ignoreFormat struct{ s string }

func (e ignoreFormat) Format(f fmt.State, c rune) {
	io.WriteString(f, e.s)
}
