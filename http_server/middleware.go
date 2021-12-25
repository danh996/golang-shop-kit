package http_server

import (
	"context"
	"net/http"
	"strings"

	"gitlab.com/canco1/canco-kit/requestinfo"
	"gitlab.com/canco1/canco-kit/token"

	"golang.org/x/time/rate"

	"google.golang.org/grpc/metadata"
)

type mapMetaDataFunc func(context.Context, *http.Request) metadata.MD

var limiter = rate.NewLimiter(1, 3)

// MapMetaDataWithBearerToken ...
func MapMetaDataWithBearerToken(authenticator token.Authenticator) mapMetaDataFunc {
	return func(ctx context.Context, r *http.Request) metadata.MD {
		md := metadata.MD{}
		authorization := r.Header.Get(requestinfo.Authorization)

		if authorization != "" {
			bearerToken := strings.Split(authorization, requestinfo.Bearer+" ")
			if len(bearerToken) < 2 {
				return md
			}
			token := bearerToken[1]
			payload, err := authenticator.Verify(token)
			if err == nil {
				if m, err := payload.ParseToJSONString(); err != nil {
					md.Append(requestinfo.InfoKey, string(m))
				}
			}
		}
		return md
	}
}

// MapMetaDataWithCookies ...
func MapMetaDataWithCookies(authenticator token.Authenticator) mapMetaDataFunc {
	return func(ctx context.Context, r *http.Request) metadata.MD {
		md := metadata.MD{}

		cookies := r.Cookies()
		for _, cookie := range cookies {
			if cookie.Name == requestinfo.CookieKey {
				payload, err := authenticator.Verify(cookie.Value)
				if err == nil {
					if m, err := payload.ParseToJSONString(); err != nil {
						md.Append(requestinfo.InfoKey, string(m))
					}
				}
			}
		}
		return md
	}
}

func AllowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, authorization")
		if r.Method != "OPTIONS" {
			h.ServeHTTP(w, r)
		}
	})
}

// GetClientIP get client IP from HTTP request
func GetClientIP(req *http.Request) string {
	clientIP := req.Header.Get("HTTP_X_FORWARDED_FOR")
	if len(clientIP) == 0 {
		clientIP = req.Header.Get("X-Forwarded-For")
	}
	if len(clientIP) == 0 {
		clientIP = req.Header.Get("REMOTE_ADDR")
	}
	return clientIP
}

func limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
