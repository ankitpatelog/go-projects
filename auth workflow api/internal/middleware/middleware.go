package middleware

import (
	"auth-workflow/internal/config"
	"net/http"
	"strings"
	"time"
)

func RateLimmiter(next http.Handler)http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//get client ip
		forwarded := r.Header.Get("X-forwardede-For")
		var clientIP string

		if forwarded!=""{
			parts := strings.Split(forwarded,",")

			if len(parts)<2 {
				http.Error(w,"client unauthorized",http.StatusBadGateway)
				return
			}
			clientIP = parts[0]
		}else{
			//if  client  ip not found then get network address
			clientIP = r.RemoteAddr
		}
		
		//generate key and increment for that key in db

		key := "Rate:"+clientIP

		count,err := config.Redis.Incr(config.Ctx,key).Result()
		if err!=nil {
			http.Error(w,"Some error occured",http.StatusInternalServerError)
			return
		}
		
		if count==1 {
			config.Redis.Expire(config.Ctx,key,time.Hour)
		}
		
		if count>20 {
			config.Redis.Del(config.Ctx,key)
			http.Error(w,"Too many request",http.StatusInternalServerError)
			return
		}
		next.ServeHTTP(w,r)
	})
}