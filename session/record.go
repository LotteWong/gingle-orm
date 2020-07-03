package session

import (
	"errors"
	"gingle-orm/clause"
	"gingle-orm/log"
	"reflect"
)

// Insert is to insert a record into database
func (s *Session) Insert(values ...interface{}) (int64, error) {
	recordValues := make([]interface{}, 0)

	// TODO: confused :(
	for _, value := range values {
		s.CallMethod(BeforeInsert, value)
		table := s.Model(value).Schema()
		s.clause.Set(clause.INSERT, table.Name, table.FieldNames)
		recordValues = append(recordValues, table.RecordValues(value))
	}
	s.clause.Set(clause.VALUES, recordValues...)

	sqlClause, sqlVars := s.clause.Build(clause.INSERT, clause.VALUES)

	res, err := s.Raw(sqlClause, sqlVars...).Exec()
	if err != nil {
		log.Errorln(err)
		return 0, err
	}
	s.CallMethod(AfterInsert, nil)
	return res.RowsAffected()
}

// Find is to select records from database
func (s *Session) Find(values interface{}) error {
	s.CallMethod(BeforeQuery, nil)
	modSlice := reflect.Indirect(reflect.ValueOf(values))
	modType := modSlice.Type().Elem()
	table := s.Model(reflect.New(modType).Elem().Interface()).Schema()

	s.clause.Set(clause.SELECT, table.Name, table.FieldNames)
	sqlClause, sqlVars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	rows, err := s.Raw(sqlClause, sqlVars...).Query()
	if err != nil {
		log.Errorln(err)
		return err
	}

	for rows.Next() {
		mod := reflect.New(modType).Elem()
		var values []interface{}
		for _, name := range table.FieldNames {
			values = append(values, mod.FieldByName(name).Addr().Interface())
		}

		err = rows.Scan(values...)
		if err != nil {
			log.Errorln(err)
			return err
		}

		s.CallMethod(AfterQuery, mod.Addr().Interface())
		// TODO: confused :(
		modSlice.Set(reflect.Append(modSlice, mod))
	}

	return rows.Close()
}

// First is to select the first record from database
func (s *Session) First(value interface{}) error {
	mod := reflect.Indirect(reflect.ValueOf(value))
	modSlice := reflect.New(reflect.SliceOf(mod.Type())).Elem()

	err := s.Limit(1).Find(modSlice.Addr().Interface())
	if err != nil {
		log.Errorln(err)
		return err
	}
	if modSlice.Len() == 0 {
		errMsg := "NOT FOUND"
		log.Errorln(errMsg)
		return errors.New(errMsg)
	}

	mod.Set(modSlice.Index(0))
	return nil
}

// Where is to add conditions for sqls
func (s *Session) Where(desc string, values ...interface{}) *Session {
	var sqlVars []interface{}
	s.clause.Set(clause.WHERE, append(append(sqlVars, desc), values...)...)
	return s
}

// Limit is to select limited number of records from database
func (s *Session) Limit(num int) *Session {
	s.clause.Set(clause.LIMIT, num)
	return s
}

// OrderBy is to add an order for sqls
func (s *Session) OrderBy(order string) *Session {
	s.clause.Set(clause.ORDERBY, order)
	return s
}

// Update is to update records in database
func (s *Session) Update(kvpair ...interface{}) (int64, error) {
	s.CallMethod(BeforeUpdate, nil)
	m, ok := kvpair[0].(map[string]interface{})
	if !ok {
		m = make(map[string]interface{})
		for i := 0; i < len(kvpair); i += 2 {
			m[kvpair[i].(string)] = kvpair[i+1]
		}
	}
	s.clause.Set(clause.UPDATE, s.Schema().Name, m)

	// TODO: confused :(
	sqlClause, sqlVars := s.clause.Build(clause.UPDATE, clause.WHERE)

	res, err := s.Raw(sqlClause, sqlVars...).Exec()
	if err != nil {
		log.Errorln(err)
		return 0, err
	}
	s.CallMethod(AfterUpdate, nil)
	return res.RowsAffected()
}

// Delete is to delete records in database
func (s *Session) Delete() (int64, error) {
	s.CallMethod(BeforeDelete, nil)
	s.clause.Set(clause.DELETE, s.Schema().Name)

	sqlClause, sqlVars := s.clause.Build(clause.DELETE, clause.WHERE)

	res, err := s.Raw(sqlClause, sqlVars).Exec()
	if err != nil {
		log.Errorln(err)
		return 0, err
	}
	s.CallMethod(AfterDelete, nil)
	return res.RowsAffected()
}

// Count is to return the number of records
func (s *Session) Count() (int64, error) {
	s.clause.Set(clause.COUNT, s.Schema().Name)

	sqlClause, sqlVars := s.clause.Build(clause.COUNT, clause.WHERE)

	row := s.Raw(sqlClause, sqlVars).QueryRow()
	var count int64
	err := row.Scan(&count)
	if err != nil {
		log.Errorln(err)
		return 0, err
	}
	return count, nil
}
