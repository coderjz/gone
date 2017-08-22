package util

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/coderjz/gone/painter"
	"github.com/coderjz/gone/state"
	"github.com/deckarep/gosx-notifier"
)

// Timer count time and update the state accordingly.
type Timer struct {
	state   *state.State
	command string
	ticker  *time.Ticker
	painter *painter.Painter
}

// NewTimer create a new timer using a state, a command to execute
// and a painter to draw the screnn
func NewTimer(s *state.State, p *painter.Painter, c string) *Timer {
	return &Timer{
		state:   s,
		painter: p,
		command: c,
	}
}

// run launch a timer and write the counter using the writer
func (t *Timer) Run() {
	//start a new timer
	t.ticker = time.NewTicker(250 * time.Millisecond)
	i := 1
	for _ = range t.ticker.C {
		t.painter.Draw()
		if i > 4 && t.state.IsRunning() {
			i = 1
			t.state.Decrease()
			if t.state.IsEnded() {
				break
			}
		} else {

			i++
		}

	}
	t.ticker.Stop()
	if t.command != "" {
		v := strings.Split(t.command, " ")
		cmd := exec.Command(v[0], v[1:]...)
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Fail to execute command : %s - %v\n", t.command, err)
		}
	}

	//t.state.WaitForConfirm(t.Run)
	//Auto-continue - switch to next state
	t.state.Next()

	//Show notification for our new state
	note := gosxnotifier.NewNotification("")

	note.Message = fmt.Sprintf("%s (%s) - Until %v", t.state.Message(), t.state.Duration(),
		time.Now().Add(time.Second*time.Duration(t.state.DurationInSecs())).Format(time.Kitchen))
	note.Title = "Pomodoro"
	err := note.Push()
	if err != nil {
		fmt.Printf("Fail to show notification: %v\n", err)
	}

	//Auto-continue - run the new state
	t.Run()
}

// Stop the timer
func (t *Timer) Stop() {
	if t.ticker != nil {
		t.ticker.Stop()
	}
}
