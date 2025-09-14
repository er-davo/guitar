package guitar

import (
	"errors"
	"fmt"
	"math"
)

const (
	stringlWeight   = 1.0  // Предпочтение вертикальным перемещениям (по струнам)
	fretWeight      = 1.0  // Вес для перемещений по ладам
	openStringBonus = -2.0 // Бонус за открытые струны
)

var notesChromo = []string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}

func MidiToNote(pitch int) (string, int) {
	if pitch < 0 || pitch > 127 {
		return "Invalid pitch", -1
	}

	note := notesChromo[pitch%12]
	octave := pitch/12 - 1

	return note, octave
}

func NoteToMidi(note string, octave int) int {
	if err := validateNote(&note); err != nil {
		return -1
	}
	for i, n := range notesChromo {
		if n == note {
			return (octave+1)*12 + i
		}
	}
	return -1
}

func validateNote(noteName *string) error {
	switch *noteName {
	case "Db", "D♭":
		*noteName = "C#"
	case "Eb", "E♭":
		*noteName = "D#"
	case "Gb", "G♭":
		*noteName = "F#"
	case "Ab", "A♭":
		*noteName = "G#"
	case "Bb", "B♭":
		*noteName = "A#"

	case "D♯":
		*noteName = "C#"
	case "E♯":
		*noteName = "D#"
	case "G♯":
		*noteName = "F#"
	case "A♯":
		*noteName = "G#"
	case "B♯":
		*noteName = "A#"

	case "C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B":
		return nil

	default:
		return fmt.Errorf("invalid note name: %s", *noteName)
	}
	return nil
}

type Note struct {
	MidiPitch int

	Fret   int
	String int

	Time float64
}

func (n Note) Name() string {
	return notesChromo[n.MidiPitch%12]
}

func (n Note) Octave() int {
	return n.MidiPitch/12 - 1
}

func (n Note) TabSymbol() string {
	return fmt.Sprintf("%d", n.Fret)
}

func (n Note) StringNumber() int {
	return n.String
}

func (n Note) FretPosition() int {
	return n.Fret
}

func (n Note) StartTime() float64 {
	return n.Time
}

func (n Note) IsValid() bool {
	return n.MidiPitch >= 0 && n.MidiPitch <= 127
}

func (n *Note) AddFret() error {
	if n.MidiPitch < 0 || n.MidiPitch > 127 {
		return fmt.Errorf("invalid pitch: %d", n.MidiPitch)
	}

	if n.Fret < 0 {
		return fmt.Errorf("invalid fret: %d", n.Fret)
	}

	if n.MidiPitch == 127 {
		return fmt.Errorf("cannot go above MIDI 127")
	}

	n.MidiPitch++
	n.Fret++
	return nil
}

func (n Note) ScoreTo(to Playable) float64 {
	// Расстояние по горизонтали (лады)
	fretDist := math.Abs(float64(n.Fret - to.FretPosition()))

	// Расстояние по вертикали (строки)
	stringDist := math.Abs(float64(n.String - to.StringNumber()))

	// Бонус за открытые струны
	openString := 0.0
	if n.Fret == 0 {
		openString = openStringBonus
	}

	score := (stringDist * stringlWeight) +
		(fretDist * fretWeight) +
		openString

	return score
}

type Notes []Note

func (n *Notes) ClosestTo(target Note) (Note, error) {
	if len(*n) == 0 {
		return Note{}, errors.New("empty notes list")
	}

	closest := Note{}
	minScore := math.MaxFloat64

	for _, candidate := range *n {
		currentScore := candidate.ScoreTo(target)

		if currentScore < minScore {
			minScore = currentScore
			closest = candidate
		}
	}

	closest.Time = target.Time

	return closest, nil
}
