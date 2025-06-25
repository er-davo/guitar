package guitar

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransitionCostTo(t *testing.T) {
	createNote := func(name string, fret, stringNum int) Playable {
		return Note{Name: name, Octave: 4, Fret: fret, String: stringNum}
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
			from:     Fingering{createNote("E", 0, 1)},
			to:       Fingering{createNote("E", 2, 4)},
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
				createNote("E", 0, 1),
			},
			to: Fingering{
				createNote("C", 5, 5),
				createNote("E", 0, 1),
			},
			expected: 0.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.from.TransitionCostTo(tc.to)
			require.InDelta(t, tc.expected, got, 0.01)
		})
	}
}

func TestOptimizePath(t *testing.T) {
	tun, _ := ParseTuning(StandardTuning)
	fb, _ := NewFingerBoard(tun, 24)
	optimizer := NewFingeringOptimizer(*fb)

	testCases := []struct {
		name     string
		input    [][]Playable
		wantErr  bool
		minSteps int
	}{
		{
			name: "Simple melody C-D-E",
			input: [][]Playable{
				{Note{Name: "C", Octave: 4}},
				{Note{Name: "D", Octave: 4}},
				{Note{Name: "E", Octave: 4}},
			},
			wantErr:  false,
			minSteps: 3,
		},
		{
			name: "Chord C-E",
			input: [][]Playable{
				{Note{Name: "C", Octave: 4}, Note{Name: "E", Octave: 4}},
			},
			wantErr:  false,
			minSteps: 1,
		},
		{
			name: "Invalid note name",
			input: [][]Playable{
				{Note{Name: "H", Octave: 4}},
			},
			wantErr: true,
		},
		{
			name: "Too many notes",
			input: [][]Playable{
				{
					Note{Name: "C", Octave: 4},
					Note{Name: "D", Octave: 4},
					Note{Name: "E", Octave: 4},
					Note{Name: "F", Octave: 4},
					Note{Name: "G", Octave: 4},
					Note{Name: "A", Octave: 4},
					Note{Name: "B", Octave: 4},
				},
			},
			wantErr: true,
		},
		{
			name: "Open strings preferred",
			input: [][]Playable{
				{Note{Name: "E", Octave: 4}},
				{Note{Name: "B", Octave: 3}},
			},
			wantErr:  false,
			minSteps: 2,
		},
		{
			name: "Chord transitions",
			input: [][]Playable{
				{Note{Name: "C", Octave: 4}, Note{Name: "E", Octave: 4}, Note{Name: "G", Octave: 4}},
				{Note{Name: "C", Octave: 4}},
				{Note{Name: "A", Octave: 4}, Note{Name: "F", Octave: 4}},
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
				{Note{Name: "A", Octave: 4}},
			},
			wantErr:  false,
			minSteps: 1,
		},
		{
			name: "With technique: Slide",
			input: [][]Playable{
				{Slide{
					NoteFrom: Note{Name: "A", Octave: 3},
					NoteTo:   Note{Name: "B", Octave: 3},
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
