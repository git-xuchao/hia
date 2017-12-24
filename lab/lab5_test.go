package lab

import (
	"fmt"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	/*
	 *a := time.Now()
	 *fmt.Println(a)
	 */

	t2, err := time.Parse("2006-01-02 15:04:05", "2016-07-27 8:46:15")
	if err != nil {
	}
	fmt.Println(t2)
	year, _, _ := t2.Date()
	hour, _, _ := t2.Clock()

	fmt.Println("year", year, "hour", hour)
}
