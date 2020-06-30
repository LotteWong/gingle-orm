package dialect

import "reflect"

var dialectMap = map[string]Dialect{}

// Dialect is related to transformation and compatiblity
type Dialect interface {
	// DataTypeOf is to convert go type to sql type
	DataTypeOf(typ reflect.Value) string
	// TableExist is to check whether the table exists
	TableExistSQL(tableName string) (string, []interface{})
}

// SetDialect is to set kvpair to dialect map
func SetDialect(name string, dialect Dialect) {
	dialectMap[name] = dialect
}

// GetDialect is to get kvpair from dialect map
func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectMap[name]
	return
}
