package good_util

import (
	"github.com/samber/lo"
	"time"
)

func loDuration() {
	duration := lo.Duration(func() {
		time.Sleep(1 * time.Second)
	})
	println(duration)
}
