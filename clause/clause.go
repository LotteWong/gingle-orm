package clause

import "strings"

// Clause is related to sql chained clause
type Clause struct {
	sqlClause map[Type]string
	sqlVars   map[Type][]interface{}
}

// Set is to generate single sql clause
func (c *Clause) Set(name Type, values ...interface{}) {
	if c.sqlClause == nil || c.sqlVars == nil {
		c.sqlClause = make(map[Type]string)
		c.sqlVars = make(map[Type][]interface{})
	}
	sqlClause, sqlVars := generators[name](values)
	c.sqlClause[name] = sqlClause
	c.sqlVars[name] = sqlVars
}

// Build is to generate chained sql clause
func (c *Clause) Build(types ...Type) (string, []interface{}) {
	var sqlClauses []string
	var sqlVars []interface{}

	for _, typ := range types {
		if sqlClause, ok := c.sqlClause[typ]; ok {
			sqlClauses = append(sqlClauses, sqlClause)
			if sqlVar, ok := c.sqlVars[typ]; ok {
				sqlVars = append(sqlVars, sqlVar)
			}
		}
	}

	return strings.Join(sqlClauses, " "), sqlVars
}
