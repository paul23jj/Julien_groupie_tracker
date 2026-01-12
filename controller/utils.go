package controller

import (
	"fmt"
	"time"
)

// FormatDate formate une date en "JJ mois AAAA" (mois en français).
func FormatDate(v interface{}) string {
	if v == nil {
		return ""
	}

	switch t := v.(type) {
	case time.Time:
		return fmt.Sprintf("%02d %s %d", t.Day(), frenchMonth(t.Month()), t.Year())
	case *time.Time:
		if t == nil {
			return ""
		}
		return fmt.Sprintf("%02d %s %d", t.Day(), frenchMonth(t.Month()), t.Year())
	}

	s, ok := v.(string)
	if !ok || s == "" {
		return ""
	}

	layouts := []string{
		"2006-01-02",
		time.RFC3339,
		"2006-01-02 15:04:05",
	}
	var parsed time.Time
	var err error
	for _, l := range layouts {
		parsed, err = time.Parse(l, s)
		if err == nil {
			return fmt.Sprintf("%02d %s %d", parsed.Day(), frenchMonth(parsed.Month()), parsed.Year())
		}
	}
	// si on n'a pas réussi à parser, renvoyer la chaîne brute
	return s
}

func frenchMonth(m time.Month) string {
	months := []string{
		"",
		"janvier",
		"février",
		"mars",
		"avril",
		"mai",
		"juin",
		"juillet",
		"août",
		"septembre",
		"octobre",
		"novembre",
		"décembre",
	}
	if int(m) >= 1 && int(m) <= 12 {
		return months[int(m)]
	}
	return ""
}
