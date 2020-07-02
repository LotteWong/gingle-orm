package session

import (
	"database/sql"
	"gingle-orm/dialect"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

var (
	TestDB      *sql.DB
	TestDial, _ = dialect.GetDialect("sqlite3")
)

func TestMain(m *testing.M) {
	// TODO: confused :(
	TestDB, _ = sql.Open("sqlite3", "../gingle.db")
	code := m.Run()
	_ = TestDB.Close()
	os.Exit(code)
}

func NewSession() *Session {
	return New(TestDB, TestDial)
}

func TestSession_Exec(t *testing.T) {
	s := NewSession()

	_, _ = s.Raw("DROP TABLE IF EXISTS USER;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()

	res, _ := s.Raw("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam").Exec()
	if count, err := res.RowsAffected(); err != nil || count != 2 {
		t.Fatal("Failed to insert into database:", err)
	}
}

func TestSession_QueryRow(t *testing.T) {
	s := NewSession()

	_, _ = s.Raw("DROP TABLE IF EXISTS USER;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()

	row := s.Raw("SELECT count(*) FROM User").QueryRow()
	count := 0
	if err := row.Scan(&count); err != nil || count != 0 {
		t.Fatal("Failed to query from database:", err)
	}
}
