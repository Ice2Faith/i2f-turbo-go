// dateproc
package packutil

//方法大写字母开头，标识包外部可见方法public，小写开头，外部不可见private
func IsLeapYear(year int) (leapYear bool) {
	if (year%400 == 0) || (year%4 == 0 && year%100 != 0) {
		leapYear = true
	} else {
		leapYear = false
	}
	return
}
