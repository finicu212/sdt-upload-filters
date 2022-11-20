package handler

import (
	"bytes"
	"github.com/kteru/reversereader"
	"io"
	"io/ioutil"
)

type IHandler interface {
	Handle(writer io.Writer, reader io.Reader)
	SetNext(handler IHandler) IHandler
}

func stringToHandler(chain string) IHandler {
	switch chain {
	case "reverser":
		return Reverser{}
	case "duplicator":
		return Duplicator{}
	case "default":
	}
	panic("wrong handler type: " + chain)
}

// StringSliceToChain returns the fully built chain's root, given an array of handler names (i.e. "[Reverser Duplicator]")
func StringSliceToChain(chain []string) IHandler {
	root, chain := stringToHandler(chain[0]), chain[1:]

	prev := root
	for _, handler := range chain {
		next := stringToHandler(handler)
		prev.SetNext(next)
		prev = next
	}
	return root
}

// -----------------------
// ---- Concrete handlers
// ----

type Reverser struct {
	next IHandler
}

func (h Reverser) Handle(writer io.Writer, reader io.Reader) {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	readSeeker := bytes.NewReader(b)
	rvrd := reversereader.NewReader(readSeeker)
	_, err = io.Copy(writer, rvrd)
	if err != nil {
		panic(err)
	}
	if h.next != nil {
		h.next.Handle(writer, reader)
	}
}

func (h Reverser) SetNext(handler IHandler) IHandler {
	h.next = handler
	return h
}

type Duplicator struct {
	next IHandler
}

func (h Duplicator) Handle(writer io.Writer, reader io.Reader) {
	io.Copy(writer, reader)
	io.Copy(writer, reader)
}

func (h Duplicator) SetNext(handler IHandler) IHandler {
	h.next = handler
	return h
}
