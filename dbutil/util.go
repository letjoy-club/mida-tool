package dbutil

import (
	"strings"

	"github.com/feiin/sqlstring"
)

func Contains(column string, value any) string {
	valueStr := sqlstring.Escape(value)
	return "JSON_CONTAINS(" + column + ", '" + strings.ReplaceAll(valueStr, "'", "\"") + "')"
}

func ArrayAppend(column string, value any) string {
	valueStr := sqlstring.Escape(value)
	return "JSON_ARRAY_APPEND(" + column + ", '" + valueStr + "')"
}
