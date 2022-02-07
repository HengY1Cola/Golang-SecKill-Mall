package common

import (
	"errors"
	"net"
)

// GetEntranceIp 拿到本机IP
func GetEntranceIp() (string, error) {
	addrLiat, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, address := range addrLiat {
		//检查Ip地址判断是否回环地址
		if ip, ok := address.(*net.IPNet); ok && !ip.IP.IsLoopback() {
			if ip.IP.To4() != nil {
				return ip.IP.String(), nil
			}
		}
	}
	return "", errors.New("获取地址异常")
}
