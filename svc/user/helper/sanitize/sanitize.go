package sanitize

import (
	"github.com/microcosm-cc/bluemonday"
)

func StrictSanitize(raw string) string {
	p := bluemonday.StrictPolicy()
	return p.Sanitize(raw)
}

func GeneralSanitize(raw string) string {
	p := bluemonday.UGCPolicy()
	return p.Sanitize(raw)
}
