package utils

import (
	"fmt"
	"net"
	"time"
)

//tcp 连接是否正常监测
func TcpStatusCheck(address string) bool {

	conn, err := net.DialTimeout("tcp", address, 5 * time.Second)

	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()

	if err != nil {
		fmt.Printf("Fatal error: %s \n", err.Error())
		return false
	}


	return true

}
