package config

import (
	"os"
	"strings"
)

func AllowedOrigins() ([]string, map[string]bool) {
	raw := os.Getenv("CORS_ORIGINS")

	list := strings.Split(raw, ",")
	exact := make(map[string]bool)

	for _, o := range list {
		o = strings.TrimSpace(o)
		if o != "" {
			exact[o] = true
		}
	}

	return list, exact
}

func LoadCorsOrigins() map[string]bool {
	raw := os.Getenv("CORS_ORIGINS")
	out := make(map[string]bool)

	for _, o := range strings.Split(raw, ",") {
		o = strings.TrimSpace(o)
		if o != "" {
			out[o] = true
		}
	}

	return out
}