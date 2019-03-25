package convert

import "time"

func IntToPtr(i int) *int {
	return &i
}

func Int64ToPtr(i int64) *int64 {
	return &i
}

func StrToPtr(s string) *string {
	return &s
}

func TimeToPtr(t time.Time) *time.Time {
	return &t
}
