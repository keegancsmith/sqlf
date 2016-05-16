sqlf [![Build Status](https://travis-ci.org/keegancsmith/sqlf.svg?branch=master)](https://travis-ci.org/) [![GoDoc](https://godoc.org/github.com/keegancsmith/sqlf?status.svg)](https://godoc.org/github.com/keegancsmith/sqlf)
======

Generate SQL Commands in Go, sprintf Style.

Quick example:

```go
// This is an example which shows off embedding SQL, which simplifies building
// complicated SQL queries
name := "John"
age, offset := 27, 100
where := sqlf.Sprintf("name=%s AND age=%d", name, age)
limit := sqlf.Sprintf("%d OFFSET %d", 10, offset)
q := sqlf.Sprintf("SELECT name FROM users WHERE %s LIMIT %s", where, limit)
fmt.Println(q.Query(sqlf.PostgresBindVar))
fmt.Println(q.Args())
// Output: SELECT name FROM users WHERE name=$1 AND age=$2 LIMIT $3 OFFSET $4
// [John 27 10 100]
```

Notice that you can pass the output of `sqlf.Sprintf` as input to itself. It
will return a flattened query string, while preserving the correct variable
binding.

See https://godoc.org/github.com/keegancsmith/sqlf for more information.
