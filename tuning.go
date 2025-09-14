package guitar

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	StandardTuning = "E4 B3 G3 D3 A2 E2"
	DropD          = "E4 B3 G3 D3 A2 D2"
)

type Tuning []Note

func (t *Tuning) NoteNames() []string {
	names := make([]string, len(*t))
	for i := range *t {
		names[i] = (*t)[i].Name()
	}
	names[0] = strings.ToLower(names[0])
	return names
}

func ParseTuning(notes string) (Tuning, error) {
	if len(notes) == 0 {
		return Tuning{}, fmt.Errorf("empty notes")
	}

	tuningNotes := strings.Split(notes, " ")
	tuning := make(Tuning, len(tuningNotes))

	for stringNumber, stringNote := range tuningNotes {
		name := stringNote[:len(stringNote)-1]
		octave, err := strconv.Atoi(stringNote[len(stringNote)-1:])
		if err != nil {
			return Tuning{}, fmt.Errorf("invalid octave at note: %s", stringNote)
		}
		note := Note{
			MidiPitch: NoteToMidi(name, octave),
			String:    stringNumber,
			Fret:      0,
		}

		if note.MidiPitch == -1 {
			return Tuning{}, fmt.Errorf("invalid note: %s", stringNote)
		}

		tuning[stringNumber] = note
	}

	return tuning, nil
}
