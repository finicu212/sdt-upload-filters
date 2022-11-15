package connection

import (
	"github.com/google/uuid"
	"io"
)

type ConnectionType uint8

const (
	FTP ConnectionType = 0
)

type IConnection interface {
	Store(reader io.Reader)
}

type Connection struct {
	UUID string
}

func (c Connection) Store(reader io.Reader) {
	//TODO implement me
	panic("implement me")
}

// NewConnection instanciates a new Connection
func NewConnection() (IConnection, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return Connection{
		UUID: id.String(),
	}, nil
}
