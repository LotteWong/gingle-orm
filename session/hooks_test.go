package session

import (
	"gingle-orm/log"
	"testing"
)

type Account struct {
	ID       int `gingleorm:"PRIMARY KEY"`
	Password string
}

func (account *Account) BeforeInsert(s *Session) error {
	log.Infoln("Before Insert:", account)
	account.ID += 1000
	return nil
}

func (account *Account) AfterQuery(s *Session) error {
	log.Infoln("After Query:", account)
	account.Password = "******"
	return nil
}

func TestSession_CallMethod(t *testing.T) {
	s := NewSession().Model(&Account{})
	dropErr := s.DropTable()
	createErr := s.CreateTable()

	affected, insertErr := s.Insert(&Account{1, "123456"}, &Account{2, "clwong"})
	account := &Account{}
	firstErr := s.First(account)

	if dropErr != nil || createErr != nil || insertErr != nil || firstErr != nil {
		t.Fatal("Failed to call hooks before insert and after query")
	}
	if affected != 2 || account.ID != 1001 || account.Password != "******" {
		t.Fatal("Failed to call hooks before insert and after query")
	}
}
