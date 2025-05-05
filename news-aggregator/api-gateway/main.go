package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var (
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	ctx = context.Background()
)

const RATE_LIMIT = 5 // số request mỗi phút

func rateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := fmt.Sprintf("rate_limit:%s", ip)

		val, err := rdb.Get(ctx, key).Int()
		if err != nil && err != redis.Nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Redis error"})
			return
		}

		if val >= RATE_LIMIT {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			return
		}

		pipe := rdb.TxPipeline()
		pipe.Incr(ctx, key)
		pipe.Expire(ctx, key, time.Minute)
		_, _ = pipe.Exec(ctx)

		c.Next()
	}
}

func main() {
	r := gin.Default()

	// Middleware giới hạn rate
	r.Use(rateLimitMiddleware())

	// Route chuyển tiếp tới article-service
	r.GET("/api/articles", func(c *gin.Context) {
		resp, err := http.Get("http://localhost:8081/articles")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
			return
		}
		defer resp.Body.Close()

		c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
	})

	fmt.Println("🚪 API Gateway chạy tại cổng 8080")
	r.Run(":8080")
}
