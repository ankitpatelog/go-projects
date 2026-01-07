package middleware

import (
	"event-processing-service/internal/config"
	"net/http"
	"strings"
	"time"
)

func RateLimiter(next http.Handler)http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//get client ip address
// 		X-Forwarded-For: 203.0.113.5, 10.0.0.1
// 		Meaning:
// 		203.0.113.5  → real client IP
//		 10.0.0.1     → proxy IP

		var clientIP string

		forwarded := r.Header.Get("X-Forwarded-For")
		if forwarded!="" {
			//means it coints the real ip and proxy ip
			parts := strings.Split(forwarded,",")
			clientIP = parts[0]
		}else{
			//means no ip found thn return the client network address
			clientIP = r.RemoteAddr
		}

		//key value pair for redis data 
		key := "rate:" + clientIP

		count,err := config.Redis.Incr(config.Ctx,key).Result()
		if err != nil {
			http.Error(w, "Redis error", http.StatusInternalServerError)
			return
		}

		if count ==1 {
			config.Redis.Expire(config.Ctx,key,time.Minute)
		}

		if count>100 {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		//allow next request
		next.ServeHTTP(w,r)

	})
	

}