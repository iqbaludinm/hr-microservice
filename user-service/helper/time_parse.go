package helper

import (
	"fmt"
	"strings"
	"time"
)

func ParseTimeToFullIndonesian(t time.Time) string {
	// dayNames := strings.Split("Minggu Senin Selasa Rabu Kamis Jumat Sabtu", " ")
	monthNames := strings.Split("Januari Februari Maret April Mei Juni Juli Agustus September Oktober November Desember", " ")

	// day := dayNames[t.Weekday()]
	month := monthNames[t.Month()-1]
	return fmt.Sprintf("%d %s %d", t.Day(), month, t.Year())
}

func ParseTimeToShortIndonesian(t time.Time) string {
	// dayNames := strings.Split("Minggu Senin Selasa Rabu Kamis Jumat Sabtu", " ")
	monthNames := strings.Split("Jan Feb Mar Apr May Jun Jul Aug Sep Okt Nov Des", " ")

	// day := dayNames[t.Weekday()]
	month := monthNames[t.Month()-1]
	return fmt.Sprintf("%d-%s-%d", t.Day(), month, t.Year())
}
