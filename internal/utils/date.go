package utils

import "time"

func FormatDate(date time.Time) string {
	return date.Format("2006-01-02 15:04:05")
}

func ParseDateToTimestamp(dateStr string) (int64, error) {
	// 解析日期字符串
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return 0, err
	}

	// 转换为秒级时间戳
	return t.Unix(), nil
}
