package timezone

import (
	"time"
)

var WIB *time.Location

func InitTimeLocation() {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		panic("failed to load Asia/Jakarta timezone: " + err.Error())
	}
	WIB = loc
}
