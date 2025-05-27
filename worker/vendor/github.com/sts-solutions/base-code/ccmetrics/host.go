package ccmetrics

import (
	"regexp"
	"strings"
)

type Host string

func (h Host) String() string {
	return string(h)
}

func (h Host) Shard() string {
	host := h.String()
	parts := strings.Split(host, ".")
	if len(parts) == 2 {
		s := regexp.MustCompile(`[0-9]+`).FindString(parts[0])
		if s != "" {
			return s
		}
	}
	return host
}
