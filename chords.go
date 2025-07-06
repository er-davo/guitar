package guitar

import (
	"strconv"
	"strings"
)

func ParseChord(chordTab string, time float64) []Playable {
	notes := strings.Split(chordTab, " ")
	chord := []Playable{}

	for i, note := range notes {
		num, err := strconv.Atoi(note)
		if err != nil {
			if note != "-" {
				return []Playable{}
			}
			continue
		}
		chord = append(chord, Note{Fret: num, String: i, Time: time})
	}

	return chord
}
