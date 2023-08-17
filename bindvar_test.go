package sqlf_test

import (
	"testing"

	"github.com/keegancsmith/sqlf"
)

func BenchmarkPostgresBindVar(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = sqlf.PostgresBindVar.BindVar(i)
	}
}

func BenchmarkOracleBindVar(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = sqlf.OracleBindVar.BindVar(i)
	}
}

func BenchmarkSQLServerBindVar(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = sqlf.SQLServerBindVar.BindVar(i)
	}
}

func BenchmarkSimpleBindVar(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = sqlf.SimpleBindVar.BindVar(i)
	}
}
