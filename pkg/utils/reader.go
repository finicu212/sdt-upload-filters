package utils

import (
	"bufio"
	"io"
	"os"
)

func GetFileReader(path string) (io.Reader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	buf := bufio.NewReader(file)
	return buf, nil
}

func PathsAsReaders(paths []string) ([]io.Reader, error) {
	var readers []io.Reader
	for _, p := range paths {
		r, err := GetFileReader(p)
		if err != nil {
			return nil, err
		}
		readers = append(readers, r)
	}
	return readers, nil
}
