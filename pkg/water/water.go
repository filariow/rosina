package water

import (
	"time"

	"github.com/filariow/rosina/internal/rpin"
)

type Waterer interface {
	Open()
	Close()
}

func New(outpin1, outpin2 rpin.OutPin) Waterer {
	return &waterer{
		pin1: outpin1,
		pin2: outpin2,
	}
}

type waterer struct {
	pin1 rpin.OutPin
	pin2 rpin.OutPin
}

func (w *waterer) waitAndReset() {
	time.Sleep(3000 * time.Millisecond)

	w.pin1.Low()
	w.pin2.Low()
}

func (w *waterer) Open() {
	defer w.waitAndReset()

	w.pin1.High()
}

func (w *waterer) Close() {
	defer w.waitAndReset()

	w.pin2.High()
}
