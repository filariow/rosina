package rpin

import "github.com/stianeikeland/go-rpio/v4"

type OutPin interface {
	High()
	Low()
}

func New(number uint8) OutPin {
	return &pin{number: number, out: rpio.Pin(number)}
}

type pin struct {
	number uint8
	out    rpio.Pin
}

func (p *pin) High() {
	p.out.High()
}

func (p *pin) Low() {
	p.out.Low()
}
