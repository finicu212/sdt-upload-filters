package handler

import (
	"io"
	"log"
	"time"
)

type LoggedReverser struct {
	Reverser
}

var _ IHandler = new(LoggedReverser)
var _ Reverser = new(LoggedReverser).Reverser

func (h LoggedReverser) Handle(writer io.Writer, reader io.Reader) {
	log.Printf("Using Reverser handler!\n")
	h.Reverser.Handle(writer, reader)
	log.Printf("Done using Reverser handler!\n")
}

type TimedSkipper struct {
	Skipper
}

var _ IHandler = new(TimedSkipper)
var _ Skipper = new(TimedSkipper).Skipper

func (h TimedSkipper) Handle(writer io.Writer, reader io.Reader) {
	start := time.Now()
	h.Skipper.Handle(writer, reader)
	log.Printf("Time elapsed for skipping headers: %s\n", time.Since(start))
}
