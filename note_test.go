package guitar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddFret(t *testing.T) {
	testCases := []struct {
		name        string
		input       Note
		expected    Note
		expectError bool
	}{
		{
			name:     "regular increment",
			input:    Note{MidiPitch: NoteToMidi("F", 2), Fret: 0},
			expected: Note{MidiPitch: NoteToMidi("F#", 2), Fret: 1},
		},
		{
			name:     "octave change",
			input:    Note{MidiPitch: NoteToMidi("B", 2), Fret: 0},
			expected: Note{MidiPitch: NoteToMidi("C", 3), Fret: 1},
		},
		{
			name:     "G# to A",
			input:    Note{MidiPitch: NoteToMidi("G#", 3), Fret: 11},
			expected: Note{MidiPitch: NoteToMidi("A", 3), Fret: 12},
		},
		{
			name:        "invalid pitch",
			input:       Note{MidiPitch: 127, Fret: 0}, // выше нельзя
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			n := tc.input
			err := n.AddFret()

			if tc.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expected.MidiPitch, n.MidiPitch)
			assert.Equal(t, tc.expected.Fret, n.Fret)

			assert.Equal(t, tc.expected.Name(), n.Name())
			assert.Equal(t, tc.expected.Octave(), n.Octave())
		})
	}
}

func TestScoreTo(t *testing.T) {
	testCases := []struct {
		name          string
		note          Note
		target        Note
		expectedScore float64
	}{
		{
			name:          "exact match",
			note:          Note{Fret: 8, String: 0},
			target:        Note{Fret: 8, String: 0},
			expectedScore: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.note.ScoreTo(tc.target)
			assert.Equal(t, tc.expectedScore, actual, "expected %f, found %f", tc.expectedScore, actual)
		})
	}
}

func TestClosestTo(t *testing.T) {
	notes := Notes{
		{MidiPitch: NoteToMidi("E", 2), Fret: 0, String: 5},
		{MidiPitch: NoteToMidi("E", 2), Fret: 2, String: 2},
		{MidiPitch: NoteToMidi("E", 2), Fret: 6, String: 4},
		{MidiPitch: NoteToMidi("E", 2), Fret: 8, String: 0},
	}

	testCases := []struct {
		name     string
		target   Note
		expected Note
	}{
		{
			name:     "exact match",
			target:   Note{Fret: 8, String: 0},
			expected: notes[3],
		},
		{
			name:     "prefer open string",
			target:   Note{Fret: 0, String: 3},
			expected: notes[0], // Should prefer open E string
		},
		{
			name:     "closest fret distance",
			target:   Note{Fret: 1, String: 5},
			expected: notes[0],
		},
		{
			name:     "tie breaker with string distance",
			target:   Note{Fret: 2, String: 3},
			expected: notes[1], // Closer string distance
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			closest, err := notes.ClosestTo(tc.target)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, closest)
		})
	}
}
