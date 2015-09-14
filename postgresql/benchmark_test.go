package postgresql

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/jmoiron/sqlx"
	"upper.io/db"
)

const (
	testRows = 1000
)

func updatedArtistN(i int) string {
	return fmt.Sprintf("Updated Artist %d", i%testRows)
}

func artistN(i int) string {
	return fmt.Sprintf("Artist %d", i%testRows)
}

func connectAndAddFakeRows() (db.Database, error) {
	var err error
	var sess db.Database

	if sess, err = db.Open(Adapter, settings); err != nil {
		return nil, err
	}

	driver := sess.Driver().(*sqlx.DB)

	if _, err = driver.Exec(`TRUNCATE TABLE "artist" RESTART IDENTITY`); err != nil {
		return nil, err
	}

	for i := 0; i < testRows; i++ {
		if _, err = driver.Exec(`INSERT INTO "artist" ("name") VALUES($1)`, artistN(i)); err != nil {
			return nil, err
		}
	}

	return sess, nil
}

// BenchmarkSQLAppend benchmarks raw INSERT SQL queries without using prepared
// statements nor arguments.
func BenchmarkSQLAppend(b *testing.B) {
	var err error
	var sess db.Database

	if sess, err = db.Open(Adapter, settings); err != nil {
		b.Fatal(err)
	}

	defer sess.Close()

	driver := sess.Driver().(*sqlx.DB)

	if _, err = driver.Exec(`TRUNCATE TABLE "artist" RESTART IDENTITY`); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err = driver.Exec(`INSERT INTO "artist" ("name") VALUES('Hayao Miyazaki')`); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkSQLAppendWithArgs benchmarks raw SQL queries with arguments but
// without using prepared statements. The SQL query looks like the one that is
// generated by upper.io/db.
func BenchmarkSQLAppendWithArgs(b *testing.B) {
	var err error
	var sess db.Database

	if sess, err = db.Open(Adapter, settings); err != nil {
		b.Fatal(err)
	}

	defer sess.Close()

	driver := sess.Driver().(*sqlx.DB)

	if _, err = driver.Exec(`TRUNCATE TABLE "artist" RESTART IDENTITY`); err != nil {
		b.Fatal(err)
	}

	args := []interface{}{
		"Hayao Miyazaki",
	}

	var rows *sqlx.Rows

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if rows, err = driver.Queryx(`INSERT INTO "artist" ("name") VALUES($1) RETURNING "id"`, args...); err != nil {
			b.Fatal(err)
		}
		rows.Close()
	}
}

// BenchmarkSQLPreparedAppend benchmarks raw INSERT SQL queries using prepared
// statements but no arguments.
func BenchmarkSQLPreparedAppend(b *testing.B) {
	var err error
	var sess db.Database

	if sess, err = db.Open(Adapter, settings); err != nil {
		b.Fatal(err)
	}

	defer sess.Close()

	driver := sess.Driver().(*sqlx.DB)

	if _, err = driver.Exec(`TRUNCATE TABLE "artist" RESTART IDENTITY`); err != nil {
		b.Fatal(err)
	}

	stmt, err := driver.Prepare(`INSERT INTO "artist" ("name") VALUES('Hayao Miyazaki')`)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err = stmt.Exec(); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkSQLAppendWithArgs benchmarks raw INSERT SQL queries with arguments
// using prepared statements. The SQL query looks like the one that is
// generated by upper.io/db.
func BenchmarkSQLPreparedAppendWithArgs(b *testing.B) {
	var err error
	var sess db.Database

	if sess, err = db.Open(Adapter, settings); err != nil {
		b.Fatal(err)
	}

	defer sess.Close()

	driver := sess.Driver().(*sqlx.DB)

	if _, err = driver.Exec(`TRUNCATE TABLE "artist"`); err != nil {
		b.Fatal(err)
	}

	stmt, err := driver.Preparex(`INSERT INTO "artist" ("name") VALUES($1) RETURNING "id"`)

	if err != nil {
		b.Fatal(err)
	}

	args := []interface{}{
		"Hayao Miyazaki",
	}

	var rows *sqlx.Rows

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if rows, err = stmt.Queryx(args...); err != nil {
			b.Fatal(err)
		}
		rows.Close()
	}
}

