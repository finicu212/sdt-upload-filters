package connection

import (
	"context"
	"errors"
	"sdt-upload-filters/pkg/file"
)

const (
	// ConnectionsLimit is used for every pool. Each pool can have its independent limit, but we use this to keep things simple
	ConnectionsLimit int = 10
)

var (
	ErrConnectionLimit = errors.New("hit connection limit, please try again later")
)

type IPool interface {
	// GetConnection creates new connection via constructor if no connections exist in the pool,
	// otherwise calls existingConnection which should return an already existing connection.
	GetConnection() (IConnection, error)
	ReleaseConnection(IConnection) error
	Q() file.IFileQueue

	// private methods
	existingConnection() IConnection
}

type Pool struct {
	queue file.IFileQueue // Queue of buffers which still need to be uploaded to this pool's server

	user        string
	pass        string
	url         string
	port        int
	connections []IConnection
}

func NewPool(user, pass, url string, port int) IPool {
	return Pool{url: url, port: port, user: user, pass: pass}
}

func (p Pool) DropConnection() {
	//TODO implement me
	panic("implement me")
}

func (p Pool) ReleaseConnection(connection IConnection) error {
	//TODO implement me
	panic("implement me")
}

func (p Pool) GetConnection() (IConnection, error) {
	if len(p.connections) > 0 {
		return p.existingConnection(), nil
	}
	if len(p.connections) < ConnectionsLimit {
		return NewFTPConnection(context.Background(), p.user, p.pass, p.url, p.port)
	}
	return nil, ErrConnectionLimit
}

func (p Pool) Q() file.IFileQueue {
	return p.queue
}

func (p Pool) existingConnection() IConnection {
	if len(p.connections) == 0 {
		panic("Not enough connections! How did this function get called!?")
	}
	var x IConnection
	x, p.connections = p.connections[0], p.connections[1:]
	return x
}
