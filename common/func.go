package common

import (
	"net"
	"strings"
)

var LocalIp string
var RemoteIp string

func FetchNetIp() {
	conn, err := net.Dial("udp", "google.com:80")
	if err == nil {
		defer conn.Close()
		LocalIp = strings.Split(conn.LocalAddr().String(), ":")[0]
		RemoteIp = strings.Split(conn.LocalAddr().String(), ":")[0]
	} else {
		LocalIPv4s()
	}
}

func LocalIPv4s() {
	var ips []string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			ips = append(ips, ipnet.IP.String())
		}
	}
	if len(ips) > 0 {
		LocalIp = strings.Split(ips[0], ":")[0]
	}
}
