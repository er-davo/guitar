package guitar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTabWriter_WriteFrames(t *testing.T) {
	tun, _ := ParseTuning(StandardTuning)
	tuningNotes := tun.NoteNames()

	testCases := []struct {
		name        string
		frames      []TabFrame
		expectedTab string
		expectError bool
	}{
		{
			name: "single note",
			frames: []TabFrame{
				{Note{MidiPitch: NoteToMidi("E", 2), Fret: 12, String: 5, Time: 0}},
			},
			expectedTab: "e|--\nB|--\nG|--\nD|--\nA|--\nE|12\n",
		},
		{
			name: "chord with time step",
			frames: []TabFrame{
				{
					Note{MidiPitch: NoteToMidi("C", 2), Fret: 3, String: 4, Time: 0},
					Note{MidiPitch: NoteToMidi("E", 2), Fret: 2, String: 3, Time: 0},
					Note{MidiPitch: NoteToMidi("G", 2), Fret: 0, String: 2, Time: 0},
				},
				{
					Note{MidiPitch: NoteToMidi("C", 2), Fret: 5, String: 2, Time: 0.4},
				},
			},
			expectedTab: "e|---\nB|---\nG|0-5\nD|2--\nA|3--\nE|---\n",
		},
		{
			name: "slide",
			frames: []TabFrame{
				{
					Slide{
						NoteFrom: Note{Fret: 5, String: 2, Time: 0},
						NoteTo:   Note{Fret: 7, String: 2, Time: 0.2},
					},
				},
			},
			expectedTab: "e|---\nB|---\nG|5/7\nD|---\nA|---\nE|---\n",
		},
		{
			name: "hammer-on",
			frames: []TabFrame{
				{
					HammerOn{
						NoteFrom: Note{Fret: 2, String: 1, Time: 0},
						NoteTo:   Note{Fret: 4, String: 1, Time: 0.2},
					},
				},
			},
			expectedTab: "e|---\nB|2h4\nG|---\nD|---\nA|---\nE|---\n",
		},
		{
			name: "invalid string",
			frames: []TabFrame{
				{Note{MidiPitch: NoteToMidi("A", 2), Fret: 1, String: 8, Time: 0}},
			},
			expectError: true,
		},
		{
			name: "time out of order",
			frames: []TabFrame{
				{Note{Fret: 0, String: 5, Time: 0.5}},
				{Note{Fret: 0, String: 4, Time: 0.3}}, // earlier than previous
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			writer, _ := NewTabWriter(tuningNotes)
			err := writer.Write(tc.frames...)

			if tc.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			assert.Equal(t, tc.expectedTab, writer.Tab())
		})
	}
}
