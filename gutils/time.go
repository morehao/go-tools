package gutils

import "time"

const (
	YYYY_MM_DD_HH_MM_SS = "2006-01-02 15:04:05"
	YYYY_MM_DD          = "2006-01-02"
	MM_DD               = "01-02"
	YYYYMMDD            = "20060102"
	MMDD                = "0102"
)

func TimeFormat(t time.Time, format string) string {
	if t.Unix() <= 0 {
		return ""
	}
	return t.Format(format)
}
