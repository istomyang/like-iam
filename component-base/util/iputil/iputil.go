package iputil

import (
	"net"
	"net/http"
)

const (
	XForwardedFor = "X-Forwarded-For"
	XRealIP       = "X-Real-IP"
	XClientIP     = "x-client-ip"
)

// RemoteIP return remote client ip throw X-Forwarded-For, X-Real-IP and x-client-ip.
// // https://zhuanlan.zhihu.com/p/21354318
func RemoteIP(r *http.Request) string {
	remoteAddr := r.RemoteAddr

	if ip := r.Header.Get(XClientIP); ip != "" {
		remoteAddr = ip
	} else if ip := r.Header.Get(XRealIP); ip != "" {
		remoteAddr = ip
	} else if ip = r.Header.Get(XForwardedFor); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}

	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}

	return remoteAddr
}

// GetLocalIP return local net ip.
// refer to sonyflake.lower16BitPrivateIP
func GetLocalIP() (ip string) {
	ip = "127.0.0.1"
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ip = ipNet.IP.String()
				return
			}
		}
	}
	return
}
