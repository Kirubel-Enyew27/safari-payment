package utils

import (
	"fmt"
	"time"
)

func ParseMpesaDate(val any) time.Time {
	str := fmt.Sprintf("%v", val)
	t, _ := time.Parse("20060102150405", str) // M-Pesa format
	return t
}
