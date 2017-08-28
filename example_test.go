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

func ExampleJoin() {
	// Our inputs
	min_quantity := 100
	name_filters := []string{"apple", "orange", "coffee"}

	var conds []*sqlf.Query
	for _, filter := range name_filters {
		conds = append(conds, sqlf.Sprintf("name LIKE %s", "%"+filter+"%"))
	}
	sub_query := sqlf.Sprintf("SELECT product_id FROM order_item WHERE quantity > %d", min_quantity)
	q := sqlf.Sprintf("SELECT name FROM product WHERE id IN (%s) AND (%s)", sub_query, sqlf.Join(conds, "OR"))

	fmt.Println(q.Query(sqlf.PostgresBindVar))
	fmt.Println(q.Args())
	// Output: SELECT name FROM product WHERE id IN (SELECT product_id FROM order_item WHERE quantity > $1) AND (name LIKE $2 OR name LIKE $3 OR name LIKE $4)
	// [100 %apple% %orange% %coffee%]
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
