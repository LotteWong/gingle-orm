package session

import (
	"fmt"
	"gingle-orm/log"
	"gingle-orm/schema"
	"reflect"
	"strings"
)

// Model is to update schema
func (s *Session) Model(model interface{}) *Session {
	if s.schema == nil || reflect.TypeOf(model) != reflect.TypeOf(s.schema.Model) {
		s.schema = schema.Parse(model, s.dialect)
	}

	return s
}

// Schema is to retrieve schema
func (s *Session) Schema() *schema.Schema {
	if s.schema == nil {
		log.Errorln("Model is not set")
	}
	return s.schema
}

// CreateTable is to create a database table from schema
func (s *Session) CreateTable() (err error) {
	table := s.Schema()
	columns := []string{}
	for _, field := range table.Fields {
		columns = append(columns, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag))
	}
	desc := strings.Join(columns, ",")
	_, err = s.Raw(fmt.Sprintf("CREATE TABLE %s (%s)", table.Name, desc)).Exec()
	if err != nil {
		log.Errorln(err)
		return
	}
	return
}

// DropTable is to drop a database table from schema
func (s *Session) DropTable() (err error) {
	table := s.Schema()
	_, err = s.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %s", table.Name)).Exec()
	if err != nil {
		log.Errorln(err)
		return
	}
	return
}

// ExistTable is to check a database table from schema
func (s *Session) ExistTable() (exist bool) {
	table := s.Schema()
	sqlClause, sqlVars := s.dialect.TableExistSQL(table.Name)
	row := s.Raw(sqlClause, sqlVars...).QueryRow()
	var tableName string
	row.Scan(&tableName)
	return tableName == table.Name
}
