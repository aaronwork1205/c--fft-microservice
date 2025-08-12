package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"golang.org/x/time/rate"
)

func main() {
	u, _ := url.Parse("http://127.0.0.1:9000")         // 你的 C++ 服务
	proxy := httputil.NewSingleHostReverseProxy(u)

	lim := rate.NewLimiter(100, 100)                   // 全局 QPS≈100
	sem := make(chan struct{}, 256)                    // 同时最多 256 个请求

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 并发节流
		sem <- struct{}{}
		defer func(){ <-sem }()

		// 限流：没令牌就 429
		if !lim.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		proxy.ServeHTTP(w, r) // 转发
	})

	http.ListenAndServe(":8080", handler)
}
