package utils

import (
	"net"
	"strconv"
)

func SplitHostPort(hostport net.IP) (string, int, error) {
	url, port, err := net.SplitHostPort(string(hostport))
	if err != nil {
		return "", 0, err
	}
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return "", 0, err
	}
	return url, portInt, nil
}
