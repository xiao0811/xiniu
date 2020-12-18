package model

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// TimeFormat 时间格式
var TimeFormat = "2006-01-02 15:04:05"

// MyTime 数据库时间格式
type MyTime struct {
	time.Time
}

// MarshalJSON 格式化成json
func (t MyTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", t.Format(TimeFormat))
	return []byte(formatted), nil
}

// Value 获取时间值
func (t MyTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan .
func (t *MyTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = MyTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

// UnmarshalJSON .
// func (t *MyTime) UnmarshalJSON(data []byte) error {
// 	if string(data) == "null" {
// 		return nil
// 	}
// 	var err error
// 	//前端接收的时间字符串
// 	str := string(data)
// 	//去除接收的str收尾多余的"
// 	timeStr := strings.Trim(str, "\"")
// 	t1, err := time.ParseInLocation(TimeFormat, timeStr, time.Local)
// 	*t = MyTime{t1}
// 	return err
// }
