package gingleorm

import (
	"database/sql"
	"fmt"
	"gingle-orm/dialect"
	"gingle-orm/log"
	"gingle-orm/session"
	"strings"
)

// Engine is related to database lifecycle
type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

// NewEngine is to return a Engine instance
func NewEngine(driver, source string) (e *Engine, err error) {
	// Open database
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Errorln(err)
		return
	}

	// Ping database
	err = db.Ping()
	if err != nil {
		log.Errorln(err)
		return
	}

	dialect, ok := dialect.GetDialect(driver)
	if !ok {
		log.Errorf("dialect %s Not Found", driver)
		return
	}

	e = &Engine{
		db:      db,
		dialect: dialect,
	}
	log.Infoln("Connect database successfully")
	return
}

// Close is to close the database instance
func (e *Engine) Close() (err error) {
	// Close database
	err = e.db.Close()
	if err != nil {
		log.Errorln(err)
		return
	}

	log.Infoln("Close database successfully")
	return
}

// NewSession is to return a Session instance
func (e *Engine) NewSession() *session.Session {
	return session.New(e.db, e.dialect)
}

// TxFunc is to callback session
type TxFunc func(*session.Session) (interface{}, error)

// Transaction is to wrap session with transaction
func (e *Engine) Transaction(f TxFunc) (res interface{}, err error) {
	s := e.NewSession()

	// Transaction begins
	if err = s.Begin(); err != nil {
		log.Errorln(err)
		return nil, err
	}

	// Transaction commits or rollbacks
	defer func() {
		if p := recover(); p != nil {
			err = s.Rollback()
			log.Errorln(err)
			panic(p)
		} else if err != nil {
			err = s.Rollback()
			log.Errorln(err)
		} else {
			err = s.Commit()
			log.Errorln(err)
		}
	}()

	return f(s)
}

func difference(a []string, b []string) (diff []string) {
	mapB := make(map[string]bool)
	for _, v := range b {
		mapB[v] = true
	}
	for _, v := range a {
		if _, ok := mapB[v]; !ok {
			diff = append(diff, v)
		}
	}
	return
}

// Migrate is to live update database when fields add or delete
func (e *Engine) Migrate(value interface{}) error {
	_, err := e.Transaction(func(s *session.Session) (res interface{}, err error) {
		if !s.Model(value).ExistTable() {
			log.Infof("Table %s doesn't exist", s.Schema().Name)
			return nil, s.CreateTable()
		}

		table := s.Schema()
		rows, _ := s.Raw(fmt.Sprintf("SELETE * FROM %s LIMIT 1", table.Name)).Query()
		columns, _ := rows.Columns()
		addCols := difference(table.FieldNames, columns)
		delCols := difference(columns, table.FieldNames)
		log.Infof("Added cols %v; Deleted cols %v", addCols, delCols)

		for _, col := range addCols {
			f, _ := table.GetField(col)
			if _, err = s.Raw(fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", table.Name, f.Name, f.Type)).Exec(); err != nil {
				return nil, err
			}
		}

		if len(delCols) == 0 {
			return nil, nil
		}
		tmpTable := "tmp_" + table.Name
		fieldStr := strings.Join(table.FieldNames, ", ")
		s.Raw(fmt.Sprintf("CREATE TABLE %s AS SELECT %s FROM %s", tmpTable, fieldStr, table.Name))
		s.Raw(fmt.Sprintf("Drop TABLE %s", table.Name))
		s.Raw(fmt.Sprintf("ALTER TABLE %s RENAME TO %s", tmpTable, table.Name))
		_, err = s.Exec()
		return nil, err
	})

	if err != nil {
		log.Errorln(err)
		return err
	}
	return nil
}
