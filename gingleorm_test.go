package gingleorm

import (
	"errors"
	"gingle-orm/session"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func OpenDB(t *testing.T) *Engine {
	t.Helper()
	engine, err := NewEngine("sqlite3", "gingleorm.db")
	if err != nil {
		t.Fatal("Failed to open database:", err)
	}
	return engine
}

type User struct {
	Name string `gingleorm:"PRIMARY KEY"`
	Age int
}

func transactionRollback(t *testing.T) {
	e := OpenDB(t)
	defer e.Close()

	s := e.NewSession()
	_ = s.Model(&User{}).DropTable()
	_, err := e.Transaction(func (s *session.Session) (res interface{}, err error) {
		_ = s.Model(&User{}).CreateTable()
		_, err = s.Insert(&User{"Tom", 18})
		return nil, err
	})

	if err == nil || s.ExistTable() {
		t.Fatal("Failed to rollback")
	}
}

func transactionCommit(t *testing.T) {
	e := OpenDB()
	defer e.Close()

	s := e.NewSession()
	_ = s.Model(&User{}).DropTable()
	_, err := e.Transaction(func (s *session.Session) (res interface{}, err error) {
		_ = s.Model(&User{}).CreateTable()
		_, err = s.Insert(&User{"Tom", 18})
		return nil, err
	})

	user := &User{}
	_ = s.First(user)
	if err != nil || user.Name != "Tom" || user.Age != 18 {
		t.Fatal("Failed to commit")
	}
}

func TestEngine_Transaction(t *testing.T) {
	t.Run("rollback", func(t *testing.T) {
		transactionRollback(t)
	})
	t.Run("commint", func(t *testing.T) {
		transactionCommit(t)
	})
}
