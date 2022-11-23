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
	connections       []connection.IConnection
	activeConnections int
}

var _ IPool = new(Pool) // Must fail to compile if Pool doesn't implement IPool

var (
	instance IPool
	lock     = &sync.Mutex{}
)

func newPool() IPool {
	p := new(Pool)
	p.connections = make([]connection.IConnection, 0)
	return p
}

func Instance() IPool {
	lock.Lock()
	defer lock.Unlock()

	if instance == nil {
		// 1. We can put concrete class pointers in interfaces.
		// 2. Thread safe, since we use a mutex lock
		fmt.Printf("New instance created\n")
		instance = newPool()
	}

	return instance
}

func (p *Pool) DropConnection() {
	//TODO Connection timeout if not used for x seconds
	panic("implement me")
}

func (p *Pool) ReleaseConnection(connection connection.IConnection) error {
	lock.Lock()
	defer lock.Unlock()
	p.connections = append(p.connections, connection)
	return nil
}

func (p *Pool) GetConnection(ctx context.Context) (connection.IConnection, error) {
	lock.Lock()
	defer lock.Unlock()
	if len(p.connections) > 0 {
		fmt.Printf("Existing connections: %d/%d\n", len(p.connections), ConnectionsLimit)
		return p.existingConnection(), nil
	}
	if p.activeConnections < ConnectionsLimit {
		conn := connection.NewConnection(ctx)
		p.activeConnections += 1
		fmt.Printf("New connection: %s! (%d/%d)\n", conn.GetUUID(), p.activeConnections, ConnectionsLimit)

		return conn, nil
	}
	return nil, ErrConnectionLimit
}

func (p *Pool) existingConnection() connection.IConnection {
	if len(p.connections) == 0 {
		panic("Not enough connections! How did this function get called!?")
	}
	if len(p.connections) == 1 {
		x := p.connections[0]
		p.connections = nil
		return x
	}
	var x connection.IConnection
	x, p.connections = p.connections[0], p.connections[1:]
	return x
}
