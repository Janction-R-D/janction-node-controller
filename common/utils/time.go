package utils

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	LayoutNormal = "2006-01-02 15:04:05"
)

func ParseLayoutNormalIgnoreNilStr(t string) (time.Time, error) {
	if t == "" {
		return time.Unix(0, 0), nil
	}
	return time.Parse(LayoutNormal, t)
}

func IsZeroTime(t time.Time) bool {
	return t.Local().Format("2006-01-02 15:04:05") == time.Time{}.Local().Format("2006-01-02 15:04:05")
}

func IsZeroDate(t time.Time) bool {
	return t.Local().Format("2006-01-02") == time.Time{}.Local().Format("2006-01-02")
}

func ConvertToPbTime(from *time.Time) (to *timestamppb.Timestamp) {
	if from == nil {
		return nil
	}
	return timestamppb.New(*from)
}

func FormatTime(t *time.Time, layout ...string) string {
	if t == nil {
		return ""
	}
	if len(layout) == 0 {
		return t.Format("2006-01-02 15:04:05")
	}
	return t.Format(layout[0])
}

func GetTimeStamp(t *time.Time) int64 {
	if t == nil {
		return 0
	}
	return t.Unix()
}

// 获取某一天的0点时间
func GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

func FormatTimeStrToLocalRFC3339(layout string, v *string) string {
	oriTm := TypeValue(v)
	t, err := time.ParseInLocation(layout, oriTm, time.Local)
	if err != nil {
		return oriTm
	}
	return t.Local().Format(time.RFC3339)
}

func FormatTimeStrToLocalRFC3339V2(layout string, v *string) string {
	oriTm := TypeValue(v)
	t, err := time.ParseInLocation(layout, oriTm, time.UTC)
	if err != nil {
		return oriTm
	}
	return t.UTC().Format(time.RFC3339)
}

func FormatRFC3339ToNormal(layout string, v *string) string {
	oriTm := TypeValue(v)
	t, err := time.ParseInLocation(layout, oriTm, time.UTC)
	if err != nil {
		return oriTm
	}
	return t.UTC().Format(LayoutNormal)
}

func FormatTimeToNormal(t *time.Time) string {
	oriTm := TypeValue(t)
	return oriTm.UTC().Format(LayoutNormal)
}

func FormatTimeToLocalRFC3339(t *time.Time) string {
	oriTm := TypeValue(t)
	return oriTm.Local().Format(time.RFC3339)
}
