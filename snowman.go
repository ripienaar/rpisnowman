package main

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	rpio "github.com/stianeikeland/go-rpio"
)

// SnowMan is a Ryanteck RTK-000-00A GPIO Snowman
type SnowMan struct {
	name    string
	leds    [9]rpio.Pin
	pins    []int
	paused  bool
	schemes []Scheme
	log     *logrus.Entry
}

// NewSnowMan sets up a new SnowMan interface
func NewSnowMan(name string, log *logrus.Entry) *SnowMan {
	return &SnowMan{
		name: name,
		leds: [9]rpio.Pin{},
		pins: []int{7, 8, 9, 17, 18, 22, 23, 24, 25},
		log:  log,
		schemes: []Scheme{
			NewSimpleScheme(),
			NewCrossScheme(),
			NewLinesScheme(),
			NewUpScheme(),
			NewDownScheme(),
		},
	}
}

func (s *SnowMan) Log() *logrus.Entry {
	return s.log
}

// Open initializes the snowman and set all pins to output
func (s *SnowMan) Open() error {
	err := rpio.Open()
	if err != nil {
		return err
	}

	for i, pin := range s.pins {
		s.leds[i] = rpio.Pin(pin)
		s.leds[i].Output()
	}

	return nil
}

// Close closes down the GPIO system
func (s *SnowMan) Close() {
	rpio.Close()
}

// Toggle flips a specific leds to their opposite
func (s *SnowMan) Toggle(leds ...int) {
	for _, i := range leds {
		s.leds[i].Toggle()
	}
}

// Off turns all given led lines low
func (s *SnowMan) Off(leds ...int) {
	for _, i := range leds {
		s.leds[i].Low()
	}
}

// On turns all given led lines high
func (s *SnowMan) On(leds ...int) {
	for _, i := range leds {
		s.leds[i].High()
	}
}

// ToggleAll toggles all leds to their opposite
func (s *SnowMan) ToggleAll() {
	for _, led := range s.leds {
		led.Toggle()
	}
}

// Flash toggles the given leds sleeps a bit and toggles them again
func (s *SnowMan) Flash(p ...int) {
	s.Toggle(p...)

	sleep()

	s.Toggle(p...)
}

// AllLow sets all led lines low
func (s *SnowMan) AllLow() {
	for _, led := range s.leds {
		led.Low()
	}
}

// Run runs the lights until interrupt
func (s *SnowMan) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	r := rand.New(rand.NewSource(time.Now().Unix()))
	s.AllLow()

	for {
		for _, i := range r.Perm(len(s.schemes)) {
			if ctx.Err() != nil {
				return
			}

			if !s.paused {
				s.schemes[i].Run(s)
				s.AllLow()
			}

			if ctx.Err() != nil {
				return
			}
			sleep()
		}
	}
}
