package water

import "github.com/filariow/rosina/internal/rpin"

type Waterer interface {
	Open()
	Close()
}

func New(outpin rpin.OutPin) Waterer {
	return &waterer{pin: outpin}
}

type waterer struct {
	pin rpin.OutPin
}

func (w *waterer) Open() {
	w.pin.High()
}

func (w *waterer) Close() {
	w.pin.Low()
}
