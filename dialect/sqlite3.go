package dialect

import (
	"fmt"
	"reflect"
	"time"
)

type sqlite3 struct{}

// Check type conversion statically
var _ Dialect = (*sqlite3)(nil)

func init() {
	SetDialect("sqlite3", &sqlite3{})
}

func (s *sqlite3) DataTypeOf(typ reflect.Value) string {
	switch typ.Kind() {
	case reflect.Bool:
		return "bool"
	case reflect.Int, reflect.Int8, reflect.Int32, reflect.Uintptr,
		reflect.Uint, reflect.Uint16, reflect.Uint32:
		return "integer"
	case reflect.Int64, reflect.Uint64:
		return "bigint"
	case reflect.Float32, reflect.Float64:
		return "real"
	case reflect.String:
		return "text"
	case reflect.Array, reflect.Slice:
		return "blob"
	case reflect.Struct:
		if _, ok := typ.Interface().(time.Time); ok {
			return "datetime"
		}
	}

	panic(fmt.Sprintf("invalid sql type %s (%s)", typ.Type().Name(), typ.Kind()))
}

func (s *sqlite3) TableExistSQL(tableName string) (string, []interface{}) {
	values := []interface{}{tableName}
	// TODO: Prevent SQL injection
	return "SELECT name FROM sqlite_master WHERE type='table' AND name = ?", values
}
