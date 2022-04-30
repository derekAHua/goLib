package utils

import (
	"net"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetLocalIp() string {
	addrS, _ := net.InterfaceAddrs()
	var ip string
	for _, addr := range addrS {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ip = ipNet.IP.String()
				if ip != "127.0.0.1" {
					return ip
				}
			}
		}
	}
	return "127.0.0.1"
}

func GetClientIp(ctx *gin.Context) string {
	clientIP := ctx.GetHeader("X-Original-Forwarded-For")
	clientIP = strings.TrimSpace(strings.Split(clientIP, ",")[0])
	if clientIP != "" {
		return clientIP
	}
	return ctx.ClientIP()
}
