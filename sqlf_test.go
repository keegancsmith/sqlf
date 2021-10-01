package sqlf_test

import (
	"reflect"
	"testing"

	"github.com/keegancsmith/sqlf"
)

func TestSprintf(t *testing.T) {
	cases := map[string]struct {
		Fmt      string
		FmtArgs  []interface{}
		Want     string
		WantArgs []interface{}
	}{
		"simple_substitute": {
			"SELECT * FROM test_table WHERE a = %s AND b = %d", []interface{}{"foo", 1},
			"SELECT * FROM test_table WHERE a = $1 AND b = $2", []interface{}{"foo", 1},
		},

		"simple_embedded": {
			"SELECT * FROM test_table WHERE a = (%s)", []interface{}{sqlf.Sprintf("SELECT b FROM b_table WHERE x = %d", 1)},
			"SELECT * FROM test_table WHERE a = (SELECT b FROM b_table WHERE x = $1)", []interface{}{1},
		},

		"embedded": {
			"SELECT * FROM test_table WHERE a = %s AND c = (%s) AND d = %s", []interface{}{"foo", sqlf.Sprintf("SELECT b FROM b_table WHERE x = %d", 1), "bar"},
			"SELECT * FROM test_table WHERE a = $1 AND c = (SELECT b FROM b_table WHERE x = $2) AND d = $3", []interface{}{"foo", 1, "bar"},
		},

		"embedded_embedded": {
			"SELECT * FROM test_table WHERE a = %s AND c = (%s) AND d = %s",
			[]interface{}{"foo", sqlf.Sprintf("SELECT b FROM b_table WHERE x = %d AND y = (%s)", 1, sqlf.Sprintf("SELECT %s", "baz")), "bar"},
			"SELECT * FROM test_table WHERE a = $1 AND c = (SELECT b FROM b_table WHERE x = $2 AND y = (SELECT $3)) AND d = $4",
			[]interface{}{"foo", 1, "baz", "bar"},
		},

		"literal_percent_operator": {
			"SELECT * FROM test_table WHERE a <<%% %s AND b = %d", []interface{}{"foo", 1},
			"SELECT * FROM test_table WHERE a <<% $1 AND b = $2", []interface{}{"foo", 1},
		},
	}
	for tn, tc := range cases {
		q := sqlf.Sprintf(tc.Fmt, tc.FmtArgs...)
		if got := q.Query(sqlf.PostgresBindVar); got != tc.Want {
			t.Errorf("%s: expected query: %q, got: %q", tn, tc.Want, got)
		}
		if got := q.Args(); !reflect.DeepEqual(got, tc.WantArgs) {
			t.Errorf("%s: expected args: %q, got: %q", tn, tc.WantArgs, got)
		}
	}
}
