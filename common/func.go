package common

import (
	"fmt"
	"net"
	"strings"
)

var LocalIp string
var RemoteIp string

func FetchNetIp() {
	conn, err := net.Dial("udp", "google.com:80")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer conn.Close()
	LocalIp = strings.Split(conn.LocalAddr().String(), ":")[0]
	RemoteIp = strings.Split(conn.LocalAddr().String(), ":")[0]
}
