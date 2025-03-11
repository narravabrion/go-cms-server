package main

import "net/http"

func (api *api) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request, retryAfter string) {
	api.logger.Warnw("rate limit exceeded", "method", r.Method, "path", r.URL.Path)

	w.Header().Set("Retry-After", retryAfter)

	writeJSONError(w, http.StatusTooManyRequests, "rate limit exceeded, retry after: "+retryAfter)
}