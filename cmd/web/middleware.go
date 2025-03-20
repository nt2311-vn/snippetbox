package main

import "net/http"

func secureHeaders(next http.Handler) http.Handler {
	securedHeaders := map[string]string{
		"Content-Security-Policy": "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com",
		"Referrer-Policy":         "origin-when-cross-origin",
		"X-Content-Type-Options":  "nosniff",
		"X-Frame-Options":         "deny",
		"X-XSS-Protection":        "0",
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for k, v := range securedHeaders {
			w.Header().Set(k, v)
		}

		next.ServeHTTP(w, r)
	})
}
