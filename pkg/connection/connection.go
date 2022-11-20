package connection

import (
	"context"
	"github.com/google/uuid"
)

type IConnection interface {
	GetUUID() string
}

type MockConnection struct {
	UUID string
}

func (m MockConnection) GetUUID() string {
	return m.UUID
}

func NewConnection(ctx context.Context) IConnection {
	return MockConnection{
		UUID: uuid.NewString(),
	}
}
