package guitar

import (
	"fmt"
)

type Harmonic struct {
	Note
}

// override
func (h Harmonic) TabSymbol() string {
	return fmt.Sprintf("<%d>", h.Fret)
}

// func (h Harmonic) StringNumber() int {
// 	return h.String
// }

// func (h Harmonic) FretPosition() int {
// 	return h.Fret
// }

// func (h Harmonic) StartTime() float64 {
// 	return h.Time
// }

// func (h Harmonic) ScoreTo(to Playable) float64 {
// 	return h.Note.ScoreTo(to)
// }

// func (h Harmonic) IsValid() bool {
// 	return h.MidiPitch >= 0 && h.MidiPitch <= 127
// }

type Slide struct {
	NoteFrom Note
	NoteTo   Note
}

func (s Slide) TabSymbol() string {
	if s.NoteFrom.Fret == -1 {
		return fmt.Sprintf("/%d", s.NoteTo.Fret)
	}
	return fmt.Sprintf("%d/%d", s.NoteFrom.Fret, s.NoteTo.Fret)
}

func (s Slide) StringNumber() int {
	return s.NoteFrom.String
}

func (s Slide) FretPosition() int {
	return s.NoteFrom.Fret
}

func (s Slide) StartTime() float64 {
	return s.NoteFrom.Time
}

func (s Slide) ScoreTo(to Playable) float64 {
	return s.NoteTo.ScoreTo(to)
}

func (s Slide) IsValid() bool {
	return (s.NoteFrom.MidiPitch >= 0 && s.NoteFrom.MidiPitch <= 127) ||
		(s.NoteTo.MidiPitch >= 0 && s.NoteTo.MidiPitch <= 127)
}

type HammerOn struct {
	NoteFrom Note
	NoteTo   Note
}

func (h HammerOn) TabSymbol() string {
	return fmt.Sprintf("%dh%d", h.NoteFrom.Fret, h.NoteTo.Fret)
}

func (h HammerOn) StringNumber() int {
	return h.NoteFrom.String
}

func (h HammerOn) FretPosition() int {
	return h.NoteFrom.Fret
}

func (h HammerOn) StartTime() float64 {
	return h.NoteFrom.Time
}

func (h HammerOn) ScoreTo(to Playable) float64 {
	return h.NoteTo.ScoreTo(to)
}

func (h HammerOn) IsValid() bool {
	return (h.NoteFrom.MidiPitch >= 0 && h.NoteFrom.MidiPitch <= 127) ||
		(h.NoteTo.MidiPitch >= 0 && h.NoteTo.MidiPitch <= 127)
}

type PullOff struct {
	NoteFrom Note
	NoteTo   Note
}

func (p PullOff) TabSymbol() string {
	return fmt.Sprintf("%dp%d", p.NoteFrom.Fret, p.NoteTo.Fret)
}

func (p PullOff) StringNumber() int {
	return p.NoteFrom.String
}

func (p PullOff) FretPosition() int {
	return p.NoteFrom.Fret
}

func (p PullOff) StartTime() float64 {
	return p.NoteFrom.Time
}

func (p PullOff) ScoreTo(to Playable) float64 {
	return p.NoteTo.ScoreTo(to)
}

func (p PullOff) IsValid() bool {
	return (p.NoteFrom.MidiPitch >= 0 && p.NoteFrom.MidiPitch <= 127) ||
		(p.NoteTo.MidiPitch >= 0 && p.NoteTo.MidiPitch <= 127)
}