// BenchmarkSQLAppendWithVariableArgs benchmarks raw INSERT SQL queries with
// arguments using prepared statements. The SQL query looks like the one that
// is generated by upper.io/db.
func BenchmarkSQLPreparedAppendWithVariableArgs(b *testing.B) {
	var err error
	var sess db.Database

	if sess, err = db.Open(Adapter, settings); err != nil {
		b.Fatal(err)
	}

	defer sess.Close()

	driver := sess.Driver().(*sqlx.DB)

	if _, err = driver.Exec(`TRUNCATE TABLE "artist"`); err != nil {
		b.Fatal(err)
	}

	stmt, err := driver.Preparex(`INSERT INTO "artist" ("name") VALUES($1) RETURNING "id"`)

	if err != nil {
		b.Fatal(err)
	}

	var rows *sqlx.Rows

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		args := []interface{}{
			fmt.Sprintf("Hayao Miyazaki %d", rand.Int()),
		}
		if rows, err = stmt.Queryx(args...); err != nil {
			b.Fatal(err)
		}
		rows.Close()
	}
}

// BenchmarkSQLPreparedAppendTransactionWithArgs benchmarks raw INSERT queries
// within a transaction block with arguments and prepared statements. SQL
// queries look like those generated by upper.io/db.
func BenchmarkSQLPreparedAppendTransactionWithArgs(b *testing.B) {
	var err error
	var sess db.Database
	var tx *sqlx.Tx

	if sess, err = db.Open(Adapter, settings); err != nil {
		b.Fatal(err)
	}

	defer sess.Close()

	driver := sess.Driver().(*sqlx.DB)

	if tx, err = driver.Beginx(); err != nil {
		b.Fatal(err)
	}

	if _, err = tx.Exec(`TRUNCATE TABLE "artist" RESTART IDENTITY`); err != nil {
		b.Fatal(err)
	}

	stmt, err := tx.Preparex(`INSERT INTO "artist" ("name") VALUES($1) RETURNING "id"`)
	if err != nil {
		b.Fatal(err)
	}

	args := []interface{}{
		"Hayao Miyazaki",
	}

	var rows *sqlx.Rows

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if rows, err = stmt.Queryx(args...); err != nil {
			b.Fatal(err)
		}
		rows.Close()
	}

	if err = tx.Commit(); err != nil {
		b.Fatal(err)
	}
}

