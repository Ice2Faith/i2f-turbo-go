// TimeProc 时间相关处理
package main

import (
	"fmt"
	"time"
)

//格林尼治时间
func UTCTime() {
	tm := time.Now()
	fmt.Printf("%d-%d-%d %d:%d:%d %d\n", tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute(), tm.Second(), tm.Nanosecond())
	fmt.Println("Unix:", tm.Unix())
}
func isLeapYear(year int64) (result bool) {
	if year%400 == 0 || (year%4 == 0 && year%100 != 0) {
		result = true
	} else {
		result = false
	}
	return
}
func DifferTimeSecond(tm1 time.Time, tm2 time.Time) (diffTime int64) {
	diffTime = tm2.Unix() - tm1.Unix()
	return
}

//计算你的出生日期到现在经过的时间
func BirthBetweenTime(birthTm time.Time) (seconds int64, minus int64, hours int64, days int64, months int64, years int64) {
	nowTm := time.Now()
	seconds = nowTm.Unix() - birthTm.Unix()
	minus = seconds / 60
	hours = seconds / 60 / 60
	days = seconds / 60 / 60 / 24
	months = int64(12 - birthTm.Month())
	var index int64 = int64(birthTm.Year())
	for index = index + 1; index < int64(nowTm.Year()); index++ {
		months += int64(12)
	}
	months += int64(nowTm.Month())
	years = int64(nowTm.Year() - birthTm.Year())
	return
}

func main() {
	UTCTime()
	fmt.Println("now year is leap year:", isLeapYear(int64(time.Now().Year())))
	//构造时间，使用本地时区，月份参数正确应该使用：time.July
	birthTm := time.Date(1999, time.Month(2), 26, 5, 15, 0, 0, time.Local)
	seconds, minus, hours, days, months, years := BirthBetweenTime(birthTm)
	fmt.Println(seconds, minus, hours, days, months, years)
}
