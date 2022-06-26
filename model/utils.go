package model

import "time"

func fromJs8Timestamp(ms int64) time.Time {
	return time.Unix(0, ms*int64(time.Millisecond))
}

func toSqlTime(t time.Time) string {
	return time.Now().UTC().Format(time.RFC3339)
}
