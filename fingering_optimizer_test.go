package guitar

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransitionCostTo(t *testing.T) {
	createNote := func(name string, fret, stringNum int) Playable {
		return makeNote(name, 4, fret, stringNum)
	}

	testCases := []struct {
		name     string
		from     Fingering
		to       Fingering
		expected float64
	}{
		{
			name:     "Same position",
			from:     Fingering{createNote("C", 3, 5)},
			to:       Fingering{createNote("C", 3, 5)},
			expected: 0,
		},
		{
			name:     "Different fret",
			from:     Fingering{createNote("C", 3, 5)},
			to:       Fingering{createNote("C", 5, 5)},
			expected: 2,
		},
		{
			name:     "Different string",
			from:     Fingering{createNote("C", 3, 5)},
			to:       Fingering{createNote("C", 3, 3)},
			expected: 2,
		},
		{
			name:     "Open string bonus",
			from:     Fingering{makeNote("E", 4, 0, 1)},
			to:       Fingering{makeNote("E", 4, 2, 4)},
			expected: 3.0, // 3 + 2 - 2 = 3
		},
		{
			name:     "Unmatched notes (from only)",
			from:     Fingering{createNote("C", 3, 5)},
			to:       Fingering{},
			expected: 3.0,
		},
		{
			name:     "Unmatched notes (to only)",
			from:     Fingering{},
			to:       Fingering{createNote("C", 3, 5)},
			expected: 3.0,
		},
		{
			name: "Multiple notes",
			from: Fingering{
				createNote("C", 3, 5),
				makeNote("E", 4, 0, 1),
			},
			to: Fingering{
				createNote("C", 5, 5),
				makeNote("E", 4, 0, 1),
			},
			expected: 0.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.from.NewTransitionCostTo(tc.to)
			require.InDelta(t, tc.expected, got, 0.01)
		})
	}
}

func TestOptimizePath(t *testing.T) {
	tun, _ := ParseTuning(StandardTuning)
	fb, _ := NewFingerBoard(tun, 24)
	optimizer := NewFingeringOptimizer(fb)

	testCases := []struct {
		name     string
		input    [][]Playable
		wantErr  bool
		minSteps int
	}{
		{
			name: "Simple melody C-D-E",
			input: [][]Playable{
				{makeNote("C", 4, 0, 1)},
				{makeNote("D", 4, 0, 1)},
				{makeNote("E", 4, 0, 1)},
			},
			wantErr:  false,
			minSteps: 3,
		},
		{
			name: "Chord C-E",
			input: [][]Playable{
				{makeNote("C", 4, 0, 1), makeNote("E", 4, 0, 1)},
			},
			wantErr:  false,
			minSteps: 1,
		},
		{
			name: "Invalid note name",
			input: [][]Playable{
				{Note{MidiPitch: -1}}, // специально неправильный pitch
			},
			wantErr: true,
		},
		{
			name: "Too many notes",
			input: [][]Playable{
				{
					makeNote("C", 4, 0, 1),
					makeNote("D", 4, 0, 1),
					makeNote("E", 4, 0, 1),
					makeNote("F", 4, 0, 1),
					makeNote("G", 4, 0, 1),
					makeNote("A", 4, 0, 1),
					makeNote("B", 4, 0, 1),
				},
			},
			wantErr: true,
		},
		{
			name: "Open strings preferred",
			input: [][]Playable{
				{makeNote("E", 4, 0, 1)},
				{makeNote("B", 3, 0, 2)},
			},
			wantErr:  false,
			minSteps: 2,
		},
		{
			name: "Chord transitions",
			input: [][]Playable{
				{makeNote("C", 4, 0, 1), makeNote("E", 4, 0, 1), makeNote("G", 4, 0, 1)},
				{makeNote("C", 4, 0, 1)},
				{makeNote("A", 4, 0, 1), makeNote("F", 4, 0, 1)},
			},
			wantErr:  false,
			minSteps: 3,
		},
		{
			name:    "Empty input",
			input:   [][]Playable{},
			wantErr: true,
		},
		{
			name: "Single note",
			input: [][]Playable{
				{makeNote("A", 4, 0, 1)},
			},
			wantErr:  false,
			minSteps: 1,
		},
		{
			name: "With technique: Slide",
			input: [][]Playable{
				{Slide{
					NoteFrom: makeNote("A", 3, 0, 1),
					NoteTo:   makeNote("B", 3, 2, 1),
				}},
			},
			wantErr:  false,
			minSteps: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			layers := []TimeLayer{}
			for _, slice := range tc.input {
				layer, err := optimizer.TimeLayer(slice)
				if tc.wantErr {
					require.Error(t, err)
					return
				}
				require.NoError(t, err)
				require.NotEmpty(t, layer)
				layers = append(layers, layer)
			}

			if tc.wantErr {
				return
			}

			path, err := optimizer.OptimizePath(layers)
			require.NoError(t, err)
			require.GreaterOrEqual(t, len(path), tc.minSteps)
		})
	}
}

func makeNote(name string, octave, fret, stringNum int) Note {
	return Note{
		MidiPitch: NoteToMidi(name, octave),
		Fret:      fret,
		String:    stringNum,
	}
}
