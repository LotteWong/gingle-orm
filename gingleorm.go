package gingleorm

import (
	"database/sql"
	"gingle-orm/dialect"
	"gingle-orm/log"
	"gingle-orm/session"
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
		log.Errorln("dialect %s Not Found", driver)
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
