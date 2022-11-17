package connection

import (
	"reflect"
	"testing"
)

func TestPool_existingConnection(t *testing.T) {
	conn1 := NewMockConnection()
	conn2 := NewMockConnection()
	type fields struct {
		connections []IConnection
	}
	tests := []struct {
		name      string
		fields    fields
		want      IConnection
		wantPanic bool
	}{
		{
			name: "Only one existing connection",
			fields: fields{
				connections: []IConnection{conn1},
			},
			want: conn1,
		},
		{
			name: "Three connections",
			fields: fields{
				connections: []IConnection{conn2, conn1, conn1, conn1},
			},
			want: conn2,
		},
		{
			name: "No connections, expect panic",
			fields: fields{
				connections: []IConnection{},
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
