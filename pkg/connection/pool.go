package connection

import (
	"context"
	"errors"
	"fmt"
	"sdt-upload-filters/pkg/file"
)

const (
	// ConnectionsLimit is used for every pool. Each pool can have its independent limit, but we use this to keep things simple
	ConnectionsLimit int = 2
)

var (
	ErrConnectionLimit = errors.New("hit connection limit, please try again later")
)

type IPool interface {
	// GetConnection creates new connection via constructor if no connections exist in the pool,
	// otherwise calls existingConnection which should return an already existing connection.
	GetConnection(ctx context.Context) (IConnection, error)
	ReleaseConnection(IConnection) error
	Q() *file.FileQueue

	// private methods
	existingConnection() IConnection
}

type Pool struct {
	queue *file.FileQueue // Queue of buffers which still need to be uploaded to this pool's server

	user        string
	pass        string
	url         string
	port        int
	connections []IConnection
}

func NewPool(user, pass, url string, port int) IPool {
	return Pool{url: url, port: port, user: user, pass: pass, queue: new(file.FileQueue)}
}

func (p Pool) DropConnection() {
	//TODO implement me
	panic("implement me")
}

func (p Pool) ReleaseConnection(connection IConnection) error {
	//TODO implement me
	panic("implement me")
}

func (p Pool) GetConnection(ctx context.Context) (IConnection, error) {
	if len(p.connections) > 0 {
		return p.existingConnection(), nil
	}
	if len(p.connections) < ConnectionsLimit {
		fmt.Printf("New FTP connection: %s on %s\n", p.user, p.url)
		return NewFTPConnection(ctx, p.user, p.pass, p.url, p.port)
	}
	return nil, ErrConnectionLimit
}

func (p Pool) Q() *file.FileQueue {
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