// BenchmarkUpperAppend benchmarks an insertion by upper.io/db.
func BenchmarkUpperAppend(b *testing.B) {

	sess, err := db.Open(Adapter, settings)
	if err != nil {
		b.Fatal(err)
	}

	defer sess.Close()

	artist, err := sess.Collection("artist")
	if err != nil {
		b.Fatal(err)
	}

	artist.Truncate()

	item := struct {
		Name string `db:"name"`
	}{"Hayao Miyazaki"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err = artist.Append(item); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkUpperAppendVariableArgs benchmarks an insertion by upper.io/db
// with variable parameters.
func BenchmarkUpperAppendVariableArgs(b *testing.B) {

	sess, err := db.Open(Adapter, settings)
	if err != nil {
		b.Fatal(err)
	}

	defer sess.Close()

	artist, err := sess.Collection("artist")
	if err != nil {
		b.Fatal(err)
	}

	artist.Truncate()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		item := struct {
			Name string `db:"name"`
		}{fmt.Sprintf("Hayao Miyazaki %d", rand.Int())}
		if _, err = artist.Append(item); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkUpperAppendTransaction benchmarks insertion queries by upper.io/db
// within a transaction operation.
func BenchmarkUpperAppendTransaction(b *testing.B) {
	var sess db.Database
	var err error

	if sess, err = db.Open(Adapter, settings); err != nil {
		b.Fatal(err)
	}

	defer sess.Close()

	var tx db.Tx
	if tx, err = sess.Transaction(); err != nil {
		b.Fatal(err)
	}
	defer tx.Close()

	var artist db.Collection
	if artist, err = tx.Collection("artist"); err != nil {
		b.Fatal(err)
	}

	if err = artist.Truncate(); err != nil {
		b.Fatal(err)
	}

	item := struct {
		Name string `db:"name"`
	}{"Hayao Miyazaki"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err = artist.Append(item); err != nil {
			b.Fatal(err)
		}
	}

	if err = tx.Commit(); err != nil {
		b.Fatal(err)
	}
}

// BenchmarkUpperAppendTransactionWithMap benchmarks insertion queries by
// upper.io/db within a transaction operation using a map instead of a struct.
func BenchmarkUpperAppendTransactionWithMap(b *testing.B) {
	var sess db.Database
	var err error

	if sess, err = db.Open(Adapter, settings); err != nil {
		b.Fatal(err)
	}

	defer sess.Close()

	var tx db.Tx
	if tx, err = sess.Transaction(); err != nil {
		b.Fatal(err)
	}
	defer tx.Close()

	var artist db.Collection
	if artist, err = tx.Collection("artist"); err != nil {
		b.Fatal(err)
	}

	if err = artist.Truncate(); err != nil {
		b.Fatal(err)
	}

	item := map[string]string{
		"name": "Hayao Miyazaki",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err = artist.Append(item); err != nil {
			b.Fatal(err)
		}
	}

	if err = tx.Commit(); err != nil {
		b.Fatal(err)
	}
}

// BenchmarkSQLSelect benchmarks SQL SELECT queries.
func BenchmarkSQLSelect(b *testing.B) {
	var err error
	var sess db.Database

	if sess, err = connectAndAddFakeRows(); err != nil {
		b.Fatal(err)
	}

	defer sess.Close()

	driver := sess.Driver().(*sqlx.DB)

	var res *sqlx.Rows

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if res, err = driver.Queryx(`SELECT * FROM "artist" WHERE "name" = $1`, artistN(i)); err != nil {
			b.Fatal(err)
		}
		res.Close()
	}
}

// BenchmarkSQLPreparedSelect benchmarks SQL select queries using prepared
// statements.
func BenchmarkSQLPreparedSelect(b *testing.B) {
	var err error
	var sess db.Database

	if sess, err = connectAndAddFakeRows(); err != nil {
		b.Fatal(err)
	}
	defer sess.Close()

	driver := sess.Driver().(*sqlx.DB)

	stmt, err := driver.Preparex(`SELECT * FROM "artist" WHERE "name" = $1`)
	if err != nil {
		b.Fatal(err)
	}

	var res *sqlx.Rows

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if res, err = stmt.Queryx(artistN(i)); err != nil {
			b.Fatal(err)
		}
		res.Close()
	}
}

// BenchmarkUpperFind benchmarks upper.io/db's One method.
func BenchmarkUpperFind(b *testing.B) {
	var err error
	var sess db.Database

	if sess, err = connectAndAddFakeRows(); err != nil {
		b.Fatal(err)
	}

	defer sess.Close()

	artist, err := sess.Collection("artist")
	if err != nil {
		b.Fatal(err)
	}

	type artistType struct {
		Name string `db:"name"`
	}

	var item artistType

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res := artist.Find(db.Cond{"name": artistN(i)})
		if err = res.One(&item); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkUpperFindWithC benchmarks upper.io/db's One method.
func BenchmarkUpperFindWithC(b *testing.B) {
	var err error
	var sess db.Database

	if sess, err = connectAndAddFakeRows(); err != nil {
		b.Fatal(err)
	}

	defer sess.Close()

	type artistType struct {
		Name string `db:"name"`
	}

	var item artistType

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res := sess.C("artist").Find(db.Cond{"name": artistN(i)})
		if err = res.One(&item); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkUpperFindAll benchmarks upper.io/db's All method.
func BenchmarkUpperFindAll(b *testing.B) {
	var err error
	var sess db.Database

	if sess, err = connectAndAddFakeRows(); err != nil {
		b.Fatal(err)
	}

	defer sess.Close()

	artist, err := sess.Collection("artist")
	if err != nil {
		b.Fatal(err)
	}

	type artistType struct {
		Name string `db:"name"`
	}

	var items []artistType

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res := artist.Find(db.Or{
			db.Cond{"name": artistN(i)},
			db.Cond{"name": artistN(i + 1)},
			db.Cond{"name": artistN(i + 2)},
		})
		if err = res.All(&items); err != nil {
			b.Fatal(err)
		}
		if len(items) != 3 {
			b.Fatal("Expecting 3 results.")
		}
	}
}

// BenchmarkSQLUpdate benchmarks SQL UPDATE queries.
func BenchmarkSQLUpdate(b *testing.B) {
	var err error
	var sess db.Database

	if sess, err = connectAndAddFakeRows(); err != nil {
		b.Fatal(err)
	}

	defer sess.Close()

	driver := sess.Driver().(*sqlx.DB)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err = driver.Exec(`UPDATE "artist" SET "name" = $1 WHERE "name" = $2`, updatedArtistN(i), artistN(i)); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkSQLPreparedUpdate benchmarks SQL UPDATE queries.
func BenchmarkSQLPreparedUpdate(b *testing.B) {
	var err error
	var sess db.Database

	if sess, err = connectAndAddFakeRows(); err != nil {
		b.Fatal(err)
	}

	defer sess.Close()

	driver := sess.Driver().(*sqlx.DB)

	stmt, err := driver.Prepare(`UPDATE "artist" SET "name" = $1 WHERE "name" = $2`)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err = stmt.Exec(updatedArtistN(i), artistN(i)); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkUpperUpdate benchmarks upper.io/db's Update method.
func BenchmarkUpperUpdate(b *testing.B) {
	var err error
	var sess db.Database

	if sess, err = connectAndAddFakeRows(); err != nil {
		b.Fatal(err)
	}

	defer sess.Close()

	artist, err := sess.Collection("artist")
	if err != nil {
		b.Fatal(err)
	}

	type artistType struct {
		Name string `db:"name"`
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		newValue := artistType{
			Name: updatedArtistN(i),
		}
		res := artist.Find(db.Cond{"name": artistN(i)})
		if err = res.Update(newValue); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkSQLDelete benchmarks SQL DELETE queries.
func BenchmarkSQLDelete(b *testing.B) {
	var err error
	var sess db.Database

	if sess, err = connectAndAddFakeRows(); err != nil {
		b.Fatal(err)
	}

	defer sess.Close()

	driver := sess.Driver().(*sqlx.DB)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err = driver.Exec(`DELETE FROM "artist" WHERE "name" = $1`, artistN(i)); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkSQLPreparedDelete benchmarks SQL DELETE queries.
func BenchmarkSQLPreparedDelete(b *testing.B) {
	var err error
	var sess db.Database

	if sess, err = connectAndAddFakeRows(); err != nil {
		b.Fatal(err)
	}
	defer sess.Close()

	driver := sess.Driver().(*sqlx.DB)

	stmt, err := driver.Prepare(`DELETE FROM "artist" WHERE "name" = $1`)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err = stmt.Exec(artistN(i)); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkUpperRemove benchmarks
func BenchmarkUpperRemove(b *testing.B) {
	var err error
	var sess db.Database

	if sess, err = connectAndAddFakeRows(); err != nil {
		b.Fatal(err)
	}

	defer sess.Close()

	artist, err := sess.Collection("artist")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res := artist.Find(db.Cond{"name": artistN(i)})
		if err = res.Remove(); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkUpperGetCollection
func BenchmarkUpperGetCollection(b *testing.B) {
	var err error
	var sess db.Database

	if sess, err = db.Open(Adapter, settings); err != nil {
		b.Fatal(err)
	}

	defer sess.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := sess.Collection("artist")
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkUpperC
func BenchmarkUpperC(b *testing.B) {
	var err error
	var sess db.Database

	if sess, err = db.Open(Adapter, settings); err != nil {
		b.Fatal(err)
	}

	defer sess.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sess.C("artist")
	}
}

// BenchmarkUpperCommitManyTransactions benchmarks
func BenchmarkUpperCommitManyTransactions(b *testing.B) {
	var sess db.Database
	var err error

	if sess, err = db.Open(Adapter, settings); err != nil {
		b.Fatal(err)
	}

	defer sess.Close()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var tx db.Tx
		if tx, err = sess.Transaction(); err != nil {
			b.Fatal(err)
		}

		var artist db.Collection
		if artist, err = tx.Collection("artist"); err != nil {
			b.Fatal(err)
		}

		if err = artist.Truncate(); err != nil {
			b.Fatal(err)
		}

		item := struct {
			Name string `db:"name"`
		}{"Hayao Miyazaki"}

		if _, err = artist.Append(item); err != nil {
			b.Fatal(err)
		}

		if err = tx.Commit(); err != nil {
			b.Fatal(err)
		}

		tx.Close()
	}
}

// BenchmarkUpperRollbackManyTransactions benchmarks
func BenchmarkUpperRollbackManyTransactions(b *testing.B) {
	var sess db.Database
	var err error

	if sess, err = db.Open(Adapter, settings); err != nil {
		b.Fatal(err)
	}

	defer sess.Close()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var tx db.Tx
		if tx, err = sess.Transaction(); err != nil {
			b.Fatal(err)
		}

		var artist db.Collection
		if artist, err = tx.Collection("artist"); err != nil {
			b.Fatal(err)
		}

		if err = artist.Truncate(); err != nil {
			b.Fatal(err)
		}

		item := struct {
			Name string `db:"name"`
		}{"Hayao Miyazaki"}

		if _, err = artist.Append(item); err != nil {
			b.Fatal(err)
		}

		if err = tx.Rollback(); err != nil {
			b.Fatal(err)
		}

		tx.Close()
	}
}
