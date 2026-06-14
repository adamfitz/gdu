package parser

import (
	"fmt"
	"os"
	"time"
)

type spinner struct {
	chars  []string
	delay  time.Duration
	pos    int
	active bool
	stop   chan struct{}
}

var defaultChars = []string{"|", "/", "-", "\\"}

func defaultSpinner(chars []string, delay time.Duration) *spinner {
	return &spinner{
		chars: chars,
		delay: delay,
		stop:  make(chan struct{}),
	}
}

// Start hides the cursor then launches a goroutine that writes the next
// spinner character to stdout on each tick, overwriting the same line
// with \r to produce the animation effect.
func (s *spinner) Start() {
	s.active = true
	fmt.Fprint(os.Stdout, "\033[?25l")
	go func() {
		for {
			select {
			case <-s.stop:
				return
			default:
				fmt.Fprintf(os.Stdout, "\r%s", s.chars[s.pos])
				s.pos = (s.pos + 1) % len(s.chars)
				time.Sleep(s.delay)
			}
		}
	}()
}

// Stop closes the stop channel to signal the animation goroutine to exit,
// restores the cursor, and clears the spinner line from the terminal.
func (s *spinner) Stop() {
	s.active = false
	close(s.stop)
	fmt.Fprint(os.Stdout, "\033[?25h")
	fmt.Fprint(os.Stdout, "\r")
}
