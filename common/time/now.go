package time

import "time"

// Now use for unit test patch.
var Now = time.Now

var dtaTimeFormat = "2006-01-02 15:04:05"

// ToDateTime layout 2006-01-02 15:04:05
func ToDateTime(t time.Time) string {
	return t.Format(dtaTimeFormat)
}

// FromDateTime layout 2006-01-02 15:04:05
func FromDateTime(str string) (time.Time, error) {
	return time.Parse(dtaTimeFormat, str)
}

// NowDateTime layout 2006-01-02 15:04:05
func NowDateTime() string {
	return ToDateTime(Now())
}
