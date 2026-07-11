package controllers

import (
	"os"
	"strings"
)

func resolveAdminStatus(stored bool, email string) bool {
	if stored {
		return true
	}
	return isBootstrapAdmin(email)
}

func isBootstrapAdmin(email string) bool {
	raw := os.Getenv("ADMIN_EMAILS")
	if raw == "" {
		return false
	}

	email = strings.ToLower(strings.TrimSpace(email))
	for part := range strings.SplitSeq(raw, ",") {
		if strings.ToLower(strings.TrimSpace(part)) == email {
			return true
		}
	}
	return false
}
