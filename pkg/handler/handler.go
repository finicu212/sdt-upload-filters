package handler

import (
	"bytes"
	"github.com/kteru/reversereader"
	"io"
)

type IHandler interface {
	Handle(writer io.Writer, reader io.Reader)
	SetNext(handler IHandler)
}

type Reverser struct {
	next IHandler
}

func (h Reverser) Handle(writer io.Writer, reader io.Reader) {
	b := make([]byte, 1024)
	_, err := reader.Read(b)
	if err != nil {
		return
	}
	readSeeker := bytes.NewReader(b)
	rvrd := reversereader.NewReader(readSeeker)
	_, err = io.Copy(writer, rvrd)
	if err != nil {
		return
	}
	if h.next != nil {
		h.next.Handle(writer, reader)
	}
}

func (h Reverser) SetNext(handler IHandler) {
	h.next = handler
}
