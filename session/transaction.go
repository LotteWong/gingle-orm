package session

import "gingle-orm/log"

// Begin is to start a transaction
func (s *Session) Begin() (err error) {
	log.Infoln("TRANSACTION BEGIN")
	if s.tx, err = s.db.Begin(); err != nil {
		return
	}
	return
}

// Commit is to commit a transaction
func (s *Session) Commit() (err error) {
	log.Infoln("TRANSACTION COMMIT")
	if err = s.tx.Commit(); err != nil {
		log.Errorln(err)
		return
	}
	return
}

// Rollback is to rollback a transaction
func (s *Session) Rollback() (err error) {
	log.Infoln("TRANSACTION ROLLBACK")
	if err = s.tx.Rollback(); err != nil {
		log.Errorln(err)
		return
	}
	return
}
