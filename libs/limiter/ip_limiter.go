package limiter

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"golib/libs/gin_helper"
	"net/http"
	"strings"
	"sync"
)

// 令牌桶容量：20，每秒扔进的令牌数：10
var limiter = NewIPRateLimiter(10, 20)

func IpLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		remoteAddr := c.ClientIP()
		addrPort := strings.Split(remoteAddr, ":")
		limit := limiter.GetLimiter(addrPort[0])
		if !limit.Allow() {
			myErr := gin_helper.LIMIT_ERROR
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
				"code":       myErr.Code,
				"msg":        myErr.Msg,
				"serverTime": myErr.ServerTime,
				"data":       myErr.Data,
			})
			return
		}

	}
}

// IPRateLimiter .
type IPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  *sync.RWMutex
	r   rate.Limit
	b   int
}

// NewIPRateLimiter .
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	i := &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b,
	}

	return i
}

// AddIP 创建了一个新的速率限制器，并将其添加到 ips 映射中,
// 使用 IP地址作为密钥
func (i *IPRateLimiter) AddIP(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limite := rate.NewLimiter(i.r, i.b)

	i.ips[ip] = limite

	return limite
}

// GetLimiter 返回所提供的IP地址的速率限制器(如果存在的话).
// 否则调用 AddIP 将 IP 地址添加到映射中
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	limite, exists := i.ips[ip]

	if !exists {
		i.mu.Unlock()
		return i.AddIP(ip)
	}

	i.mu.Unlock()

	return limite
}
