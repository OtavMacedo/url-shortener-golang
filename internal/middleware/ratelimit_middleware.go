package middleware

import (
	"net/http"

	"github.com/OtavMacedo/url-shortener-golang/internal/ratelimiter"
	"github.com/gin-gonic/gin"
)

func RateLimiterMiddleware(rl *ratelimiter.RateLimiter) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()
		allowed, err := rl.Allow(ctx.Request.Context(), ip)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		if !allowed {
			ctx.Header("Retry-After", "60")
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests, try again later",
			})
			return
		}
		ctx.Next()
	}
}
