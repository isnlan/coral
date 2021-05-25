package utils

import (
	"time"
)

func GetHourTime() int64 {
	now := time.Now()
	zero := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())
	return zero.Unix()
}

func MakeHourBegin(begin time.Time) int64 {
	return begin.Unix() - int64(begin.Minute()*60) - int64(begin.Second())
}

func MakeDayBegin(begin time.Time) int64 {
	return begin.Unix() - int64(begin.Hour()*60*60) - int64(begin.Minute()*60) - int64(begin.Second())
}

// 获取当前时间最接近的10分钟， 例如：14:35 -14:30, 11:07->11:10
func Make10MinuteBegin(begin time.Time) int64 {
	hour := begin.Unix() - int64(begin.Minute()*60) - int64(begin.Second())
	m := int64(begin.Minute() / 10)
	return hour + m*600
}

func GetTime(i int64) time.Time {
	return time.Unix(i, 0)
}

func Millisecond(t time.Time) int64 {
	return t.UnixNano() / 1e6
}
