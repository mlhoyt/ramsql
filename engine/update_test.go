package engine_test

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/mlhoyt/ramsql/engine/log"

	_ "github.com/mlhoyt/ramsql/driver"
)

func TestUpdateSimple(t *testing.T) {
	log.UseTestLogger(t)

	db, err := sql.Open("ramsql", "TestUpdateSimple")
	if err != nil {
		t.Fatalf("sql.Open : Error : %s\n", err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE account (id INT AUTOINCREMENT, email TEXT)")
	if err != nil {
		t.Fatalf("sql.Exec: Error: %s\n", err)
	}

	_, err = db.Exec("INSERT INTO account ('email') VALUES ('foo@bar.com')")
	if err != nil {
		t.Fatalf("Cannot insert into table account: %s", err)
	}

	_, err = db.Exec("INSERT INTO account ('email') VALUES ('leon@bar.com')")
	if err != nil {
		t.Fatalf("Cannot insert into table account: %s", err)
	}

	_, err = db.Exec("UPDATE account SET email = 'roger@gmail.com' WHERE id = 2")
	if err != nil {
		t.Fatalf("Cannot update table account: %s", err)
	}

	row := db.QueryRow("SELECT * FROM account WHERE id = 2")
	if row == nil {
		t.Fatalf("sql.Query failed")
	}

	var email string
	var id int
	err = row.Scan(&id, &email)
	if err != nil {
		t.Fatalf("row.Scan: %s", err)
	}

	if email != "roger@gmail.com" {
		t.Fatalf("Expected email 'roger@gmail.com', got '%s'", email)
	}
}

func TestUpdateSimpleOtherAutoIncrement(t *testing.T) {
	log.UseTestLogger(t)

	db, err := sql.Open("ramsql", "TestUpdateSimpleOtherAutoIncrement")
	if err != nil {
		t.Fatalf("sql.Open : Error : %s\n", err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE account (id INT AUTO_INCREMENT, email TEXT)")
	if err != nil {
		t.Fatalf("sql.Exec: Error: %s\n", err)
	}

	_, err = db.Exec("INSERT INTO account ('email') VALUES ('foo@bar.com')")
	if err != nil {
		t.Fatalf("Cannot insert into table account: %s", err)
	}

	_, err = db.Exec("INSERT INTO account ('email') VALUES ('leon@bar.com')")
	if err != nil {
		t.Fatalf("Cannot insert into table account: %s", err)
	}

	_, err = db.Exec("UPDATE account SET email = 'roger@gmail.com' WHERE id = 2")
	if err != nil {
		t.Fatalf("Cannot update table account: %s", err)
	}

	row := db.QueryRow("SELECT * FROM account WHERE id = 2")
	if row == nil {
		t.Fatalf("sql.Query failed")
	}

	var email string
	var id int
	err = row.Scan(&id, &email)
	if err != nil {
		t.Fatalf("row.Scan: %s", err)
	}

	if email != "roger@gmail.com" {
		t.Fatalf("Expected email 'roger@gmail.com', got '%s'", email)
	}
}

func TestUpdateIsNull(t *testing.T) {
	log.UseTestLogger(t)

	db, err := sql.Open("ramsql", "TestUpdateIsNull")
	if err != nil {
		t.Fatalf("sql.Open : Error : %s\n", err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE account (id INT AUTOINCREMENT, email TEXT, creation_date TIMESTAMP WITH TIME ZONE)")
	if err != nil {
		t.Fatalf("sql.Exec: Error: %s\n", err)
	}

	_, err = db.Exec("INSERT INTO account ('email') VALUES ('foo@bar.com')")
	if err != nil {
		t.Fatalf("Cannot insert into table account: %s", err)
	}

	_, err = db.Exec("INSERT INTO account ('email') VALUES ('leon@bar.com')")
	if err != nil {
		t.Fatalf("Cannot insert into table account: %s", err)
	}

	res, err := db.Exec("UPDATE account SET email = 'roger@gmail.com', creation_date = $1 WHERE id = 2 AND creation_date IS NULL", time.Now())
	if err != nil {
		t.Fatalf("Cannot update table account: %s", err)
	}

	ra, err := res.RowsAffected()
	if err != nil {
		t.Fatalf("Cannot check number of rows affected: %s", err)
	}
	if ra != 1 {
		t.Fatalf("Expected 1 row, affected. Got %d", ra)
	}

	rows, err := db.Query(`SELECT id FROM account WHERE creation_date IS NULL`)
	if err != nil {
		t.Fatalf("cannot select null columns: %s", err)
	}

	var n, id int64
	for rows.Next() {
		n++
		err = rows.Scan(&id)
		if err != nil {
			t.Fatalf("cannot scan null columns: %s", err)
		}
	}
	rows.Close()
	if n != 1 {
		t.Fatalf("Expected 1 rows, got %d", n)
	}

	rows, err = db.Query(`SELECT id FROM account WHERE creation_date IS NOT NULL`)
	if err != nil {
		t.Fatalf("cannot select not null columns: %s", err)
	}

	n = 0
	for rows.Next() {
		n++
		err = rows.Scan(&id)
		if err != nil {
			t.Fatalf("cannot scan null columns: %s", err)
		}
	}
	rows.Close()
	if n != 1 {
		t.Fatalf("Expected 1 rows, got %d", n)
	}

}

func TestUpdateNotNull(t *testing.T) {
	log.UseTestLogger(t)

	db, err := sql.Open("ramsql", "TestUpdateNotNull")
	if err != nil {
		t.Fatalf("sql.Open : Error : %s\n", err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE account (id INT AUTOINCREMENT, email TEXT, creation_date TIMESTAMP WITH TIME ZONE)")
	if err != nil {
		t.Fatalf("sql.Exec: Error: %s\n", err)
	}

	_, err = db.Exec("INSERT INTO account ('email') VALUES ('foo@bar.com')")
	if err != nil {
		t.Fatalf("Cannot insert into table account: %s", err)
	}

	_, err = db.Exec("INSERT INTO account ('email') VALUES ('leon@bar.com')")
	if err != nil {
		t.Fatalf("Cannot insert into table account: %s", err)
	}

	_, err = db.Exec("UPDATE account SET email = 'roger@gmail.com' WHERE id = 2 AND creation_date IS NOT NULL")
	if err != nil {
		t.Fatalf("Cannot update table account: %s", err)
	}

}

func TestInsertDefaultCurrentTimestamp(t *testing.T) {
	log.UseTestLogger(t)

	db, err := sql.Open("ramsql", "TestInsertDefaultCurrentTimestamp")
	if err != nil {
		t.Fatalf("sql.Open : Error : %s\n", err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE account (id INT AUTOINCREMENT, email TEXT, creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP)")
	if err != nil {
		t.Fatalf("sql.Exec: Error: %s\n", err)
	}

	_, err = db.Exec("INSERT INTO account ('email') VALUES ('foo@bar.com')")
	if err != nil {
		t.Fatalf("Cannot insert into table account: %s", err)
	}

	rows, err := db.Query(`SELECT id, creation_date FROM account WHERE email = 'foo@bar.com'`)
	if err != nil {
		t.Fatalf("cannot select row: %s", err)
	}

	var id int
	var creationDate string

	n := 0
	for rows.Next() {
		n++

		err = rows.Scan(&id, &creationDate)
		if err != nil {
			t.Fatalf("cannot scan row %d: %s", n, err)
		}

		if creationDate == "CURRENT_TIMESTAMP" {
			t.Fatalf("Expected timestamp value for creation_date but found string literal 'CURRENT_TIMESTAMP'")
		}
	}
	rows.Close()

	if n != 1 {
		t.Fatalf("Expected 1 rows, got %d", n)
	}
}

func TestUpdateDefaultCurrentTimestampOnUpdateCurrentTimestamp(t *testing.T) {
	log.UseTestLogger(t)

	db, err := sql.Open("ramsql", "TestUpdateDefaultCurrentTimestampOnUpdateCurrentTimestamp")
	if err != nil {
		t.Fatalf("sql.Open : Error : %s\n", err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE account (id INT AUTOINCREMENT, email TEXT, modified_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP)")
	if err != nil {
		t.Fatalf("sql.Exec: Error: %s\n", err)
	}

	_, err = db.Exec("INSERT INTO account ('email') VALUES ('foo@bar.com')")
	if err != nil {
		t.Fatalf("Cannot insert into table account: %s", err)
	}

	var id int
	var modifiedDate string

	rows, err := db.Query(`SELECT id, modified_date FROM account WHERE email = 'foo@bar.com'`)
	if err != nil {
		t.Fatalf("cannot select row: %s", err)
	}

	n := 0
	for rows.Next() {
		n++
		err = rows.Scan(&id, &modifiedDate)
		if err != nil {
			t.Fatalf("cannot scan row %d: %s", n, err)
		}
	}
	rows.Close()

	if n != 1 {
		t.Fatalf("Expected 1 rows, got %d", n)
	}

	insertModifiedDate := modifiedDate

	_, err = db.Exec(fmt.Sprintf("UPDATE account SET email = 'roger@gmail.com' WHERE id = %d", id))
	if err != nil {
		t.Fatalf("Cannot update table account: %s", err)
	}

	rows, err = db.Query(fmt.Sprintf("SELECT id, modified_date FROM account WHERE email = 'roger@gmail.com'"))
	if err != nil {
		t.Fatalf("cannot select row: %s", err)
	}

	n = 0
	for rows.Next() {
		n++
		err = rows.Scan(&id, &modifiedDate)
		if err != nil {
			t.Fatalf("cannot scan row %d: %s", n, err)
		}
	}
	rows.Close()

	if n != 1 {
		t.Fatalf("Expected 1 rows, got %d", n)
	}

	updateModifiedDate := modifiedDate

	if insertModifiedDate == updateModifiedDate {
		t.Fatalf("Expected insert modified_date (%s) and update modified_date (%s) to be different", insertModifiedDate, updateModifiedDate)
	}
}
