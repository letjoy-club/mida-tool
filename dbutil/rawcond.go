package dbutil

import (
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm/clause"
)

type Tabler interface {
	Alias() string
	TableName() string
}

type Field struct {
	Field field.Expr
	Table Tabler
}

type rawCond struct {
	field.Field
	sql  string
	args []interface{}
}

func (m rawCond) BeCond() interface{} {
	args := []interface{}{}
	for _, v := range m.args {
		switch arg := v.(type) {
		case Field:
			column := clause.Column{
				Name: arg.Field.ColumnName().String(),
				Raw:  false,
			}
			if arg.Table != nil {
				column.Table = arg.Table.TableName()
				column.Alias = arg.Table.Alias()
			}
			args = append(args, column)
		case field.Expr:
			column := clause.Column{
				Name: arg.ColumnName().String(),
				Raw:  false,
			}
			args = append(args, column)
		default:
			args = append(args, v)
		}
	}

	expr := clause.NamedExpr{SQL: m.sql}
	expr.Vars = append(expr.Vars, args...)
	return expr
}

func (rawCond) CondError() error { return nil }

func RawCond(sql string, args ...interface{}) gen.Condition {
	return &rawCond{
		sql:  sql,
		args: args,
	}
}
