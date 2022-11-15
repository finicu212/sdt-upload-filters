package connection

import (
	"errors"
)

const (
	ConnectionsLimit int = 10
)

var (
	ErrConnectionLimit = errors.New("hit connection limit, please try again later")
)

type IPool interface {
	// GetConnection calls newConnection if no connections exist in the pool,
	// otherwise calls existingConnection which should return an already existing connection.
	GetConnection() (IConnection, error)
	TimeoutConnection()
	ReleaseConnection(IConnection) error

	// private methods
	newConnection() (IConnection, error)
	existingConnection() (IConnection, error)
}

type Pool struct {
	connections []IConnection
}

func (p Pool) GetConnection() (IConnection, error) {
	if len(p.connections) > 0 {
		return p.existingConnection(), nil
	}
	if len(p.connections) < ConnectionsLimit {
		return p.newConnection()
	}
	return nil, ErrConnectionLimit
}

func (p Pool) newConnection() (IConnection, error) {
	return NewConnection() // Call the constructor with some params
}

func (p Pool) existingConnection() IConnection {
	if len(p.connections) == 0 {
		panic("Not enough connections! How did this function get called!?")
	}
	var x IConnection
	x, p.connections = p.connections[0], p.connections[1:]
	//p.connections = cs
	return x
}
