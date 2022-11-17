package utils

import (
	"net"
	"strconv"
)

func SplitHostPort(hostport string) (string, int, error) {
	url, port, err := net.SplitHostPort(hostport)
	if err != nil {
		return "", 0, err
	}
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return "", 0, err
	}
	return url, portInt, nil
}

func Repeated[T any](val T, times int) []T {
	arr := make([]T, times)
	for i := 0; i < times; i++ {
		arr[i] = val
	}
	return arr
}
