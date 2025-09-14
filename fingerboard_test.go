package guitar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFingerBoard(t *testing.T) {
	standardTun, _ := ParseTuning(StandardTuning)
	t.Run("valid frets", func(t *testing.T) {
		fb, err := NewFingerBoard(standardTun, 24)
		assert.NoError(t, err)
		assert.Equal(t, 24, fb.frets)
	})

	t.Run("negative frets", func(t *testing.T) {
		_, err := NewFingerBoard(standardTun, -5)
		assert.ErrorContains(t, err, "frets value can not be negative")
	})
}

func TestFingerBoard_FindNotes(t *testing.T) {
	tun, _ := ParseTuning(StandardTuning)
	fb, _ := NewFingerBoard(tun, 24)

	results := fb.FindNotes(Note{MidiPitch: NoteToMidi("A", 2)})

	var found bool
	for _, p := range results {
		n := p.(Note)
		if n.Fret == 0 && n.String == 4 {
			found = true
			break
		}
	}
	assert.True(t, found, "A2 at string 4 fret 0 not found")
}

func TestFingerBoard_FindSlide(t *testing.T) {
	tun, _ := ParseTuning(StandardTuning)
	fb, _ := NewFingerBoard(tun, 24)

	slide := Slide{
		NoteFrom: Note{MidiPitch: NoteToMidi("A", 2)},
		NoteTo:   Note{MidiPitch: NoteToMidi("B", 2)},
	}

	results := fb.FindSlides(slide)

	assert.NotEmpty(t, results, "Expected at least one slide result")
	for _, p := range results {
		s := p.(Slide)
		assert.Equal(t, s.NoteFrom.String, s.NoteTo.String, "Slide should be on same string")
		assert.NotEqual(t, s.NoteFrom.Fret, s.NoteTo.Fret, "Slide should go to different fret")
	}
}

func TestFingerBoard_FindHammerOn(t *testing.T) {
	tun, _ := ParseTuning(StandardTuning)
	fb, _ := NewFingerBoard(tun, 24)

	ho := HammerOn{
		NoteFrom: Note{MidiPitch: NoteToMidi("C", 3)},
		NoteTo:   Note{MidiPitch: NoteToMidi("D", 3)},
	}

	results := fb.FindHammerOns(ho)

	assert.NotEmpty(t, results, "Expected hammer-on positions")
	for _, p := range results {
		h := p.(HammerOn)
		assert.Equal(t, h.NoteFrom.String, h.NoteTo.String, "Hammer-on must be on same string")
		assert.True(t, h.NoteTo.Fret > h.NoteFrom.Fret, "Hammer-on should go to higher fret")
	}
}

func TestFingerBoard_FindPullOff(t *testing.T) {
	tun, _ := ParseTuning(StandardTuning)
	fb, _ := NewFingerBoard(tun, 24)

	po := PullOff{
		NoteFrom: Note{MidiPitch: NoteToMidi("D", 3)},
		NoteTo:   Note{MidiPitch: NoteToMidi("C", 3)},
	}

	results := fb.FindPullOffs(po)

	assert.NotEmpty(t, results, "Expected pull-off positions")
	for _, p := range results {
		po := p.(PullOff)
		assert.Equal(t, po.NoteFrom.String, po.NoteTo.String, "Pull-off must be on same string")
		assert.True(t, po.NoteTo.Fret < po.NoteFrom.Fret, "Pull-off should go to lower fret")
	}
}

func TestFingerBoard_FindHarmonics(t *testing.T) {
	tun, _ := ParseTuning(StandardTuning)
	fb, _ := NewFingerBoard(tun, 24)

	h := Harmonic{Note: Note{MidiPitch: NoteToMidi("A", 2)}}

	results := fb.FindHarmonics(h)

	assert.NotEmpty(t, results, "Expected harmonic positions")
	for _, p := range results {
		hh := p.(Harmonic)
		assert.Equal(t, "A", hh.Name())
		assert.Equal(t, 2, hh.Octave())
	}
}

func TestFingerBoard_Find(t *testing.T) {
	tun, _ := ParseTuning(StandardTuning)
	fb, _ := NewFingerBoard(tun, 24)

	testCases := []struct {
		name     string
		input    Playable
		validate func([]Playable) bool
	}{
		{
			name:  "Find Note A2",
			input: Note{MidiPitch: NoteToMidi("A", 2)},
			validate: func(results []Playable) bool {
				for _, r := range results {
					n := r.(Note)
					if n.Name() == "A" && n.Octave() == 2 {
						return true
					}
				}
				return false
			},
		},
		{
			name: "Find Slide A2 -> B2",
			input: Slide{
				NoteFrom: Note{MidiPitch: NoteToMidi("A", 2)},
				NoteTo:   Note{MidiPitch: NoteToMidi("B", 2)},
			},
			validate: func(results []Playable) bool {
				for _, r := range results {
					s := r.(Slide)
					if s.NoteFrom.Name() == "A" && s.NoteTo.Name() == "B" &&
						s.NoteFrom.String == s.NoteTo.String &&
						s.NoteFrom.Fret != s.NoteTo.Fret {
						return true
					}
				}
				return false
			},
		},
		{
			name: "Find HammerOn C3 -> D3",
			input: HammerOn{
				NoteFrom: Note{MidiPitch: NoteToMidi("C", 3)},
				NoteTo:   Note{MidiPitch: NoteToMidi("D", 3)},
			},
			validate: func(results []Playable) bool {
				for _, r := range results {
					h := r.(HammerOn)
					if h.NoteFrom.Name() == "C" && h.NoteTo.Name() == "D" &&
						h.NoteFrom.String == h.NoteTo.String &&
						h.NoteFrom.Fret < h.NoteTo.Fret {
						return true
					}
				}
				return false
			},
		},
		{
			name: "Find PullOff D3 → C3",
			input: PullOff{
				NoteFrom: Note{MidiPitch: NoteToMidi("D", 3)},
				NoteTo:   Note{MidiPitch: NoteToMidi("C", 3)},
			},
			validate: func(results []Playable) bool {
				for _, r := range results {
					p := r.(PullOff)
					if p.NoteFrom.Name() == "D" && p.NoteTo.Name() == "C" &&
						p.NoteFrom.String == p.NoteTo.String &&
						p.NoteFrom.Fret > p.NoteTo.Fret {
						return true
					}
				}
				return false
			},
		},
		{
			name: "Find Harmonics A2",
			input: Harmonic{
				Note: Note{MidiPitch: NoteToMidi("A", 2)},
			},
			validate: func(results []Playable) bool {
				for _, r := range results {
					h := r.(Harmonic)
					if h.Name() == "A" && h.Octave() == 2 {
						return true
					}
				}
				return false
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			results := fb.Find(tc.input)
			assert.NotEmpty(t, results, "Expected some playable results")
			assert.True(t, tc.validate(results), "Validation function failed for test case")
		})
	}
}
