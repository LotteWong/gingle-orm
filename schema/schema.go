package schema

import (
	"gingle-orm/dialect"
	"go/ast"
	"reflect"
)

// Field is to convert struct member to database column
type Field struct {
	Name string
	Type string
	Tag  string
}

// Schema is to convert go interface to database table
type Schema struct {
	Model      interface{}
	Name       string
	FieldNames []string
	Fields     []*Field
	fieldMap   map[string]*Field
}

// SetField is to set kvpair to field map
func (s *Schema) SetField(name string, field *Field) {
	s.fieldMap[name] = field
}

// GetField is to get kvpair from field map
func (s *Schema) GetField(name string) (field *Field, ok bool) {
	field, ok = s.fieldMap[name]
	return
}

// Parse is to do object relation mapping
func Parse(mod interface{}, dial dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(mod)).Type()
	schema := &Schema{
		Model:    mod,
		Name:     modelType.Name(),
		fieldMap: make(map[string]*Field),
	}

	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)

		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				Name: p.Name,
				Type: dial.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}

			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.fieldMap[p.Name] = field

			if val, ok := p.Tag.Lookup("gingleorm"); ok {
				field.Tag = val
			}
		}
	}

	return schema
}
