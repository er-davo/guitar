package guitar

import (
	"fmt"
	"math"
	"strings"
)

type TabWriter struct {
	time     float64
	timeStep float64

	lines []strings.Builder
}

type Playable interface {
	TabSymbol() string
	StringNumber() int
	FretPosition() int

	StartTime() float64

	ScoreTo(Playable) float64

	IsValid() bool
}

type TabFrame []Playable

func (tf TabFrame) Tab(numStrings int) ([]string, error) {
	tabs := make([]string, numStrings)

	maxTabLen := 0

	for _, p := range tf {
		str := p.StringNumber()
		if str < 0 || str >= numStrings {
			return nil, fmt.Errorf("invalid string number %d (max %d)", str, numStrings-1)
		}
		symbol := p.TabSymbol()
		tabs[str] = symbol
		maxTabLen = max(maxTabLen, len(symbol))
	}

	for i := range tabs {
		if tabs[i] == "" {
			tabs[i] = strings.Repeat("-", maxTabLen)
		} else if len(tabs[i]) < maxTabLen {
			tabs[i] += strings.Repeat("-", maxTabLen-len(tabs[i]))
		}
	}

	return tabs, nil
}

func (tf TabFrame) Time() float64 {
	if len(tf) == 0 {
		return 0
	}
	return tf[0].StartTime()
}

func NewTabWriter(tuningNotes []string, opts ...TabOption) (*TabWriter, error) {
	tb := &TabWriter{
		time:     0,
		timeStep: 0.2,
		lines:    make([]strings.Builder, len(tuningNotes)),
	}

	for _, opt := range opts {
		opt(tb)
	}

	if err := tb.addNotes(tuningNotes); err != nil {
		return nil, err
	}

	return tb, nil
}

func (tb *TabWriter) Tab() string {
	tab := strings.Builder{}

	for i := range len(tb.lines) {
		tab.WriteString(tb.lines[i].String() + "\n")
	}

	return tab.String()
}

func (tb *TabWriter) Write(frames ...TabFrame) error {
	if len(frames) == 0 {
		return nil
	}

	for _, frame := range frames {
		frameTime := frame.Time()
		if frameTime < tb.time-tb.timeStep {
			return fmt.Errorf("frame time %f is earlier than current time %f", frameTime, tb.time)
		}

		// if frameTime < tb.time {
		// 	for i := range tb.lines {
		// 		tb.lines[i].WriteString("-")
		// 	}
		// }

		if frameTime > tb.time {
			padding := int((frameTime - tb.time) / tb.timeStep)
			for i := range tb.lines {
				tb.lines[i].WriteString(strings.Repeat("-", padding))
			}
		} else if tb.lines[0].Len() == 0 {
			for i := range tb.lines {
				tb.lines[i].WriteString("-")
			}
		}

		tabSymbols, err := frame.Tab(len(tb.lines))
		if err != nil {
			return err
		}

		for i := range tb.lines {
			tb.lines[i].WriteString(tabSymbols[i])
		}

		// tb.time = frameTime + tb.timeStep
		tb.time = frameTime + (tb.timeStep - math.Mod(frameTime, tb.timeStep))
	}

	return nil
}

func (tb *TabWriter) addNotes(notes []string) error {
	if len(notes) != len(tb.lines) {
		return fmt.Errorf("invalid tuning notes count")
	}

	for i := range notes {
		tb.lines[i].WriteString(notes[i] + "|")
	}
	return nil
}

type TabOption func(*TabWriter)

func WithTimeStep(step float64) TabOption {
	return func(tb *TabWriter) {
		tb.timeStep = step
	}
}

func WithDefaultTimeStep() TabOption {
	return func(tb *TabWriter) {
		tb.timeStep = 0.2
	}
}
