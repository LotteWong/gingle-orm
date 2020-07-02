package clause

import (
	"fmt"
	"strings"
)

type Type int

const (
	INSERT Type = iota
	VALUES
	SELECT
	LIMIT
	WHERE
	ORDERBY
	UPDATE
	DELETE
	COUNT
)

type generator func(values ...interface{}) (string, []interface{})

var generators map[Type]generator

func init() {
	generators = make(map[Type]generator)
	generators[INSERT] = _insert
	generators[VALUES] = _values
	generators[SELECT] = _select
	generators[WHERE] = _where
	generators[LIMIT] = _limit
	generators[ORDERBY] = _orderBy
}

func genBindVars(num int) string {
	vars := []string{}
	for i := 0; i < num; i++ {
		vars = append(vars, "?")
	}
	return strings.Join(vars, ", ")
}

func _insert(values ...interface{}) (string, []interface{}) {
	// INSERT INTO $tableName ($fields)
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("INSERT INTO %s (%v)", tableName, fields), []interface{}{}
}

func _values(values ...interface{}) (string, []interface{}) {
	// VALUES ($v1), ($v2), ...
	var bindStr string
	var sqlClause strings.Builder
	var sqlVars []interface{}

	sqlClause.WriteString("VALUES ")
	for idx, val := range sqlVars {
		v := val.([]interface{})

		if bindStr == "" {
			bindStr = genBindVars(len(v)) // placeholder
		}
		sqlClause.WriteString(fmt.Sprintf("(%v)", bindStr))
		if idx != len(values)-1 {
			sqlClause.WriteString(", ")
		}

		sqlVars = append(sqlVars, v...) // real params
	}

	return sqlClause.String(), sqlVars
}

func _select(values ...interface{}) (string, []interface{}) {
	// SELECT $fields FROM $tableName
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ", ")
	return fmt.Sprintf("SELECT %v FROM %s", fields, tableName), []interface{}{}
}

func _where(values ...interface{}) (string, []interface{}) {
	// WHERE $cond
	cond, vals := values[0], values[1:]
	return fmt.Sprintf("WHERE %s", cond), vals
}

func _limit(values ...interface{}) (string, []interface{}) {
	// LIMIT $num
	return "LIMIT ?", values
}

func _orderBy(values ...interface{}) (string, []interface{}) {
	// ORDER BY $order
	order := values[0]
	return fmt.Sprintf("ORDER BY %s", order), []interface{}{}
}
