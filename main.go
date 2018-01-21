package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	rpio "github.com/stianeikeland/go-rpio"
)

var (
	leds   = [9]rpio.Pin{}
	pins   = []int{7, 8, 9, 17, 18, 22, 23, 24, 25}
	paused = false
	err    error
)

func sleep() {
	time.Sleep(1 * time.Second)
}

func flip(p ...int) {
	for _, i := range p {
		leds[i].Toggle()
	}

	sleep()

	for _, i := range p {
		leds[i].Toggle()
	}
}

func simpleScheme() {
	seq := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 8, 7, 6, 5, 4, 3, 2, 1, 0}
	for _, pin := range seq {
		leds[pin].Toggle()
		sleep()
	}
	allToggle()
}

func crossScheme() {
	flip(0, 5)
	flip(1, 4)
	flip(2, 3)
	flip(6, 7)
	flip(0, 1, 2, 3, 4, 5, 6, 7, 8)
}

func linesScheme() {
	flip(0, 1, 2)
	flip(3, 4, 5)
	flip(6, 7, 8)
}

func upScheme() {
	flip(2, 5)
	flip(1, 4)
	flip(0, 3)
	flip(8)
	flip(6, 7)
}

func downScheme() {
	flip(6, 7)
	flip(8)
	flip(0, 3)
	flip(1, 4)
	flip(2, 5)
}

func allScheme() {
	allToggle()
	sleep()
	allToggle()
}

func allToggle() {
	for _, led := range leds {
		led.Toggle()
	}
}

func allOff() {
	for _, led := range leds {
		led.Low()
	}
}

func main() {
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer rpio.Close()

	if enableChoria {
		go runChoria()
	}

	for i, pin := range pins {
		leds[i] = rpio.Pin(pin)
		leds[i].Output()
	}

	schemes := []func(){upScheme, crossScheme, linesScheme, allScheme, downScheme, simpleScheme}
	r := rand.New(rand.NewSource(time.Now().Unix()))

	for {
		for _, i := range r.Perm(len(schemes)) {
			allOff()
			if !paused {
				schemes[i]()
			}
			sleep()
		}
	}
}
