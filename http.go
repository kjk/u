package u

import (
	"net/http"
	"strings"
)

// Request.RemoteAddress contains port, which we want to remove i.e.:
// "[::1]:58292" => "[::1]"
func ipAddrFromRemoteAddr(s string) string {
	idx := strings.LastIndex(s, ":")
	if idx == -1 {
		return s
	}
	return s[:idx]
}

// RequestGetRemoteAddress returns ip address of the client making the request,
// taking into account http proxies
func RequestGetRemoteAddress(r *http.Request) string {
	hdr := r.Header
	hdrRealIP := hdr.Get("X-Real-Ip")
	hdrForwardedFor := hdr.Get("X-Forwarded-For")
	if hdrRealIP == "" && hdrForwardedFor == "" {
		return ipAddrFromRemoteAddr(r.RemoteAddr)
	}
	if hdrForwardedFor != "" {
		// X-Forwarded-For is potentially a list of addresses separated with ","
		parts := strings.Split(hdrForwardedFor, ",")
		for i, p := range parts {
			parts[i] = strings.TrimSpace(p)
		}
		// TODO: should return first non-local address
		return parts[0]
	}
	return hdrRealIP
}

// RequestGetProtocol returns protocol under which the request is being served i.e. "http" or "https"
func RequestGetProtocol(r *http.Request) string {
	hdr := r.Header
	// X-Forwarded-Proto is set by proxies e.g. CloudFlare
	forwardedProto := strings.TrimSpace(strings.ToLower(hdr.Get("X-Forwarded-Proto")))
	if forwardedProto != "" {
		if forwardedProto == "http" || forwardedProto == "https" {
			return forwardedProto
		}
	}
	if r.TLS != nil {
		return "https"
	}
	return "http"
}

// RequestGetFullHost returns full host name e.g. "https://blog.kowalczyk.info/"
func RequestGetFullHost(r *http.Request) string {
	return RequestGetProtocol(r) + "://" + r.Host
}
