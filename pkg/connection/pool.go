package connection

import (
	"context"
	"errors"
	"net"
	"sdt-upload-filters/pkg/utils"
)

const (
	ConnectionsLimit int = 10
)

var (
	ErrConnectionLimit = errors.New("hit connection limit, please try again later")
)

func ConnectToIPs(ips []net.IP) []IPool {
	for _, ip := range ips {
		url, port, err := utils.SplitHostPort(ip)
		if err != nil {
			return err
		}
		pool := connection.NewPool(url, port)
		c, err := pool.GetConnection("root", "root")
		if err != nil {
			c = connection.NewMockConnection()
		}
		conns = append(conns, c)
	}
}

type IPool interface {
	// GetConnection calls newConnection if no connections exist in the pool,
	// otherwise calls existingConnection which should return an already existing connection.
	GetConnection(username string, password string) (IConnection, error)
	//TimeoutConnection()
	ReleaseConnection(IConnection) error

	// private methods
	existingConnection() IConnection
}

type Pool struct {
	url         string
	port        int
	connections []IConnection
}

func NewPool(url string, port int) IPool {
	return Pool{url: url, port: port}
}

func (p Pool) DropConnection() {
	//TODO implement me
	panic("implement me")
}

func (p Pool) ReleaseConnection(connection IConnection) error {
	//TODO implement me
	panic("implement me")
}

func (p Pool) GetConnection(username, password string) (IConnection, error) {
	if len(p.connections) > 0 {
		return p.existingConnection(), nil
	}
	if len(p.connections) < ConnectionsLimit {
		return NewFTPConnection(context.Background(), username, password, p.url, p.port)
	}
	return nil, ErrConnectionLimit
}

func (p Pool) existingConnection() IConnection {
	if len(p.connections) == 0 {
		panic("Not enough connections! How did this function get called!?")
	}
	var x IConnection
	x, p.connections = p.connections[0], p.connections[1:]
	return x
}
