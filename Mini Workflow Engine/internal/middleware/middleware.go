package middleware

import (
	"net/http"
	"strings"
	"time"
	"workflow-engine/internal/config"
)

func RateLimitter(next http.Handler)http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		forwarded := r.Header.Get("X-Forwarded-For")
		var clientIP string
		if forwarded!="" {
			//means it coints the real ip and proxy ip
			parts := strings.Split(forwarded,",")
			clientIP = parts[0]
		}else{
			//means no ip found thn return the client network address
			clientIP = r.RemoteAddr
		}

		//add this key value pair in redis and adjyst the rate limmiter
		key := "Ratelimit:"+clientIP

		count,err := config.Redis.Incr(config.Ctx,key).Result()
		if err != nil {
			http.Error(w, "Redis error", http.StatusInternalServerError)
			return
		}

		if count==1 {
			//rate count stored for the first time
			config.Redis.Expire(config.Ctx,key,time.Minute)
		}

		if count>25 {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w,r)
	})
}