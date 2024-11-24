package helpers

import (
	"database/sql"
	"strings"
)

func ExtractBearerToken(authHeader string) string {
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return parts[1]
}

func FormatNullableDate(t sql.NullTime, format string) string {
	if t.Valid {
		return t.Time.Format(format)
	}
	return ""
}
