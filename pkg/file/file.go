package file

import (
	"bufio"
	"io"
	"os"
)

// FileDetails is a logical grouping of a buffer to a file, it's local path, and desired path on the remote, after upload
type FileDetails struct {
	DataReader io.Reader
	LocalPath  string
	RemotePath string
}

func getFileReader(path string) (io.Reader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	buf := bufio.NewReader(file)
	return buf, nil
}

func NewFileDetails(path string) (fd FileDetails, err error) {
	fd = *new(FileDetails)
	rd, err := getFileReader(path)
	if err != nil {
		return fd, err
	}
	fd.DataReader = rd
	fd.LocalPath = path
	fd.RemotePath = path // placeholder
	return fd, nil
}

func ManyFileDetails(paths []string) (fds []FileDetails, err error) {
	for _, p := range paths {
		fd, err := NewFileDetails(p)
		if err != nil {
			return nil, err
		}
		fds = append(fds, fd)
	}
	return fds, nil
}
