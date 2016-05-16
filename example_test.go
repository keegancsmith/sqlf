package sqlf_test

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/keegancsmith/sqlf"
)

func Example() {
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
}

var db *sql.DB

func Example_dbquery() {
	age := 27
	// The next two lines are the only difference from the dabatabase/sql/db.Query example.
	// Original is rows, err := db.Query("SELECT name FROM users WHERE age=?", age)
	s := sqlf.Sprintf("SELECT name FROM users WHERE age=%d", age)
	rows, err := db.Query(s.Query(sqlf.SimpleBindVar), s.Args()...)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s is %d\n", name, age)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
