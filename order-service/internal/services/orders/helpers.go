package orders

import "strings"

func key(idempotenceKey, userID string) string {
	sb := strings.Builder{}

	sb.WriteString(idempotenceKey)
	sb.WriteString("_")
	sb.WriteString(userID)

	return sb.String()
}
