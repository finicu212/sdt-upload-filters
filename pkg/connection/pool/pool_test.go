package pool

import (
	"context"
	"reflect"
	"sdt-upload-filters/pkg/connection"
	"testing"
)

func TestPool_existingConnection(t *testing.T) {
	conn1 := connection.NewConnection(context.TODO())
	conn2 := connection.NewConnection(context.TODO())
	type fields struct {
		connections []connection.IConnection
	}
	tests := []struct {
		name      string
		fields    fields
		want      connection.IConnection
		wantPanic bool
	}{
		{
			name: "Only one existing connection",
			fields: fields{
				connections: []connection.IConnection{conn1},
			},
			want: conn1,
		},
		{
			name: "Three connections",
			fields: fields{
				connections: []connection.IConnection{conn2, conn1, conn1, conn1},
			},
			want: conn2,
		},
		{
			name: "No connections, expect panic",
			fields: fields{
				connections: []connection.IConnection{},
			},
			want:      nil,
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Pool{
				connections: tt.fields.connections,
			}
			defer func() {
				if r := recover(); r != nil && !tt.wantPanic {
					t.Errorf("Panicked when no panic expected")
				}
			}()
			got := p.existingConnection()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("existingConnection() got = %v, want %v", got, tt.want)
			}
		})
	}
}
