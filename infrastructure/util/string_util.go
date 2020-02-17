package util

import "database/sql"

//Convert Nullstring to string
func NullStringToString(s sql.NullString) string {
	if s.Valid {
		return s.String
	}
	return ""
}
