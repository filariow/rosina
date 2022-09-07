package rpin

import (
	"github.com/warthog618/gpiod"
)

type OutPin interface {
	High()
	Low()
}

func New(number uint8) (OutPin, error) {
	p, err := gpiod.RequestLine("gpiochip0", int(number), gpiod.AsOutput(0))
	if err != nil {
		return nil, err
	}

	return &pin{
		number: number,
		out:    p,
	}, nil
}

type pin struct {
	number uint8
	out    *gpiod.Line
}

func (p *pin) High() {
	p.out.SetValue(1)
}

func (p *pin) Low() {
	p.out.SetValue(0)
}
