package times

import (
	"strings"
	"time"
)

const (
	TmFmtWithMS1 = "2006-01-02 15:04:05.000"
	TmFmtWithMS2 = "2006/01/02 15:04:05.000"
	TmFmtWithS1  = "2006-01-02 15:04:05"
	TmFmtWithS2  = "2006/01/02 15:04:05"
	TmFmtWithD1  = "2006-01-02"
	TmFmtWithD2  = "2006/01/02"
)

var (
	Locate, _ = time.LoadLocation("Asia/Shanghai")
)

func SubMinuteByTime(t time.Time) float64 {
	return time.Now().Sub(t).Minutes()
}

func SubMinuteByString(t string) (float64, error) {
	timeObj, err := time.ParseInLocation(TmFmtWithS1, t, Locate)
	if err != nil {
		return 0, err
	}
	return time.Now().Sub(timeObj).Minutes(), nil
}

func String2Time(s string) (time.Time, error) {
	timeObj, err := time.ParseInLocation(TmFmtWithS1, s, Locate)
	if err != nil {
		return time.Time{}, err
	}
	return timeObj, nil
}

func FormatSubTime(t string) string {
	result := strings.Replace(t, "h", "时", 1)
	result = strings.Replace(result, "m", "分钟", 1)
	lastIdx := strings.LastIndex(result, "钟")
	if lastIdx > 0 {
		return result[:lastIdx]
	} else {
		return "1分"
	}
}
