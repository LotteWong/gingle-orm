package session

import (
	"database/sql"
	"gingle-orm/dialect"
	"gingle-orm/log"
	"gingle-orm/schema"
	"strings"
)

// Session is related to sql execution or sql query
type Session struct {
	db        *sql.DB
	dialect   dialect.Dialect
	schema    *schema.Schema
	sqlClause strings.Builder
	sqlVars   []interface{}
}

// New is to return a Session instance
func New(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{
		db:      db,
		dialect: dialect,
	}
}

// DB is to return a database instance
func (s *Session) DB() *sql.DB {
	return s.db
}

// Raw is to set sql clause and variables
func (s *Session) Raw(sql string, values ...interface{}) *Session {
	s.sqlClause.WriteString(sql)
	s.sqlClause.WriteString(" ")
	s.sqlVars = append(s.sqlVars, values...)
	return s
}

// Clear is to reset sql clause and variables
func (s *Session) Clear() {
	s.sqlClause.Reset()
	s.sqlVars = nil
}

// Exec raw sql clause with sql variables
func (s *Session) Exec() (res sql.Result, err error) {
	defer s.Clear()
	log.Infoln(s.sqlClause.String(), s.sqlVars)

	if res, err = s.DB().Exec(s.sqlClause.String(), s.sqlVars...); err != nil {
		log.Errorln(err)
	}

	return
}

// QueryRow retrieves a record from database
func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	log.Infoln(s.sqlClause, s.sqlVars)

	row := s.DB().QueryRow(s.sqlClause.String(), s.sqlVars...)
	return row
}

// QueryRows retrieves a list of records from database
func (s *Session) QueryRows() *sql.Rows {
	defer s.Clear()
	log.Infoln(s.sqlClause, s.sqlVars)

	rows, _ := s.DB().Query(s.sqlClause.String(), s.sqlVars...)
	return rows
}
