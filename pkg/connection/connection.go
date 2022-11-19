package connection

import (
	"context"
	"crypto/tls"
	"github.com/google/uuid"
	"github.com/kardianos/ftps"
	"io"
	"path"
)

// TODO: Custom Connection Types
//type Type uint8
//
//const (
//	FTP Type = 0
//)

type IConnection interface {
	Store(filename string, reader io.Reader) error
	GetUUID() string
}

type FTPConnection struct {
	UUID   string
	client *ftps.Client
	ctx    context.Context
}

func (c FTPConnection) Store(filepath string, reader io.Reader) error {
	filename := path.Base(filepath)
	if filename != filepath {
		// If filepath consists of a directory structure
		err := c.client.Chdir(path.Dir(filepath))
		if err != nil {
			return err
		}
	}
	return c.client.Upload(c.ctx, filename, reader)
}

func (c FTPConnection) GetUUID() string {
	return c.UUID
}

// NewFTPConnection instantiates a new Connection
func NewFTPConnection(ctx context.Context, username, password, url string, port int) (IConnection, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	c, err := ftps.Dial(ctx, ftps.DialOptions{
		Username:    username,
		Passowrd:    password,
		Host:        url,
		Port:        port,
		ExplicitTLS: true, // TODO: If you are forking this repo, review these critical security settings!
		TLSConfig: &tls.Config{
			ServerName:         url,
			InsecureSkipVerify: false,            // TODO: If you are forking this repo, review these critical security settings!
			MaxVersion:         tls.VersionTLS12, // TODO: If you are forking this repo, review these critical security settings!
		},
	})
	if err != nil {
		return nil, err
	}

	return FTPConnection{
		UUID:   id.String(),
		ctx:    ctx,
		client: c,
	}, nil
}

// Mocks. To be auto gen

type MockConnection struct {
	UUID string
}

func (m MockConnection) GetUUID() string {
	return m.UUID
}

func (m MockConnection) Store(_ string, _ io.Reader) error {
	panic("this is just a mock!")
}

func NewMockConnection() IConnection {
	id, _ := uuid.NewRandom()
	return MockConnection{
		UUID: id.String(),
	}
}
