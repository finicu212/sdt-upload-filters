package pool

import (
	"context"
	"errors"
	"fmt"
	"sdt-upload-filters/pkg/connection"
	"sync"
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
	GetConnection(ctx context.Context) (connection.IConnection, error)
	ReleaseConnection(connection.IConnection) error
	// private methods
	existingConnection() connection.IConnection
}

type Pool struct {
	connections []connection.IConnection
}

var (
	instance IPool
	lock     = &sync.Mutex{}
)

func Instance() IPool {
	lock.Lock()
	defer lock.Unlock()

	if instance == nil {
		// 1. We can put concrete class pointers in interfaces.
		// 2. Thread safe, since we use a mutex lock
		fmt.Printf("New instance created\n")
		instance = new(Pool)
	}

	return instance
}

func (p Pool) DropConnection() {
	//TODO implement me
	panic("implement me")
}

func (p Pool) ReleaseConnection(connection connection.IConnection) error {
	//TODO implement me
	panic("implement me")
}

func (p Pool) GetConnection(ctx context.Context) (connection.IConnection, error) {
	lock.Lock()
	defer lock.Unlock()
	if len(p.connections) > 0 {
		fmt.Printf("Existing connections: %d/%d\n", len(p.connections), ConnectionsLimit)
		return p.existingConnection(), nil
	}
	if len(p.connections) < ConnectionsLimit {
		conn := connection.NewConnection(ctx)
		p.connections = append(p.connections, conn)
		fmt.Printf("New connection: %s! (%d/%d)\n", conn.GetUUID(), len(p.connections), ConnectionsLimit)

		return conn, nil
	}
	return nil, ErrConnectionLimit
}

func (p Pool) existingConnection() connection.IConnection {
	lock.Lock()
	defer lock.Unlock()
	if len(p.connections) == 0 {
		panic("Not enough connections! How did this function get called!?")
	}
	var x connection.IConnection
	x, p.connections = p.connections[0], p.connections[1:]
	fmt.Printf("Using existing connection: %s! (%d/%d)\n", x.GetUUID(), len(p.connections), ConnectionsLimit)
	return x
}
