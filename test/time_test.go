package test

import (
	"fmt"
	"testing"
	"time"
)

func TestTimeNow(t *testing.T) {
	now := time.Now()
	curYearMonth := now.In(time.Local).Format("2006_01")

	fmt.Println(now, curYearMonth)
}
