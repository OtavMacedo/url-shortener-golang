package middleware

import (
	"net/http"

	"github.com/OtavMacedo/url-shortener-golang/internal/cache"
	"github.com/gin-gonic/gin"
)

func RateLimit(cache *cache.RedisCache) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()

		limited, err := cache.IsRateLimited(ctx, ip)
		if err != nil {
			ctx.Next()
			return
		}

		if limited {
			ctx.Header("Retry-After", "60")
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests, try again in 60 seconds",
			})
			return
		}

		ctx.Next()
	}
}
