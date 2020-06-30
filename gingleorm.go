package gingleorm

import (
	"database/sql"
	"gingle-orm/log"
	"gingle-orm/session"
)

// Engine is related to database lifecycle
type Engine struct {
	db *sql.DB
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

	e = &Engine{
		db: db,
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
	return session.New(e.db)
}
