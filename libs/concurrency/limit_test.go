package concurrency

import (
	"golang.org/x/time/rate"
	"testing"
	"time"
)

func TestLimit(t *testing.T) {
	limiter := rate.NewLimiter(rate.Limit(1), 1)
	allow := limiter.AllowN(time.Now(), 2)
	t.Log(allow)
}
