package session

import (
	"testing"
)

type User struct {
	Name string `gingleorm:"PRIMARY KEY"`
	Age  int
}

func TestSession_Model(t *testing.T) {
	s := NewSession().Model(&User{})
	userTable := s.Schema()
	s.Model(&Session{})
	sessionTable := s.Schema()
	if userTable.Name != "User" || sessionTable.Name != "Session" {
		t.Fatal("Failed to change model")
	}
}

func TestSession_CreateTable(t *testing.T) {
	s := NewSession().Model(&User{})
	_ = s.DropTable()
	_ = s.CreateTable()
	if !s.ExistTable() {
		t.Fatal("Failed to create table")
	}
}
