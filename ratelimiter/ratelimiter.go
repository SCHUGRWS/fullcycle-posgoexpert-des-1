package ratelimiter

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type RateLimiter struct {
	client *redis.Client
}

func NewRateLimiter(client *redis.Client) *RateLimiter {
	return &RateLimiter{
		client: client,
	}
}

func getLimitForKey(key string, isToken bool) int {
	var defaultLimit int
	if isToken {
		defaultLimit = 20
	} else {
		defaultLimit = 5
	}

	if isToken {
		tokens, ok := viper.Get("rate_limit.tokens").([]interface{})
		if !ok {
			return defaultLimit
		}

		for _, t := range tokens {
			tokenConfig, ok := t.(map[string]interface{})
			if !ok {
				continue
			}
			if tokenConfig["token"].(string) == key {
				limit, ok := tokenConfig["limit"].(int)
				if !ok {
					continue
				}

				return limit
			}
		}
	} else {
		ips, ok := viper.Get("rate_limit.ips").([]interface{})
		if !ok {
			return defaultLimit
		}

		for _, ip := range ips {
			ipConfig, ok := ip.(map[string]interface{})
			if !ok {
				continue
			}
			if ipConfig["ip"].(string) == key {
				limit, ok := ipConfig["limit"].(int)
				if !ok {
					continue
				}

				return limit
			}
		}
	}

	return defaultLimit
}

func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("API_KEY")
		var key string
		var isToken bool
		if token != "" {
			key = token
			isToken = true
		} else {
			requesterIP := r.Header.Get("X-Forwarded-For")
			if requesterIP == "" {
				requesterIP = r.RemoteAddr
			}
			key = strings.Split(requesterIP, ":")[0]
			isToken = false
		}

		limit := getLimitForKey(key, isToken)

		count, err := rl.client.Get(key).Int()
		if err == redis.Nil {
			count = 0
		} else if err != nil {
			fmt.Printf("Error accessing Redis: %v\n", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if count >= limit {
			http.Error(w, "You have reached the maximum number of requests allowed within a certain time frame", http.StatusTooManyRequests)
			return
		}

		_, err = rl.client.Incr(key).Result()
		if err != nil {
			fmt.Printf("Error incrementing count in Redis: %v\n", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if count == 0 {
			timeout := getenvInt("REDIS_EXPIRE_TIMEOUT_SECONDS")
			rl.client.Expire(key, time.Second*time.Duration(timeout)).Result()
		}

		next.ServeHTTP(w, r)
	})
}

func getenvInt(key string) int {
	v, _ := strconv.Atoi(os.Getenv(key))
	return v
}
