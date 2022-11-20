package handler

import (
	"reflect"
	"testing"
)

func TestStringSliceToChain(t *testing.T) {
	r := func() IHandler {
		return Reverser{}
	}
	rr := func() IHandler {
		rev := r()
		rev.SetNext(r())
		return rev
	}
	rrs := func() IHandler {
		rev := rr()
		rev.SetNext(Skipper{})
		return rev
	}

	tests := []struct {
		name string
		args []string
		want IHandler
	}{
		{
			name: "reverser -> reverser",
			args: []string{"reverser", "reverser"},
			want: rr(),
		},
		{
			name: "reverser",
			args: []string{"reverser"},
			want: r(),
		},
		{
			name: "reverser -> reverser -> skipper",
			args: []string{"reverser", "reverser", "skipper"},
			want: rrs(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringSliceToChain(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringSliceToChain() = %v, want %v", got, tt.want)
			}
		})
	}
}
