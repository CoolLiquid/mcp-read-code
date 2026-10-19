package main

import (
	"net/http"
	"strings"
)

func extractGitLabToken(r *http.Request) string {
	if t := strings.TrimSpace(r.Header.Get("X-Gitlab-Token")); t != "" {
		return t
	}
	auth := r.Header.Get("Authorization")
	if len(auth) > 7 && strings.EqualFold(auth[:7], "Bearer ") {
		return strings.TrimSpace(auth[7:])
	}
	return ""
}

func gitLabURLFromRequest(r *http.Request, defaultURL string) string {
	if u := strings.TrimSpace(r.Header.Get("X-Gitlab-Url")); u != "" {
		return u
	}
	return defaultURL
}

type requireGitLabToken struct {
	inner      http.Handler
	defaultURL string
}

func (m *requireGitLabToken) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if extractGitLabToken(r) == "" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("missing GitLab token: set X-Gitlab-Token or Authorization: Bearer <token>\n"))
		return
	}
	m.inner.ServeHTTP(w, r)
}
